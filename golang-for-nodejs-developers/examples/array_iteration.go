package main

import (
	"fmt"
	"strings"
)

func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

func Filter[T any](s []T, f func(T, int) bool) []T {
	var r []T
	for i, v := range s {
		if f(v, i) {
			r = append(r, v)
		}
	}
	return r
}

func Reduce[T, U any](s []T, f func(U, T, int) U, init U) U {
	r := init
	for i, v := range s {
		r = f(r, v, i)
	}
	return r
}

func main() {
	array := []string{"a", "b", "c"}

	for i, value := range array {
		fmt.Println(i, value)
	}

	mapped := Map(array, func(value string) string {
		return strings.ToUpper(value)
	})

	fmt.Println(mapped)

	filtered := Filter(array, func(value string, i int) bool {
		return i%2 == 0
	})

	fmt.Println(filtered)

	reduced := Reduce(array, func(acc []string, value string, i int) []string {
		if i%2 == 0 {
			acc = append(acc, strings.ToUpper(value))
		}
		return acc
	}, []string{})

	fmt.Println(reduced)
}
