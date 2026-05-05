package main

import (
	"fmt"
	"sync"
	"time"
)

// Goroutines: lightweight threads managed by Go runtime
// Channels: typed conduits for communication between goroutines
// Principle: "Do not communicate by sharing memory; share memory by communicating"

// --- 1. Basic goroutine ---

func greet(name string) {
	fmt.Println("Hello,", name)
}

// --- 2. Channel — send and receive ---

func produce(ch chan<- int, n int) {
	for i := 0; i < n; i++ {
		ch <- i
		time.Sleep(50 * time.Millisecond)
	}
	close(ch) // signal no more values
}

// --- 3. Fan-out: one producer, many workers ---

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("worker %d processing job %d\n", id, j)
	}
}

// --- 4. Done channel pattern (manual cancel) ---

func generate(done <-chan struct{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; ; i++ {
			select {
			case out <- i:
			case <-done:
				fmt.Println("generator cancelled")
				return
			}
		}
	}()
	return out
}

func main() {
	fmt.Println("=== 1. Basic Goroutine ===")
	var wg sync.WaitGroup
	for _, name := range []string{"Alice", "Bob", "Carol"} {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			greet(n)
		}(name)
	}
	wg.Wait()

	fmt.Println("\n=== 2. Buffered Channel ===")
	ch := make(chan int, 5)
	go produce(ch, 5)
	for v := range ch {
		fmt.Println("received:", v)
	}

	fmt.Println("\n=== 3. Fan-out Worker Pool ===")
	jobs := make(chan int, 10)
	var wg2 sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg2.Add(1)
		go worker(w, jobs, &wg2)
	}
	for j := 0; j < 9; j++ {
		jobs <- j
	}
	close(jobs)
	wg2.Wait()

	fmt.Println("\n=== 4. Done Channel Cancel ===")
	done := make(chan struct{})
	nums := generate(done)
	for i := 0; i < 5; i++ {
		fmt.Println("got:", <-nums)
	}
	close(done)
	time.Sleep(10 * time.Millisecond)
}
