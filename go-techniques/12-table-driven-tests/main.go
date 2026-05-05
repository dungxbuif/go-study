package main

import "fmt"

// Table-driven tests are the idiomatic Go testing convention.
// Define a slice of test cases (table), then range over them.
// This file shows the functions under test; see main_test.go for tests.

func Add(a, b int) int { return a + b }

func FizzBuzz(n int) string {
	switch {
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return fmt.Sprintf("%d", n)
	}
}

func IsPalindrome(s string) bool {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(FizzBuzz(15)) // FizzBuzz
	fmt.Println(FizzBuzz(9))  // Fizz
	fmt.Println(FizzBuzz(10)) // Buzz
	fmt.Println(FizzBuzz(7))  // 7

	fmt.Println(IsPalindrome("racecar")) // true
	fmt.Println(IsPalindrome("hello"))   // false
}
