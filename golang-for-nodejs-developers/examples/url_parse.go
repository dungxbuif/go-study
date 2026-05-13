package main

import (
	"fmt"
	"net/url"
)

func main() {
	urlstr := "http://bob:secret@sub.example.com:8080/somepath?foo=bar"

	u, err := url.Parse(urlstr)
	if err != nil {
		panic(err)
	}

	fmt.Println(u.Scheme + ":") // To match JS output if protocol includes colon
	fmt.Println(u.User)
	fmt.Println(u.Port())
	fmt.Println(u.Hostname())
	fmt.Println(u.Path)
	fmt.Println(u.Query())
}
