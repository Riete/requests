package requests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	ContentTypeJson string = "application/json;charset=utf-8"
	ContentTypeForm string = "application/x-www-form-urlencoded"
)

type Response struct {
	Resp *http.Response
}

func (r *Response) Content() string {
	defer r.Resp.Body.Close()
	body, _ := ioutil.ReadAll(r.Resp.Body)
	return string(body)
}

func (r *Response) StatusCode() int {
	return r.Resp.StatusCode
}

func (r *Response) Status() string {
	return r.Resp.Status
}

func (r *Response) Cookies() []*http.Cookie {
	return r.Resp.Cookies()
}

func Get(originUrl string) (*Response, error) {
	r := &Response{}
	resp, err := http.Get(originUrl)
	if err != nil {
		return nil, err
	}
	r.Resp = resp
	return r, err
}

func GetWithParams(originUrl string, params map[string]string) (*Response, error) {
	r := &Response{}
	sendUrl, err := url.Parse(originUrl)
	if err != nil {
		return nil, err
	}
	p := url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	sendUrl.RawQuery = p.Encode()
	resp, err := http.Get(sendUrl.String())
	if err != nil {
		return nil, err
	}
	r.Resp = resp
	return r, nil
}

func Post(originUrl string) (*Response, error) {
	r := &Response{}
	resp, err := http.Post(originUrl, ContentTypeJson, nil)
	if err != nil {
		return nil, err
	}
	r.Resp = resp
	return r, nil
}

func PostJson(originUrl string, data map[string]interface{}) (*Response, error) {
	r := &Response{}
	jsonStr, _ := json.Marshal(data)
	resp, err := http.Post(originUrl, ContentTypeJson, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	r.Resp = resp
	return r, nil
}

func PostForm(originUrl string, data map[string]string) (*Response, error) {
	r := &Response{}
	formData := url.Values{}
	for k, v := range data {
		formData.Add(k, v)
	}
	resp, err := http.PostForm(originUrl, formData)
	if err != nil {
		return nil, err
	}
	r.Resp = resp
	return r, nil
}
