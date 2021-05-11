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
	HttpGet         string = "GET"
	HttpPost        string = "POST"
)

type Request struct {
	Req        *http.Request
	Client     *http.Client
	Resp       *http.Response
	Content    string
	StatusCode int
	Status     string
}

func NewRequest() *Request {
	r := &Request{}
	r.Client = &http.Client{}
	r.Req, _ = http.NewRequest("", "", nil)
	return r
}

func (r *Request) ParseUrl(originUrl string) {
	sendUrl, err := url.Parse(originUrl)
	if err != nil {
		panic(err)
	}
	r.Req.URL = sendUrl
}

func (r *Request) Do() error {
	resp, err := r.Client.Do(r.Req)
	if err != nil {
		return err
	}
	r.Resp = resp
	r.StatusCode = resp.StatusCode
	r.Status = resp.Status
	defer r.Resp.Body.Close()
	body, err := ioutil.ReadAll(r.Resp.Body)
	if err != nil {
		return err
	}
	r.Content = string(body)
	return nil
}

func (r *Request) Get(originUrl string) error {
	r.Req.Method = HttpGet
	r.ParseUrl(originUrl)
	return r.Do()
}

func (r *Request) GetWithParams(originUrl string, params map[string]string) error {
	r.Req.Method = HttpGet
	r.ParseUrl(originUrl)
	p := url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	r.Req.URL.RawQuery = p.Encode()
	return r.Do()
}

func (r *Request) Post(originUrl string) error {
	r.Req.Method = HttpPost
	r.ParseUrl(originUrl)
	return r.Do()
}

func (r *Request) PostJson(originUrl string, data map[string]interface{}) error {
	r.Req.Method = HttpPost
	r.ParseUrl(originUrl)
	jsonStr, _ := json.Marshal(data)
	r.Req.Header.Set("Content-Type", ContentTypeJson)
	r.Req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonStr))
	return r.Do()
}

func (r *Request) PostForm(originUrl string, data map[string]string) error {
	r.Req.Method = HttpPost
	r.ParseUrl(originUrl)
	r.Req.Header.Set("Content-Type", ContentTypeForm)
	formData := url.Values{}
	for k, v := range data {
		formData.Add(k, v)
	}
	r.Req.Body = ioutil.NopCloser(strings.NewReader(formData.Encode()))
	return r.Do()
}

func Get(originUrl string) error {
	return NewRequest().Get(originUrl)
}

func GetWithParams(originUrl string, params map[string]string) error {
	return NewRequest().GetWithParams(originUrl, params)
}

func Post(originUrl string) error {
	return NewRequest().Post(originUrl)
}

func PostJson(originUrl string, data map[string]interface{}) error {
	return NewRequest().PostJson(originUrl, data)
}

func PostForm(originUrl string, data map[string]string) error {
	return NewRequest().PostForm(originUrl, data)
}
