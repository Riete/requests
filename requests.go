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
	contentTypeJson string = "application/json;charset=utf-8"
	contentTypeForm string = "application/x-www-form-urlencoded"
	httpGet         string = "GET"
	httpPost        string = "POST"
)

type request struct {
	httpreq *http.Request
	client  *http.Client
	httprsp *response
}

type response struct {
	resp *http.Response
}

func newRequest() *request {
	r := &request{}
	r.client = &http.Client{}
	r.httpreq, _ = http.NewRequest("", "", nil)
	r.httprsp = &response{}
	return r
}

func (r *request) do() (*response, error) {
	resp, err := r.client.Do(r.httpreq)
	r.httprsp.resp = resp
	return r.httprsp, err
}

func (r *request) get(originUrl string) (*response, error) {
	r.httpreq.Method = httpGet
	r.httpreq.URL = parseUrl(originUrl)
	return r.do()
}

func (r *request) getWithParams(originUrl string, params map[string]string) (*response, error) {
	r.httpreq.Method = httpGet
	r.httpreq.URL = parseUrl(originUrl)
	sendUrl := parseUrl(originUrl)
	p := url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	sendUrl.RawQuery = p.Encode()
	r.httpreq.URL = sendUrl
	return r.do()
}

func (r *request) post(originUrl string) (*response, error) {
	r.httpreq.Method = httpPost
	r.httpreq.URL = parseUrl(originUrl)
	return r.do()
}

func (r *request) postJson(originUrl string, data map[string]interface{}) (*response, error) {
	r.httpreq.Method = httpPost
	r.httpreq.URL = parseUrl(originUrl)
	jsonStr, _ := json.Marshal(data)
	r.httpreq.Header.Set("Content-Type", contentTypeJson)
	r.httpreq.Body = ioutil.NopCloser(bytes.NewBuffer(jsonStr))
	return r.do()
}

func (r *request) postForm(originUrl string, data map[string]string) (*response, error) {
	r.httpreq.Method = httpPost
	r.httpreq.URL = parseUrl(originUrl)
	r.httpreq.Header.Set("Content-Type", contentTypeForm)
	formData := url.Values{}
	for k, v := range data {
		formData.Add(k, v)
	}
	r.httpreq.Body = ioutil.NopCloser(strings.NewReader(formData.Encode()))
	return r.do()
}

func (r *response) Content() string {
	defer r.resp.Body.Close()
	body, _ := ioutil.ReadAll(r.resp.Body)
	return string(body)
}

func (r *response) StatusCode() int {
	return r.resp.StatusCode
}

func (r *response) Status() string {
	return r.resp.Status
}

func (r *response) Cookies() []*http.Cookie {
	return r.resp.Cookies()
}

func Get(originUrl string) (*response, error) {
	r := newRequest()
	return r.get(originUrl)
}

func GetWithParams(originUrl string, params map[string]string) (*response, error) {
	r := newRequest()
	return r.getWithParams(originUrl, params)
}

func Post(originUrl string) (*response, error) {
	r := newRequest()
	return r.post(originUrl)
}

func PostJson(originUrl string, data map[string]interface{}) (*response, error) {
	r := newRequest()
	return r.postJson(originUrl, data)
}

func PostForm(originUrl string, data map[string]string) (*response, error) {
	r := newRequest()
	return r.postForm(originUrl, data)
}
