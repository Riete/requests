package main

import (
	"fmt"
	"time"

	"github.com/riete/requests"
)

func main() {
	s := requests.NewSession()
	s.SetBasicAuth("aaa", "bbb")
	data := make(map[string]interface{})
	data["a"] = "1"
	data["b"] = "2"

	data2 := map[string]string{"aa": "bb"}

	resp, err := s.PostJson("http://xxx", data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Content(), resp.StatusCode(), resp.Status())

	resp, err = s.PostForm("http://xxx", data2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Content(), resp.StatusCode(), resp.Status())

	s.SetBearerTokenAuth("xxx")
	resp, err = s.PostJson("http://xxx", data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Content(), resp.StatusCode(), resp.Status())

	resp, err = s.PostForm("http://xxx", data2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Content(), resp.StatusCode(), resp.Status())

	resp, err = s.JsonAuth("http://xxx", data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Content(), resp.StatusCode(), resp.Status())

	resp, err = s.FormAuth("http://xxx", data2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Content(), resp.StatusCode(), resp.Status())
	s.SetTimeout(1 * time.Microsecond)
	s.SkipTLSVerify()
	proxy := map[string]string{
		"http_proxy":  "http://xxx",
		"https_proxy": "http://xxxx",
	}
	s.SetProxy(proxy)
	resp, err = s.Get("https://xxx")
	fmt.Println(resp)
	fmt.Println(err)
}
