## Requests

Go HTTP Client Library

## Usage
```
// Default Request
url := "http://x.x.x.x"
requests.Get(url, nil)
requests.Get(url, map[string]string{"a": "1", "b": "2"})
requests.Post(url, map[string]interface{}{"a": "1", "b": "2"})
requests.PostForm(url, map[string]string{"a": "1", "b": "2"})
requests.Put(url, map[string]interface{}{"a": "1", "b": "2"})
requests.Delete(url)

// Default Session
loginUrl := "http://x.x.x.x/login"
s := requests.Session()
s.Post(loginUrl, map[string]interface{}{"user": "username", "password": "password"})
s.Get(url, nil)
s.Get(url, map[string]string{"a": "1", "b": "2"})
s.Post(url, map[string]interface{}{"a": "1", "b": "2"})
s.PostForm(url, map[string]string{"a": "1", "b": "2"})
s.Put(url, map[string]interface{}{"a": "1", "b": "2"})
s.Delete(url)

// New Request
r := NewRequest(DefaultConfig)
// New Session
s := NewSession(DefaultConfig)
```