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
	"time"
	"unsafe"

	"github.com/juju/ratelimit"
)

// DefaultClient set transport to clone of http.DefaultTransport
var DefaultClient = NewClient()

// NewTransport return clone of http.DefaultTransport
func NewTransport() *http.Transport {
	return http.DefaultTransport.(*http.Transport).Clone()
}

func NewClient() *http.Client {
	return &http.Client{Transport: NewTransport()}
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

func (r *Request) SetSkipTLS() {
	r.client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func (r *Request) UnsetProxy() {
	_ = os.Unsetenv("HTTP_PROXY")
	_ = os.Unsetenv("http_proxy")
	_ = os.Unsetenv("HTTPS_PROXY")
	_ = os.Unsetenv("https_proxy")
	r.client.Transport.(*http.Transport).Proxy = nil
}

func (r *Request) SetProxyEnv(proxy map[string]string) {
	for k, v := range proxy {
		_ = os.Setenv(k, v)
	}
}

func (r *Request) SetProxyFunc(f func(*http.Request) (*url.URL, error)) {
	r.client.Transport.(*http.Transport).Proxy = f
}

func (r *Request) SetProxyURL(proxy *url.URL) {
	r.SetProxyFunc(http.ProxyURL(proxy))
}

func (r *Request) parseURL(originURL string) error {
	var err error
	r.req.URL, err = url.Parse(originURL)
	return err
}

func (r *Request) do(options ...MethodOption) error {
	for _, option := range options {
		option(r)
	}
	var err error
	r.resp, err = r.client.Do(r.req)
	if err != nil {
		return err
	}
	defer r.resp.Body.Close()
	r.content, err = io.ReadAll(r.resp.Body)
	return err
}

func (r *Request) Content() []byte {
	return r.content
}

func (r *Request) ContentToString() string {
	return *(*string)(unsafe.Pointer(&r.content))
}

func (r *Request) Json() (map[string]any, error) {
	m := make(map[string]any)
	return m, json.Unmarshal(r.content, &m)
}

func (r *Request) JsonTo(dest any) error {
	return json.Unmarshal(r.content, dest)
}

func (r *Request) Status() (int, string) {
	return r.resp.StatusCode, r.resp.Status
}

func (r *Request) Response() *http.Response {
	return r.resp
}

func (r *Request) Request() *http.Request {
	return r.req
}

func (r *Request) Do(originURL string, options ...MethodOption) error {
	if err := r.parseURL(originURL); err != nil {
		return err
	}
	return r.do(options...)
}

func (r *Request) Get(originURL string, options ...MethodOption) error {
	return r.Do(originURL, append(options, WithMethod(http.MethodGet))...)
}

func (r *Request) Post(originURL string, options ...MethodOption) error {
	return r.Do(originURL, append(options, WithMethod(http.MethodPost))...)
}

func (r *Request) Put(originURL string, options ...MethodOption) error {
	return r.Do(originURL, append(options, WithMethod(http.MethodPut))...)
}

func (r *Request) Patch(originURL string, options ...MethodOption) error {
	return r.Do(originURL, append(options, WithMethod(http.MethodPatch))...)
}

func (r *Request) Delete(originURL string, options ...MethodOption) error {
	return r.Do(originURL, append(options, WithMethod(http.MethodDelete))...)
}

// Stream return io.ReadCloser, use ReadStream to read stream data
func (r *Request) Stream(originURL string, options ...MethodOption) (io.ReadCloser, error) {
	var err error
	r.req.Method = http.MethodGet
	if err = r.parseURL(originURL); err != nil {
		return nil, err
	}
	for _, option := range options {
		option(r)
	}
	r.resp, err = r.client.Do(r.req)
	return r.resp.Body, err
}

func (r *Request) DownloadToWriter(originURL string, w io.Writer, options ...MethodOption) (int64, error) {
	r.req.Method = http.MethodGet
	err := r.parseURL(originURL)
	if err != nil {
		return 0, err
	}
	for _, option := range options {
		option(r)
	}
	r.resp, err = r.client.Do(r.req)
	if err != nil {
		return 0, err
	}
	defer r.resp.Body.Close()
	return io.Copy(w, r.resp.Body)
}

// Download rate is download speed per second, e.g. 1024 ==> 1KiB/s, 1024*1024 ==> 1MiB/s, if rate <= 0 it means no limit
func (r *Request) Download(filePath, originURL string, rate int64, options ...MethodOption) (int64, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	if rate > 0 {
		bucket := ratelimit.NewBucketWithRate(float64(rate), rate)
		return r.DownloadToWriter(originURL, ratelimit.Writer(f, bucket), options...)
	}
	return r.DownloadToWriter(originURL, f, options...)
}

// Upload rate is upload speed per second, e.g. 1024 ==> 1KiB, 1024*1024 ==> 1MiB/s, if rate <= 0 it means no limit
func (r *Request) Upload(originURL string, data map[string]string, rate int64, fileFieldName string, filePaths []string, options ...MethodOption) error {
	r.req.Method = http.MethodPost
	if err := r.parseURL(originURL); err != nil {
		return err
	}
	for _, option := range options {
		option(r)
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

		writer, err := w.CreateFormFile(fileFieldName, filepath.Base(fp))
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
	if rate > 0 {
		bucket := ratelimit.NewBucketWithRate(float64(rate), rate)
		r.req.Body = io.NopCloser(ratelimit.Reader(body, bucket))
	}
	return r.do()
}

func (r *Request) CloseIdleConnections() {
	r.client.CloseIdleConnections()
}

// NewRequest default use DefaultClient to do http request, RequestOption can be provided to set Request properties
func NewRequest(options ...RequestOption) *Request {
	r := &Request{client: DefaultClient}
	r.req, _ = http.NewRequest("", "", nil)
	for _, option := range options {
		option(r)
	}
	return r
}

func NewSession(options ...RequestOption) *Request {
	r := NewRequest(options...)
	r.client.Jar, _ = cookiejar.New(nil)
	return r
}
