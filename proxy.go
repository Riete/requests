package requests

import (
	"net/http"
	"net/url"

	"golang.org/x/net/http/httpproxy"
)

const (
	ProxyHttp  = "http"
	ProxySocks = "socks5"
)

type Proxy map[string]string

func newProxy(prefix, proxy string) Proxy {
	p := prefix + "://" + proxy
	return Proxy{"http_proxy": p, "https_proxy": p}
}

func newProxyWithAuth(prefix, proxy, username, password string) Proxy {
	auth := url.QueryEscape(username) + ":" + url.QueryEscape(password)
	p := prefix + "://" + auth + "@" + proxy
	return Proxy{"http_proxy": p, "https_proxy": p}
}

func NewHttpProxy(proxy string) Proxy {
	return newProxy(ProxyHttp, proxy)
}

func NewSocksProxy(proxy string) Proxy {
	return newProxy(ProxySocks, proxy)
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
