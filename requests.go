package requests

import (
	"bytes"
	"crypto/tls"
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

var (
	// DefaultTransport is clone of http.DefaultTransport
	DefaultTransport = NewTransport()
	// DefaultClient set Transport to DefaultTransport
	DefaultClient = NewClient()
)

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

func (r *Request) SetProxyEnv(proxy Proxy) {
	for k, v := range proxy {
		_ = os.Setenv(k, v)
	}
}

func (r *Request) SetProxyFunc(f func(*http.Request) (*url.URL, error)) {
	r.client.Transport.(*http.Transport).Proxy = f

}

func (r *Request) SetProxyUrl(proxy *url.URL) {
	r.SetProxyFunc(http.ProxyURL(proxy))
}

func (r *Request) parseUrl(originUrl string) error {
	if sendUrl, err := url.Parse(originUrl); err != nil {
		return err
	} else {
		r.req.URL = sendUrl
		return nil
	}
}

func (r *Request) do(options ...MethodOption) error {
	for _, option := range options {
		option(r)
	}
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

func (r *Request) Get(originUrl string, options ...MethodOption) error {
	r.req.Method = http.MethodGet
	if err := r.parseUrl(originUrl); err != nil {
		return err
	}
	return r.do(options...)
}

func (r *Request) Post(originUrl string, options ...MethodOption) error {
	r.req.Method = http.MethodPost
	if err := r.parseUrl(originUrl); err != nil {
		return err
	}
	return r.do(options...)
}

func (r *Request) Put(originUrl string, options ...MethodOption) error {
	r.req.Method = http.MethodPut
	if err := r.parseUrl(originUrl); err != nil {
		return err
	}
	return r.do(options...)
}

func (r *Request) Delete(originUrl string, options ...MethodOption) error {
	r.req.Method = http.MethodDelete
	if err := r.parseUrl(originUrl); err != nil {
		return err
	}
	return r.do(options...)
}

func (r *Request) DownloadToWriter(originUrl string, w io.Writer) error {
	r.req.Method = http.MethodGet
	if err := r.parseUrl(originUrl); err != nil {
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

// Download rate is download speed per second, e.g. 1024 ==> 1KiB/s, 1024*1024 ==> 1MiB/s, if rate <= 0 it means no limit
func (r *Request) Download(filePath, originUrl string, rate int64) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	if rate > 0 {
		bucket := ratelimit.NewBucketWithRate(float64(rate), rate)
		return r.DownloadToWriter(originUrl, ratelimit.Writer(f, bucket))
	}
	return r.DownloadToWriter(originUrl, f)
}

// Upload rate is upload speed per second, e.g. 1024 ==> 1KiB, 1024*1024 ==> 1MiB/s, if rate <= 0 it means no limit
func (r *Request) Upload(originUrl string, data map[string]string, rate int64, filePaths ...string) error {
	r.req.Method = http.MethodPost
	if err := r.parseUrl(originUrl); err != nil {
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
	if rate > 0 {
		bucket := ratelimit.NewBucketWithRate(float64(rate), rate)
		r.req.Body = io.NopCloser(ratelimit.Reader(body, bucket))
	}
	return r.do()
}

func (r *Request) CloseIdleConnections() {
	r.client.CloseIdleConnections()
}

// NewRequest use DefaultClient to do http request, RequestOption can be provided to set Request properties
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
