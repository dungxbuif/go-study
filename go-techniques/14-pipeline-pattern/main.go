package main

import (
	"fmt"
	"sync"
)

// Pipeline: chain of stages connected by channels
// Each stage: receives from upstream, transforms, sends to downstream
// Enables composable, concurrent data processing

// --- Stage helpers ---

// generator: seed values into the pipeline
func generator(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

// sq: squares each number
func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

// filter: passes only values satisfying predicate
func filter(done <-chan struct{}, in <-chan int, pred func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if pred(n) {
				select {
				case out <- n:
				case <-done:
					return
				}
			}
		}
	}()
	return out
}

// --- Fan-out / Fan-in for parallel work ---

func fanOut(done <-chan struct{}, in <-chan int, n int) []<-chan int {
	channels := make([]<-chan int, n)
	for i := 0; i < n; i++ {
		channels[i] = sq(done, in)
	}
	return channels
}

func merge(done <-chan struct{}, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(ch <-chan int) {
		defer wg.Done()
		for n := range ch {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)

	fmt.Println("=== 1. Linear Pipeline ===")
	// generator → sq → filter (keep only > 10) → print
	nums := generator(done, 1, 2, 3, 4, 5, 6)
	squares := sq(done, nums)
	bigSquares := filter(done, squares, func(n int) bool { return n > 10 })

	for v := range bigSquares {
		fmt.Println(v)
	}

	fmt.Println("\n=== 2. Fan-out / Fan-in ===")
	// Split work across 3 workers then merge results
	done2 := make(chan struct{})
	defer close(done2)

	src := generator(done2, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	workers := fanOut(done2, src, 3)
	results := merge(done2, workers...)

	var out []int
	for v := range results {
		out = append(out, v)
	}
	fmt.Println("merged results:", out)
}
