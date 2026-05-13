package main

import (
	"fmt"
	"iter"
)

// Generator returns a Go 1.23 iterator (iter.Seq)
func Generator() iter.Seq[string] {
	return func(yield func(string) bool) {
		if !yield("hello") {
			return
		}
		if !yield("world") {
			return
		}
	}
}

func main() {
	// Manual iteration using iter.Pull to match original behavior of gen.next()
	next, stop := iter.Pull(Generator())
	defer stop()

	for {
		value, ok := next()
		fmt.Println(value, ok)
		if !ok {
			break
		}
	}

	// Idiomatic range over iterator
	for value := range Generator() {
		fmt.Println(value)
	}

	// Repeating manual iteration to match the original example structure
	next2, stop2 := iter.Pull(Generator())
	defer stop2()
	for {
		value, ok := next2()
		fmt.Println(value, ok)
		if !ok {
			break
		}
	}
}
