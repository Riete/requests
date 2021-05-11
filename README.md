## requests

Go HTTP library

## requests
```
package main

import (
    "fmt"
	
    "github.com/riete/requests"
)

// http get http://example.com
resp, _ := requests.Get("http://example.com")
fmt.Println(resp.Content())

// http get http://example.com?key1=value1&key2=value2
data := map[string]string{"key1": "value1", "key2": "value2"}
resp, _ := requests.GetWithParams("http://example.com", data)
fmt.Println(resp.Content())

// http post http://example.com
resp, _ := requests.Post("http://example.com")
fmt.Println(resp.Content())

// post json
data := make(map[string]interface{})
data["a"] = "1"
data["b"] = "2"
resp, _ := requests.PostJson("http://example.com", data)
fmt.Println(resp.Content())


// post form
data := map[string]string{"key1": "value1", "key2": "value2"}
resp, err = requests.PostForm("http://example.com", data)
fmt.Println(resp.Content())

```

## session
```
package main

import (
    "fmt"
    "time"
	
    "github.com/riete/requests"
)

// new session
s := requests.NewSession()

// set timeout to 5 second
s.SetTimeout(5 * time.Second)

// set skip tls verify
s.SkipTLSVerify()

// set proxy
proxy := map[string]string{
    "http_proxy":  "http://xxx",
    "https_proxy": "http://xxxx",
    }
s.SetProxy(proxy)

// set basic auth
s.SetBasicAuth("username", "password")

// set bearer token auth
s.SetBearerTokenAuth("token")

// login via api interface
// post json
auth := make(map[string]interface{})
auth["username"] = "xxx"
auth["password"] = "xxx"
s.JsonAuth("http://example.com", auth)

// post form
auth := map[string]string{"username": "xxx", "password": "xxx"}
s.FormAuth("http://example.com", auth)

// get
s.Get("http://example.com")

data := map[string]string{"key1": "value1", "key2": "value2"}
s.GetWithParams("http://example.com", data)

// post
s.Post("http://example.com")

data := make(map[string]interface{})
data["a"] = "1"
data["b"] = "2"
s.PostJson("http://example.com", data)

data := map[string]string{"key1": "value1", "key2": "value2"}
s.PostForm("http://example.com", data)
```