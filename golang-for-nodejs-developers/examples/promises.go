package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func asyncMethod(value string) chan any {
	ch := make(chan any, 1)
	go func() {
		time.Sleep(1 * time.Second)
		ch <- "resolved: " + value
		close(ch)
	}()

	return ch
}

func resolveAll(ch ...chan any) chan any {
	var wg sync.WaitGroup
	res := make([]string, len(ch))
	resCh := make(chan any, 1)

	go func() {
		for i, c := range ch {
			wg.Add(1)
			go func(j int, ifcCh chan any) {
				ifc := <-ifcCh
				switch v := ifc.(type) {
				case error:
					resCh <- v
				case string:
					res[j] = v
				}
				wg.Done()
			}(i, c)
		}

		wg.Wait()
		resCh <- res
		close(resCh)
	}()

	return resCh
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		result := <-asyncMethod("foo")
		switch v := result.(type) {
		case string:
			fmt.Println(v)
		case error:
			log.Println(v)
		}

		wg.Done()
	}()

	go func() {
		result := <-resolveAll(
			asyncMethod("A"),
			asyncMethod("B"),
			asyncMethod("C"),
		)

		switch v := result.(type) {
		case []string:
			fmt.Println(v)
		case error:
			log.Println(v)
		}

		wg.Done()
	}()

	wg.Wait()
}
