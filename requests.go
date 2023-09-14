package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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
	ContentTypeJson = "application/json;charset=utf-8"
	ContentTypeForm = "application/x-www-form-urlencoded"
	ProxyHttp       = "http"
	ProxySocks      = "socks5"
)

type Proxy map[string]string

func newProxy(prefix, proxy string) Proxy {
	p := prefix + "://" + proxy
	return Proxy{"http_proxy": p, "https_proxy": p}
}

func NewHttpProxy(proxy string) Proxy {
	return newProxy(ProxyHttp, proxy)
}

func NewSocksProxy(proxy string) Proxy {
	return newProxy(ProxySocks, proxy)
}

func newProxyWithAuth(prefix, proxy, username, password string) Proxy {
	auth := url.QueryEscape(username) + ":" + url.QueryEscape(password)
	p := prefix + "://" + auth + "@" + proxy
	return Proxy{"http_proxy": p, "https_proxy": p}
}

func NewHttpProxyWithAuth(proxy, username, password string) Proxy {
	return newProxyWithAuth(ProxyHttp, proxy, username, password)
}

func NewSocksProxyWithAuth(proxy, username, password string) Proxy {
	return newProxyWithAuth(ProxySocks, proxy, username, password)
}

// ProxyFromEnvironment read proxy form env for every request
// http.ProxyFromEnvironment read only once
func ProxyFromEnvironment(req *http.Request) (*url.URL, error) {
	return httpproxy.FromEnvironment().ProxyFunc()(req.URL)
}

// DefaultTransport is clone of http.DefaultTransport
var DefaultTransport = http.DefaultTransport.(*http.Transport).Clone()

var DefaultClient = &http.Client{Transport: DefaultTransport}

func NewDefaultTransport() *http.Transport {
	return http.DefaultTransport.(*http.Transport).Clone()
}

func NewDefaultClient() *http.Client {
	return &http.Client{Transport: NewDefaultTransport()}
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

func (r *Request) SetTransport(rt http.RoundTripper) {
	r.client.Transport = rt
}

func (r *Request) SetClient(client *http.Client) {
	r.client = client
}

func (r *Request) SkipTLSVerify() {
	r.client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func (r *Request) UnsetProxy() {
	_ = os.Unsetenv("HTTP_PROXY")
	_ = os.Unsetenv("http_proxy")
	_ = os.Unsetenv("HTTPS_PROXY")
	_ = os.Unsetenv("https_proxy")
}

func (r *Request) SetProxy(proxy Proxy) {
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
	r.req.Method = http.MethodGet
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
	r.req.Method = http.MethodPost
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
	r.req.Method = http.MethodPost
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
	r.req.Method = http.MethodPut
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	jsonStr, _ := json.Marshal(data)
	r.req.Header.Set("Content-Type", ContentTypeJson)
	r.req.Body = io.NopCloser(bytes.NewBuffer(jsonStr))
	return r.do()
}

func (r *Request) Delete(originUrl string, data map[string]interface{}) error {
	r.req.Method = http.MethodDelete
	if err := r.ParseUrl(originUrl); err != nil {
		return err
	}
	jsonStr, _ := json.Marshal(data)
	r.req.Header.Set("Content-Type", ContentTypeJson)
	r.req.Body = io.NopCloser(bytes.NewBuffer(jsonStr))
	return r.do()
}

func (r *Request) download(originUrl string, w io.Writer) error {
	r.req.Method = http.MethodGet
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
	r.req.Method = http.MethodPost
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
	r := &Request{client: DefaultClient}
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
