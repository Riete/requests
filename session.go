package requests

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

type Session struct {
	Request *Request
	Cookies []*http.Cookie
}

func NewSession() *Session {
	s := &Session{NewRequest(), nil}
	s.Request.Client = &http.Client{}
	jar, _ := cookiejar.New(nil)
	s.Request.Client.Jar = jar
	return s
}

func (s *Session) SetBasicAuth(username, password string) {
	s.Request.Req.SetBasicAuth(username, password)
}

func (s *Session) SetBearerTokenAuth(token string) {
	s.Request.Req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (s *Session) SetCookies(originUrl string) {
	s.Request.ParseUrl(originUrl)
	s.Request.Client.Jar.SetCookies(s.Request.Req.URL, s.Cookies)
}

func (s *Session) SetTimeout(t time.Duration) {
	s.Request.Client.Timeout = t
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
	s.Request.Client.Transport = tr
}

func (s *Session) SetProxy(proxy map[string]string) {
	for k, v := range proxy {
		err := os.Setenv(k, v)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Session) Get(originUrl string) error {
	s.SetCookies(originUrl)
	return s.Request.Get(originUrl)
}

func (s *Session) GetWithParams(originUrl string, params map[string]string) error {
	s.SetCookies(originUrl)
	return s.Request.GetWithParams(originUrl, params)
}

func (s *Session) Post(originUrl string) error {
	s.SetCookies(originUrl)
	return s.Request.Post(originUrl)
}

func (s *Session) PostJson(originUrl string, data map[string]interface{}) error {
	s.SetCookies(originUrl)
	return s.Request.PostJson(originUrl, data)
}

func (s *Session) PostForm(originUrl string, data map[string]string) error {
	s.SetCookies(originUrl)
	return s.Request.PostForm(originUrl, data)
}
