package main

import (
	"fmt"
	"slices"
)

func main() {
	array := []int{1, 2, 3, 4, 5}
	fmt.Println(array)

	clone := slices.Clone(array)
	fmt.Println(clone)

	sub := array[2:4]
	fmt.Println(sub)

	concatenated := slices.Concat(array, []int{6, 7})
	fmt.Println(concatenated)

	prepended := slices.Concat([]int{-2, -1, 0}, concatenated)
	fmt.Println(prepended)
}
