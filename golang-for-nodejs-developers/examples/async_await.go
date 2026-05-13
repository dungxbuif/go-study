package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func hello(name string) chan any {
	ch := make(chan any, 1)
	go func() {
		time.Sleep(1 * time.Second)
		if name == "fail" {
			ch <- errors.New("failed")
		} else {
			ch <- "hello " + name
		}
		close(ch)
	}()

	return ch
}

func main() {
	result := <-hello("bob")
	switch v := result.(type) {
	case string:
		fmt.Println(v)
	case error:
		log.Println(v)
	}

	result = <-hello("fail")
	switch v := result.(type) {
	case string:
		fmt.Println(v)
	case error:
		log.Println(v.Error())
	}
}
