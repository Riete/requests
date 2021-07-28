## Requests

Go HTTP Client Library

## Usage
```
// request
url := "http://x.x.x.x"
r := NewRequest(DefaultConfig)
r.Get(url, nil)
r.Get(url, map[string]string{"a": "1", "b": "2"})
r.Post(url, map[string]interface{}{"a": "1", "b": "2"})
r.PostForm(url, map[string]string{"a": "1", "b": "2"})
r.Put(url, map[string]interface{}{"a": "1", "b": "2"})
r.Delete(url)

// session
loginUrl := "http://x.x.x.x/login"
s := NewSession(DefaultConfig)
s.Post(loginUrl, map[string]interface{}{"user": "username", "password": "password"})
s.Get(url, nil)
s.Get(url, map[string]string{"a": "1", "b": "2"})
s.Post(url, map[string]interface{}{"a": "1", "b": "2"})
s.PostForm(url, map[string]string{"a": "1", "b": "2"})
s.Put(url, map[string]interface{}{"a": "1", "b": "2"})
s.Delete(url)
```