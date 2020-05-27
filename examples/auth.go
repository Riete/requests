package main

import (
	"fmt"

	"github.com/riete/requests"
)

func main() {
	s := requests.NewSession()
	auth := make(map[string]interface{})
	auth["username"] = "xxx"
	auth["password"] = "xxx"
	resp, _ := s.JsonAuth("http://xxx", auth)
	fmt.Println(resp.Content())
	resp, _ = s.Get("http://xxxx")
	fmt.Println(resp.Content())
}
