package main

import (
	"fmt"
	"os"
)

// defer: schedule a function call to run when the surrounding function returns
//   — runs even on panic, in LIFO order
// panic: aborts normal execution (like an unchecked exception)
// recover: catches a panic inside a deferred function — must be called directly in a defer

// --- 1. Defer for cleanup (resource management) ---

func readFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close() // guaranteed to run regardless of what happens next

	fmt.Println("reading file:", f.Name())
	return nil
}

// --- 2. Defer LIFO order ---

func deferOrder() {
	fmt.Println("start")
	defer fmt.Println("first defer")
	defer fmt.Println("second defer")
	defer fmt.Println("third defer")
	fmt.Println("end")
}

// --- 3. Panic + Recover — use sparingly, mainly for boundary protection ---

func safeDiv(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	return a / b, nil // panics if b == 0
}

// --- 4. Wrapping a library that may panic ---

func safeCall(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	fn()
	return nil
}

func main() {
	fmt.Println("=== 1. Defer Cleanup ===")
	if err := readFile("/tmp/test.txt"); err != nil {
		fmt.Println("expected error:", err)
	}

	fmt.Println("\n=== 2. Defer LIFO ===")
	deferOrder()

	fmt.Println("\n=== 3. Panic + Recover ===")
	result, err := safeDiv(10, 2)
	fmt.Printf("10/2 = %d, err = %v\n", result, err)

	result, err = safeDiv(10, 0)
	fmt.Printf("10/0 = %d, err = %v\n", result, err)

	fmt.Println("\n=== 4. SafeCall ===")
	err = safeCall(func() {
		var s []int
		_ = s[0] // index out of range panic
	})
	fmt.Println("safeCall error:", err)
}
