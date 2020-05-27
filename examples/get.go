package main

import (
	"fmt"

	"github.com/riete/requests"
)

func main() {
	resp, err := requests.Get("https://xxx")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Content(), resp.StatusCode(), resp.Status())

	data := map[string]string{
		"a": "1",
		"b": "2",
	}
	resp, err = requests.GetWithParams("http://xxxx", data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Content(), resp.StatusCode(), resp.Status())
}
