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

type session struct {
	cookies []*http.Cookie
	request
}

func NewSession() *session {
	s := &session{}
	s.client = &http.Client{}
	s.httprsp = &response{}
	s.httpreq, _ = http.NewRequest("", "", nil)
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

func (s *session) SetBasicAuth(username, password string) {
	s.httpreq.SetBasicAuth(username, password)
}

func (s *session) SetBearTokenAuth(token string) {
	s.httpreq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (s *session) JsonAuth(originUrl string, data map[string]interface{}) (*response, error) {
	resp, err := PostJson(originUrl, data)
	if err != nil {
		return nil, err
	}
	s.httprsp = resp
	s.cookies = s.httprsp.Cookies()
	return s.httprsp, nil
}

func (s *session) FormAuth(originUrl string, data map[string]string) (*response, error) {
	resp, err := PostForm(originUrl, data)
	if err != nil {
		return nil, err
	}
	s.httprsp = resp
	s.cookies = s.httprsp.Cookies()
	return s.httprsp, nil
}

func (s *session) setCookies(originUrl string) {
	sendUrl := parseUrl(originUrl)
	s.client.Jar.SetCookies(sendUrl, s.cookies)
	s.httpreq.URL = sendUrl
}

func (s *session) do() (*response, error) {
	resp, err := s.client.Do(s.httpreq)
	s.httprsp.resp = resp
	return s.httprsp, err
}

func (s *session) SetTimeout(t time.Duration) {
	s.client.Timeout = t * time.Second
}

func (s *session) SkipTLSVerify() {
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

func (s *session) SetProxy(proxy map[string]string) {
	for k, v := range proxy {
		err := os.Setenv(k, v)
		if err != nil {
			panic(err)
		}
	}
}

func (s *session) Get(originUrl string) (*response, error) {
	s.setCookies(originUrl)
	s.httpreq.Method = httpGet
	return s.do()
}

func (s *session) GetWithParams(originUrl string, params map[string]string) (*response, error) {
	s.setCookies(originUrl)
	s.httpreq.Method = httpGet
	sendUrl := parseUrl(originUrl)
	p := url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	sendUrl.RawQuery = p.Encode()
	s.httpreq.URL = sendUrl
	return s.do()
}

func (s *session) Post(originUrl string) (*response, error) {
	s.setCookies(originUrl)
	s.httpreq.Method = httpPost
	return s.do()
}

func (s *session) PostJson(originUrl string, data map[string]interface{}) (*response, error) {
	s.setCookies(originUrl)
	s.httpreq.Method = httpPost
	jsonStr, _ := json.Marshal(data)
	s.httpreq.Header.Set("Content-Type", contentTypeJson)
	s.httpreq.Body = ioutil.NopCloser(bytes.NewBuffer(jsonStr))
	return s.do()
}

func (s *session) PostForm(originUrl string, data map[string]string) (*response, error) {
	s.setCookies(originUrl)
	s.httpreq.Method = httpPost
	s.httpreq.Header.Set("Content-Type", contentTypeForm)
	formData := url.Values{}
	for k, v := range data {
		formData.Add(k, v)
	}
	s.httpreq.Body = ioutil.NopCloser(strings.NewReader(formData.Encode()))
	return s.do()
}
