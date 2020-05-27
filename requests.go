package requests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	ContentTypeJson string = "application/json;charset=utf-8"
	ContentTypeForm string = "application/x-www-form-urlencoded"
)

type Request struct {
	request  *http.Request
	client   *http.Client
	response *Response
}

type Response struct {
	resp *http.Response
}

func NewRequest() *Request {
	r := &Request{}
	r.client = &http.Client{}
	r.request, _ = http.NewRequest("", "", nil)
	r.response = &Response{}
	return r
}

func (r *Request) do() (*Response, error) {
	resp, err := r.client.Do(r.request)
	r.response.resp = resp
	return r.response, err
}

func (r *Request) Get(originUrl string) (*Response, error) {
	r.request.Method = GET
	r.request.URL = parseUrl(originUrl)
	return r.do()
}

func (r *Request) GetWithParams(originUrl string, params map[string]string) (*Response, error) {
	r.request.Method = GET
	r.request.URL = parseUrl(originUrl)
	sendUrl := parseUrl(originUrl)
	p := url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	sendUrl.RawQuery = p.Encode()
	r.request.URL = sendUrl
	return r.do()
}

func (r *Request) Post(originUrl string) (*Response, error) {
	r.request.Method = POST
	r.request.URL = parseUrl(originUrl)
	return r.do()
}

func (r *Request) PostJson(originUrl string, data map[string]interface{}) (*Response, error) {
	r.request.Method = POST
	r.request.URL = parseUrl(originUrl)
	jsonStr, _ := json.Marshal(data)
	r.request.Header.Set("Content-Type", ContentTypeJson)
	r.request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonStr))
	return r.do()
}

func (r *Request) PostForm(originUrl string, data map[string]string) (*Response, error) {
	r.request.Method = POST
	r.request.URL = parseUrl(originUrl)
	r.request.Header.Set("Content-Type", ContentTypeForm)
	formData := url.Values{}
	for k, v := range data {
		formData.Add(k, v)
	}
	r.request.Body = ioutil.NopCloser(strings.NewReader(formData.Encode()))
	return r.do()
}

func (r *Response) Content() string {
	defer r.resp.Body.Close()
	body, _ := ioutil.ReadAll(r.resp.Body)
	return string(body)
}

func (r *Response) StatusCode() int {
	return r.resp.StatusCode
}

func (r *Response) Status() string {
	return r.resp.Status
}

func (r *Response) Cookies() []*http.Cookie {
	return r.resp.Cookies()
}

func Get(originUrl string) (*Response, error) {
	r := NewRequest()
	return r.Get(originUrl)
}

func GetWithParams(originUrl string, params map[string]string) (*Response, error) {
	r := NewRequest()
	return r.GetWithParams(originUrl, params)
}

func Post(originUrl string) (*Response, error) {
	r := NewRequest()
	return r.Post(originUrl)
}

func PostJson(originUrl string, data map[string]interface{}) (*Response, error) {
	r := NewRequest()
	return r.PostJson(originUrl, data)
}

func PostForm(originUrl string, data map[string]string) (*Response, error) {
	r := NewRequest()
	return r.PostForm(originUrl, data)
}
