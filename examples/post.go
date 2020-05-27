package main

import (
	"fmt"

	"github.com/riete/requests"
)

func main() {
	resp, err := requests.Post("http://xxxx")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Content(), resp.StatusCode(), resp.Status())

	data := make(map[string]interface{})
	data["a"] = "1"
	data["b"] = "2"
	resp, err = requests.PostJson("http://xxxx", data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Content(), resp.StatusCode(), resp.Status())

	data2 := map[string]string{
		"a": "1",
		"b": "2",
	}
	resp, err = requests.PostForm("http://xxxx", data2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Content(), resp.StatusCode(), resp.Status())
}
