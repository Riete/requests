package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	GET  string = "GET"
	POST string = "POST"
)

type Session struct {
	cookies []*http.Cookie
	Request
}

func NewSession() *Session {
	s := &Session{}
	s.client = &http.Client{}
	s.response = &Response{}
	s.request, _ = http.NewRequest("", "", nil)
	jar, _ := cookiejar.New(nil)
	s.client.Jar = jar
	return s
}

func parseUrl(originUrl string) *url.URL {
	sendUrl, err := url.Parse(originUrl)
	if err != nil {
		panic(err)
	}
	return sendUrl
}

func (s *Session) SetBasicAuth(username, password string) {
	s.request.SetBasicAuth(username, password)
}

func (s *Session) SetBearTokenAuth(token string) {
	s.request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (s *Session) JsonAuth(originUrl string, data map[string]interface{}) (*Response, error) {
	resp, err := PostJson(originUrl, data)
	if err != nil {
		return nil, err
	}
	s.response = resp
	s.cookies = s.response.Cookies()
	return s.response, nil
}

func (s *Session) FormAuth(originUrl string, data map[string]string) (*Response, error) {
	resp, err := PostForm(originUrl, data)
	if err != nil {
		return nil, err
	}
	s.response = resp
	s.cookies = s.response.Cookies()
	return s.response, nil
}

func (s *Session) setCookies(originUrl string) {
	sendUrl := parseUrl(originUrl)
	s.client.Jar.SetCookies(sendUrl, s.cookies)
	s.request.URL = sendUrl
}

func (s *Session) do() (*Response, error) {
	resp, err := s.client.Do(s.request)
	s.response.resp = resp
	return s.response, err
}

func (s *Session) SetTimeout(t time.Duration) {
	s.client.Timeout = t * time.Second
}

func (s *Session) SkipTLSVerify() {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	s.client.Transport = tr
}

func (s *Session) SetProxy(proxy map[string]string) {
	for k, v := range proxy {
		err := os.Setenv(k, v)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Session) Get(originUrl string) (*Response, error) {
	s.setCookies(originUrl)
	s.request.Method = GET
	return s.do()
}

func (s *Session) GetWithParams(originUrl string, params map[string]string) (*Response, error) {
	s.setCookies(originUrl)
	s.request.Method = GET
	sendUrl := parseUrl(originUrl)
	p := url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	sendUrl.RawQuery = p.Encode()
	s.request.URL = sendUrl
	return s.do()
}

func (s *Session) Post(originUrl string) (*Response, error) {
	s.setCookies(originUrl)
	s.request.Method = POST
	return s.do()
}

func (s *Session) PostJson(originUrl string, data map[string]interface{}) (*Response, error) {
	s.setCookies(originUrl)
	s.request.Method = POST
	jsonStr, _ := json.Marshal(data)
	s.request.Header.Set("Content-Type", ContentTypeJson)
	s.request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonStr))
	return s.do()
}

func (s *Session) PostForm(originUrl string, data map[string]string) (*Response, error) {
	s.setCookies(originUrl)
	s.request.Method = POST
	s.request.Header.Set("Content-Type", ContentTypeForm)
	formData := url.Values{}
	for k, v := range data {
		formData.Add(k, v)
	}
	s.request.Body = ioutil.NopCloser(strings.NewReader(formData.Encode()))
	return s.do()
}
