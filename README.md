## Requests

Go HTTP Client Library

## Usage

### New Request/Session or (Pool)
```
NewRequest(options ...RequestOption)
NewSession(options ...RequestOption)

// Pool
NewRequestPool(options ...RequestOption)
NewSessionPool(options ...RequestOption)
```

### RequestOption
* WithTimeout(t time.Duration)
* WithHeader(headers ...map[string]string)
* WithProxyEnv(proxy map[string]string)
* WithProxyURL(proxy *url.URL)
* WithProxyFunc(f func(*http.Request) (*url.URL, error))
* WithUnsetProxy()
* WithSkipTLS()
* WithBasicAuth(username, password string)
* WithBearerTokenAuth(token string)
* WithTransport(tr http.RoundTripper)
* WithDefaultTransport()
* WithClient(client *http.Client)
* WithDefaultClient()

### DoHttpMethod
```
r := NewRequest(RequestOptions...)
r.SetXXX() // if needed
r.Get(originURL string, options ...MethodOption)
r.Post(originURL string, options ...MethodOption)
r.Put(originURL string, options ...MethodOption)
r.Delete(originURL string, options ...MethodOption)
r.CloseIdleConnections() // if needed
```

### Upload/Download
```
// rate is speed per second, e.g. 1024 ==> 1KiB, if rate <= 0 it means no limit
r.DownloadToWriter(originURL string, w io.Writer)
r.Download(filePath, originURL string, rate int64)
r.Upload(originURL string, data map[string]string, rate int64, filePaths ...string) 
```

### Response
```
r.Status() // status code and status
r.Content() // []byte data
r.ContentToString() // string data
```

### PoolMode
```
rp := NewRequestPool(RequestOptions...)
r := p.Get()
defer rp.Put(r)
// Do Http Method
...
```

### MethodOption
* WithParams(params map[string]string)
* WithJsonData(data map[string]any) MethodOption
* WithFormData(data map[string]string)

### Proxy
```
NewProxy(scheme, addr string, auth *Auth)
NewHttpProxy(addr string, auth *Auth)
NewSocks5Proxy(addr string, auth *Auth)
ProxyFromEnvironment(req *http.Request)
```