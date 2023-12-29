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

func WithProxyEnv(proxy Proxy) Option {
	return func(r *Request) {
		r.SetProxyEnv(proxy)
	}
}

func WithProxyUrl(proxy *url.URL) Option {
	return func(r *Request) {
		r.SetProxyUrl(proxy)
	}
}

func WithProxyFunc(f func(*http.Request) (*url.URL, error)) Option {
	return func(r *Request) {
		r.SetProxyFunc(f)
	}
}

func WithUnsetProxy() Option {
	return func(r *Request) {
		r.UnsetProxy()
	}
}

func WithSkipTLS() Option {
	return func(r *Request) {
		r.SetSkipTLS()
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

func WithDefaultTransport() Option {
	return func(r *Request) {
		r.SetTransport(DefaultTransport)
	}
}

func WithClient(client *http.Client) Option {
	return func(r *Request) {
		r.SetClient(client)
	}
}

func WithDefaultClient() Option {
	return func(r *Request) {
		r.SetClient(DefaultClient)
	}
}
