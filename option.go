package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RequestOption func(*Request)

func WithTimeout(t time.Duration) RequestOption {
	return func(r *Request) {
		r.SetTimeout(t)
	}
}

func WithHeader(headers ...map[string]string) RequestOption {
	return func(r *Request) {
		for _, header := range headers {
			r.SetHeader(header)
		}
	}
}

func WithProxyEnv(proxy map[string]string) RequestOption {
	return func(r *Request) {
		r.SetProxyEnv(proxy)
	}
}

func WithProxyURL(proxy *url.URL) RequestOption {
	return func(r *Request) {
		r.SetProxyURL(proxy)
	}
}

func WithProxyFunc(f func(*http.Request) (*url.URL, error)) RequestOption {
	return func(r *Request) {
		r.SetProxyFunc(f)
	}
}

func WithUnsetProxy() RequestOption {
	return func(r *Request) {
		r.UnsetProxy()
	}
}

func WithSkipTLS() RequestOption {
	return func(r *Request) {
		r.SetSkipTLS()
	}
}

func WithBasicAuth(username, password string) RequestOption {
	return func(r *Request) {
		r.SetBasicAuth(username, password)
	}
}

func WithBearerTokenAuth(token string) RequestOption {
	return func(r *Request) {
		r.SetBearerTokenAuth(token)
	}
}

func WithTransport(tr http.RoundTripper) RequestOption {
	return func(r *Request) {
		r.SetTransport(tr)
	}
}

func WithDefaultTransport() RequestOption {
	return func(r *Request) {
		r.SetTransport(DefaultTransport)
	}
}

func WithClient(client *http.Client) RequestOption {
	return func(r *Request) {
		r.SetClient(client)
	}
}

func WithDefaultClient() RequestOption {
	return func(r *Request) {
		r.SetClient(DefaultClient)
	}
}

type MethodOption func(r *Request)

func WithMethod(method string) MethodOption {
	return func(r *Request) {
		r.req.Method = method
	}
}

func WithParams(params map[string]string) MethodOption {
	return func(r *Request) {
		p := url.Values{}
		for k, v := range params {
			p.Add(k, v)
		}
		r.req.URL.RawQuery = p.Encode()
	}
}

func WithJsonData(data map[string]any) MethodOption {
	return func(r *Request) {
		jsonStr, _ := json.Marshal(data)
		r.req.Header.Set("Content-Type", "application/json;charset=utf-8")
		r.req.Body = io.NopCloser(bytes.NewBuffer(jsonStr))
	}
}

func WithFormData(data map[string]string) MethodOption {
	return func(r *Request) {
		formData := url.Values{}
		for k, v := range data {
			formData.Add(k, v)
		}
		r.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.req.Body = io.NopCloser(strings.NewReader(formData.Encode()))
	}
}
