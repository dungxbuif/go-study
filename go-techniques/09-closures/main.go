package main

import "fmt"

// Closures: functions that capture variables from their enclosing scope
// In Go, functions are first-class values — assign, pass, return them

// --- 1. Counter — closure captures state ---

func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// --- 2. Memoize — closure captures a cache ---

func memoize(fn func(int) int) func(int) int {
	cache := make(map[int]int)
	return func(n int) int {
		if v, ok := cache[n]; ok {
			fmt.Printf("cache hit for %d\n", n)
			return v
		}
		result := fn(n)
		cache[n] = result
		return result
	}
}

func slowFib(n int) int {
	if n <= 1 {
		return n
	}
	return slowFib(n-1) + slowFib(n-2)
}

// --- 3. Middleware / decorator via closure ---

func withLogging(name string, fn func(int) int) func(int) int {
	return func(n int) int {
		fmt.Printf("calling %s(%d)\n", name, n)
		result := fn(n)
		fmt.Printf("%s(%d) = %d\n", name, n, result)
		return result
	}
}

// --- 4. Closure in loops — common gotcha ---

func main() {
	fmt.Println("=== 1. Counter ===")
	c1 := makeCounter()
	c2 := makeCounter() // independent state
	fmt.Println(c1(), c1(), c1()) // 1 2 3
	fmt.Println(c2(), c2())       // 1 2 — c2 has its own count

	fmt.Println("\n=== 2. Memoize ===")
	fastFib := memoize(slowFib)
	fmt.Println(fastFib(10))
	fmt.Println(fastFib(10)) // cache hit

	fmt.Println("\n=== 3. Decorator ===")
	loggedDouble := withLogging("double", func(n int) int { return n * 2 })
	loggedDouble(5)

	fmt.Println("\n=== 4. Closure in Loop — correct way ===")
	funcs := make([]func(), 3)
	for i := 0; i < 3; i++ {
		i := i // shadow i — each closure captures its own copy
		funcs[i] = func() { fmt.Println(i) }
	}
	for _, f := range funcs {
		f() // prints 0, 1, 2 — not 3, 3, 3
	}
}
