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
r.PostFrom(url, map[string]string{})
r.Put(url, map[string]interface{}{})
r.Delete(url, map[string]interface{}{})

// download
r.DownloadToWriter(originUrl, io.Writer)
r.Download(filePath, originUrl, rate)

// upload
r.Upload(originUrl, map[string]string{}, rate, filepath1, filepath2 ...)

// response
r.Content()
r.ContentToString()
r.Status() // status code, status
r.Resp() // raw http.Response
```