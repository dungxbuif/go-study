package main

import (
	"fmt"
	"time"
)

// select: like a switch but for channel operations
// Blocks until one case is ready; picks randomly if multiple are ready

// --- 1. Timeout with select ---

func slowOp() <-chan string {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- "result"
	}()
	return ch
}

// --- 2. Fan-in: merge multiple channels into one ---

func fanIn(ch1, ch2 <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for {
			select {
			case v, ok := <-ch1:
				if !ok {
					ch1 = nil
				} else {
					out <- v
				}
			case v, ok := <-ch2:
				if !ok {
					ch2 = nil
				} else {
					out <- v
				}
			}
			if ch1 == nil && ch2 == nil {
				close(out)
				return
			}
		}
	}()
	return out
}

func source(name string, msgs []string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, m := range msgs {
			time.Sleep(50 * time.Millisecond)
			ch <- fmt.Sprintf("[%s] %s", name, m)
		}
	}()
	return ch
}

// --- 3. Non-blocking send/receive with default ---

func tryReceive(ch <-chan int) {
	select {
	case v := <-ch:
		fmt.Println("received:", v)
	default:
		fmt.Println("no value ready, moving on")
	}
}

func main() {
	fmt.Println("=== 1. Timeout ===")
	result := slowOp()
	select {
	case v := <-result:
		fmt.Println("got:", v)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("timed out!")
	}

	// with enough timeout
	result2 := slowOp()
	select {
	case v := <-result2:
		fmt.Println("got:", v)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("timed out!")
	}

	fmt.Println("\n=== 2. Fan-in ===")
	a := source("A", []string{"one", "two"})
	b := source("B", []string{"alpha", "beta", "gamma"})
	merged := fanIn(a, b)
	for msg := range merged {
		fmt.Println(msg)
	}

	fmt.Println("\n=== 3. Non-blocking ===")
	ch := make(chan int, 1)
	tryReceive(ch)   // default fires — channel empty
	ch <- 42
	tryReceive(ch)   // receives 42
}
