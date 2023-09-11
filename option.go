package requests

import (
	"net/http"
	"net/url"
	"time"
)

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

func WithProxy(proxy Proxy) Option {
	return func(r *Request) {
		r.SetProxy(proxy)
	}
}

func WithProxyFunc(f func(*http.Request) (*url.URL, error)) Option {
	return func(r *Request) {
		r.client.Transport.(*http.Transport).Proxy = f
	}
}

func WithUnsetProxy() Option {
	return func(r *Request) {
		r.UnsetProxy()
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

func WithTransport(tr http.RoundTripper) Option {
	return func(r *Request) {
		r.SetTransport(tr)
	}
}
