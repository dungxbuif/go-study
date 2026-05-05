package main

import "fmt"

// Generics (Go 1.18+): write functions and types that work with any type
// Use type parameters [T constraint] to express this

// --- 1. Generic function ---

func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func Filter[T any](slice []T, fn func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U {
	acc := initial
	for _, v := range slice {
		acc = fn(acc, v)
	}
	return acc
}

// --- 2. Type constraint ---

type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

func Sum[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

func Min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// --- 3. Generic data structure ---

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T)       { s.items = append(s.items, v) }
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top, true
}
func (s *Stack[T]) Len() int { return len(s.items) }

func main() {
	fmt.Println("=== 1. Map / Filter / Reduce ===")
	nums := []int{1, 2, 3, 4, 5}

	doubled := Map(nums, func(n int) int { return n * 2 })
	fmt.Println("doubled:", doubled)

	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("evens:", evens)

	sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	fmt.Println("sum:", sum)

	words := []string{"go", "is", "great"}
	lengths := Map(words, func(s string) int { return len(s) })
	fmt.Println("lengths:", lengths)

	fmt.Println("\n=== 2. Numeric Constraint ===")
	fmt.Println("sum int:", Sum([]int{1, 2, 3}))
	fmt.Println("sum float:", Sum([]float64{1.1, 2.2, 3.3}))
	fmt.Println("min:", Min(3.14, 2.71))

	fmt.Println("\n=== 3. Generic Stack ===")
	s := &Stack[string]{}
	s.Push("a")
	s.Push("b")
	s.Push("c")
	for s.Len() > 0 {
		v, _ := s.Pop()
		fmt.Print(v, " ")
	}
	fmt.Println()
}
