package requests

import (
	"net/http"
	"net/url"

	"golang.org/x/net/http/httpproxy"
)

type ProxyScheme string

func (p ProxyScheme) String() string {
	return string(p)
}

const (
	HttpProxy   ProxyScheme = "http"
	Socks5Proxy ProxyScheme = "socks5"
)

type Auth struct {
	Username string
	Password string
}

// Proxy addr is host:port
type Proxy struct {
	addr     string
	scheme   ProxyScheme
	auth     *Auth
	proxyURL string
}

func (p *Proxy) buildProxyURL() {
	if p.auth != nil {
		credential := url.QueryEscape(p.auth.Username) + ":" + url.QueryEscape(p.auth.Password)
		p.proxyURL = p.scheme.String() + "://" + credential + "@" + p.addr
	} else {
		p.proxyURL = p.scheme.String() + "://" + p.addr
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

func newProxy(scheme ProxyScheme, addr string, auth *Auth) *Proxy {
	p := &Proxy{addr: addr, scheme: scheme, auth: auth}
	p.buildProxyURL()
	return p
}

func NewHttpProxy(addr string, auth *Auth) *Proxy {
	return newProxy(HttpProxy, addr, auth)
}

func NewSocks5Proxy(addr string, auth *Auth) *Proxy {
	return newProxy(Socks5Proxy, addr, auth)
}

// ProxyFromEnvironment read proxy form env for every request
// http.ProxyFromEnvironment read only once
func ProxyFromEnvironment(req *http.Request) (*url.URL, error) {
	return httpproxy.FromEnvironment().ProxyFunc()(req.URL)
}
