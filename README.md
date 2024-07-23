## Requests

Go HTTP Client Library

## Usage

### New Request/Session or (Pool)
```
NewRequest(RequestOptions...)
NewSession(RequestOptions...)

// Pool
NewRequestPool(RequestOptions...)
NewSessionPool(RequestOptions...)
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
r.Get(url, MethodOptions...)
r.Post(url, MethodOptions...)
r.Put(url, MethodOptions...)
r.Delete(url, MethodOptions...)
r.CloseIdleConnections() // if needed
```

### Upload/Download
```
// rate is speed per second, e.g. 1024 ==> 1KiB, if rate <= 0 it means no limit
r.Download(filePath, url, rate)
r.Upload(url, data, rate, filePaths ...) 
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
NewHttpProxy(addr, auth)
NewSocks5Proxy(addr, auth)
ProxyFromEnvironment(req *http.Request)
```