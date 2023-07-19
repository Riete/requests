## Requests

Go HTTP Client Library

## Usage
```
// New Request
r := NewRequest(WithTimout(time.Second))

// New Session
r := NewSession(WithTimout(time.Second))

// do http method
r.Get(url, map[string]string{})
r.Post(url, map[string]interface{}{})
r.Post(url, map[string]interface{}{})
r.PostFrom(url, map[string]string{})
r.Put(url, map[string]interface{}{})
r.Delete(url, map[string]interface{}{})

// download
r.Download(filePath, originUrl)
r.DownloadWithRateLimit(filePath, originUrl, rate)

// response
r.Content()
r.ContentToString()
r.Status
r.StatusCode
r.Resp // raw http.Response
```