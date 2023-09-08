package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"github.com/juju/ratelimit"
	"golang.org/x/net/http/httpproxy"
)

const (
	ContentTypeJson string = "application/json;charset=utf-8"
	ContentTypeForm string = "application/x-www-form-urlencoded"
	HttpGet         string = "GET"
	HttpPost        string = "POST"
	HttpPut         string = "PUT"
	HttpDelete      string = "DELETE"
)

func ProxyFromEnvironment(req *http.Request) (*url.URL, error) {
	return httpproxy.FromEnvironment().ProxyFunc()(req.URL)
}

func DefaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     false,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

type Option func(*Request)

func WithTimeout(t time.Duration) Option {
	return func(r *Request) {
		r.SetTimeout(t)
	}
}

func WithHeader(headers ...map[string]string) Option {
	return func(r *Request) {
		for _, header := range headers {
			r.SetHeader(header)
		}
	}
}

func WithProxy(proxy map[string]string) Option {
	return func(r *Request) {
		r.SetProxy(proxy)
	}
}

func WithCancelProxy() Option {
	return func(r *Request) {
		r.CancelProxy()
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

func WithTransport(transport *http.Transport) Option {
	return func(r *Request) {
		r.SetTransport(transport)
	}
}

type Request struct {
	req     *http.Request
	client  *http.Client
	resp    *http.Response
	content []byte
}

func (r *Request) SetHeader(header map[string]string) {
	for k, v := range header {
		r.req.Header.Set(k, v)
	}
}

func (r *Request) SetBasicAuth(username, password string) {
	r.req.SetBasicAuth(username, password)
}

func (r *Request) SetBearerTokenAuth(token string) {
	r.req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (r *Request) SetTimeout(t time.Duration) {
	r.client.Timeout = t
}

func (r *Request) SetTransport(transport *http.Transport) {
	r.client.Transport = transport
}

func (r *Request) SkipTLSVerify() {
	transport := DefaultTransport()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	r.SetTransport(transport)
}

func (r *Request) CancelProxy() {
	_ = os.Setenv("HTTP_PROXY", "")
	_ = os.Setenv("http_proxy", "")
	_ = os.Setenv("HTTPS_PROXY", "")
	_ = os.Setenv("https_proxy", "")

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
		r.req.URL = sendUrl
		return nil
	}
}

func (r *Request) do() error {
	resp, err := r.client.Do(r.req)
	if err != nil {
		return err
	}
	r.resp = resp
	defer r.resp.Body.Close()
	r.content, err = io.ReadAll(r.resp.Body)
	return err
}

func (r Request) Content() []byte {
	return r.content
}

func (r Request) ContentToString() string {
	return *(*string)(unsafe.Pointer(&r.content))
}

func (r Request) Status() (int, string) {
	return r.resp.StatusCode, r.resp.Status
}

func (r Request) Response() *http.Response {
	return r.resp
}

func (r Request) Request() *http.Request {
	return r.req
}

func (r *Request) Get(originUrl string, params map[string]string) error {
	r.req.Method = HttpGet
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	p := url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	r.req.URL.RawQuery = p.Encode()
	return r.do()
}

func (r *Request) Post(originUrl string, data map[string]interface{}) error {
	r.req.Method = HttpPost
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.req.Header.Set("Content-Type", ContentTypeJson)
	r.req.Body = io.NopCloser(bytes.NewBuffer(jsonStr))
	return r.do()
}

func (r *Request) PostForm(originUrl string, data map[string]string) error {
	r.req.Method = HttpPost
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	r.req.Header.Set("Content-Type", ContentTypeForm)
	formData := url.Values{}
	for k, v := range data {
		formData.Add(k, v)
	}
	r.req.Body = io.NopCloser(strings.NewReader(formData.Encode()))
	return r.do()
}

func (r *Request) Put(originUrl string, data map[string]interface{}) error {
	r.req.Method = HttpPut
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	jsonStr, _ := json.Marshal(data)
	r.req.Header.Set("Content-Type", ContentTypeJson)
	r.req.Body = io.NopCloser(bytes.NewBuffer(jsonStr))
	return r.do()
}

func (r *Request) Delete(originUrl string, data map[string]interface{}) error {
	r.req.Method = HttpDelete
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	jsonStr, _ := json.Marshal(data)
	r.req.Header.Set("Content-Type", ContentTypeJson)
	r.req.Body = io.NopCloser(bytes.NewBuffer(jsonStr))
	return r.do()
}

func (r *Request) download(originUrl string, w io.Writer) error {
	r.req.Method = HttpGet
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	resp, err := r.client.Do(r.req)
	if err != nil {
		return err
	}
	r.resp = resp
	defer r.resp.Body.Close()
	_, err = io.Copy(w, r.resp.Body)
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

func (r *Request) Upload(originUrl string, data map[string]string, filePaths ...string) error {
	r.req.Method = HttpPost
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for k, v := range data {
		if err := w.WriteField(k, v); err != nil {
			return err
		}
	}

	var fileCloser []io.Closer
	defer func() {
		for _, f := range fileCloser {
			_ = f.Close()
		}
	}()

	for _, fp := range filePaths {
		f, err := os.Open(fp)
		if err != nil {
			return err
		}
		fileCloser = append(fileCloser, f)

		writer, err := w.CreateFormFile("file", filepath.Base(fp))
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, f)
		if err != nil {
			return err
		}
	}

	if err := w.Close(); err != nil {
		return err
	}

	r.req.Header.Set("Content-Type", w.FormDataContentType())
	r.req.Body = io.NopCloser(body)
	return r.do()
}

func NewRequest(options ...Option) *Request {
	r := &Request{client: &http.Client{}}
	r.client.Transport = DefaultTransport()
	r.req, _ = http.NewRequest("", "", nil)
	for _, option := range options {
		option(r)
	}
	return r
}

func NewSession(options ...Option) *Request {
	r := NewRequest(options...)
	r.client.Jar, _ = cookiejar.New(nil)
	return r
}
