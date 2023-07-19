package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/juju/ratelimit"
)

const (
	ContentTypeJson string = "application/json;charset=utf-8"
	ContentTypeForm string = "application/x-www-form-urlencoded"
	HttpGet         string = "GET"
	HttpPost        string = "POST"
	HttpPut         string = "PUT"
	HttpDelete      string = "DELETE"
)

type Option func(*Request)

func WithTimeout(t time.Duration) Option {
	return func(r *Request) {
		r.SetTimeout(t)
	}
}

func WithHeader(header map[string]string) Option {
	return func(r *Request) {
		r.SetHeader(header)
	}
}

func WithProxy(proxy map[string]string) Option {
	return func(r *Request) {
		r.SetProxy(proxy)
	}
}

func WithSkipTLSVerify() Option {
	return func(r *Request) {
		r.SkipTLSVerify()
	}
}

func WithBasicAuth(username, password string) Option {
	return func(r *Request) {
		r.SetBasicAuth(username, password)
	}
}

func WithBearerTokenAuth(token string) Option {
	return func(r *Request) {
		r.SetBearerTokenAuth(token)
	}
}

type Request struct {
	Req        *http.Request
	Client     *http.Client
	Resp       *http.Response
	Content    []byte
	StatusCode int
	Status     string
}

func NewRequest(options ...Option) *Request {
	r := &Request{Client: &http.Client{}}
	r.Req, _ = http.NewRequest("", "", nil)
	for _, option := range options {
		option(r)
	}
	return r
}

func NewSession(options ...Option) *Request {
	r := NewRequest(options...)
	r.Client.Jar, _ = cookiejar.New(nil)
	return r
}

func (r *Request) SetHeader(headers map[string]string) {
	for k, v := range headers {
		r.Req.Header.Set(k, v)
	}
}

func (r *Request) SetBasicAuth(username, password string) {
	r.Req.SetBasicAuth(username, password)
}

func (r *Request) SetBearerTokenAuth(token string) {
	r.Req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (r *Request) SetTimeout(t time.Duration) {
	r.Client.Timeout = t
}

func (r *Request) SkipTLSVerify() {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	r.Client.Transport = tr
}

func (r *Request) SetProxy(proxy map[string]string) {
	for k, v := range proxy {
		_ = os.Setenv(k, v)
	}
}

func (r *Request) ParseUrl(originUrl string) error {
	if sendUrl, err := url.Parse(originUrl); err != nil {
		return err
	} else {
		r.Req.URL = sendUrl
		return nil
	}
}

func (r *Request) do() error {
	resp, err := r.Client.Do(r.Req)
	if err != nil {
		return err
	}
	r.Resp = resp
	r.StatusCode = resp.StatusCode
	r.Status = resp.Status
	defer r.Resp.Body.Close()
	r.Content, err = io.ReadAll(r.Resp.Body)
	return err
}

func (r *Request) Get(originUrl string, params map[string]string) error {
	r.Req.Method = HttpGet
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	p := url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	r.Req.URL.RawQuery = p.Encode()
	return r.do()
}

func (r *Request) Post(originUrl string, data map[string]interface{}) error {
	r.Req.Method = HttpPost
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.Req.Header.Set("Content-Type", ContentTypeJson)
	r.Req.Body = io.NopCloser(bytes.NewBuffer(jsonStr))
	return r.do()
}

func (r *Request) PostForm(originUrl string, data map[string]string) error {
	r.Req.Method = HttpPost
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	r.Req.Header.Set("Content-Type", ContentTypeForm)
	formData := url.Values{}
	for k, v := range data {
		formData.Add(k, v)
	}
	r.Req.Body = io.NopCloser(strings.NewReader(formData.Encode()))
	return r.do()
}

func (r *Request) Put(originUrl string, data map[string]interface{}) error {
	r.Req.Method = HttpPut
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	jsonStr, _ := json.Marshal(data)
	r.Req.Header.Set("Content-Type", ContentTypeJson)
	r.Req.Body = io.NopCloser(bytes.NewBuffer(jsonStr))
	return r.do()
}

func (r *Request) Delete(originUrl string, data map[string]interface{}) error {
	r.Req.Method = HttpDelete
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	jsonStr, _ := json.Marshal(data)
	r.Req.Header.Set("Content-Type", ContentTypeJson)
	r.Req.Body = io.NopCloser(bytes.NewBuffer(jsonStr))
	return r.do()
}

func (r *Request) download(originUrl string, w io.Writer) error {
	r.Req.Method = HttpGet
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	resp, err := r.Client.Do(r.Req)
	if err != nil {
		return err
	}
	r.Resp = resp
	r.StatusCode = resp.StatusCode
	r.Status = resp.Status
	defer r.Resp.Body.Close()
	_, err = io.Copy(w, r.Resp.Body)
	return err
}

func (r *Request) Download(filePath, originUrl string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return r.download(originUrl, f)
}

func (r *Request) DownloadWithRateLimit(filePath, originUrl string, rate int64) error {
	if rate <= 0 {
		return r.Download(filePath, originUrl)
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	bucket := ratelimit.NewBucketWithRate(float64(rate), rate)
	w := ratelimit.Writer(f, bucket)
	return r.download(originUrl, w)
}

func (r Request) ContentToString() string {
	return *(*string)(unsafe.Pointer(&r.Content))
}
