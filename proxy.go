package requests

import (
	"net/http"
	"net/url"

	"golang.org/x/net/http/httpproxy"
)

type Auth struct {
	Username string
	Password string
}

// Proxy addr is host:port
type Proxy struct {
	addr     string
	scheme   string
	auth     *Auth
	proxyURL string
}

func (p *Proxy) buildProxyURL() {
	if p.auth != nil {
		credential := url.QueryEscape(p.auth.Username) + ":" + url.QueryEscape(p.auth.Password)
		p.proxyURL = p.scheme + "://" + credential + "@" + p.addr
	} else {
		p.proxyURL = p.scheme + "://" + p.addr
	}
}

func (p *Proxy) ProxyRawURL() string {
	return p.proxyURL
}

func (p *Proxy) ProxyURL() (*url.URL, error) {
	return url.Parse(p.proxyURL)
}

func (p *Proxy) ProxyEnv() map[string]string {
	return map[string]string{"http_proxy": p.proxyURL, "https_proxy": p.proxyURL}
}

func NewProxy(scheme, addr string, auth *Auth) *Proxy {
	p := &Proxy{addr: addr, scheme: scheme, auth: auth}
	p.buildProxyURL()
	return p
}

func NewHttpProxy(addr string, auth *Auth) *Proxy {
	return NewProxy("http", addr, auth)
}

func NewSocks5Proxy(addr string, auth *Auth) *Proxy {
	return NewProxy("socks5", addr, auth)
}

// ProxyFromEnvironment read proxy form env for every request
// http.ProxyFromEnvironment read only once
func ProxyFromEnvironment(req *http.Request) (*url.URL, error) {
	return httpproxy.FromEnvironment().ProxyFunc()(req.URL)
}
