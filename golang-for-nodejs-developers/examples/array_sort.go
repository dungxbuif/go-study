package main

import (
	"cmp"
	"fmt"
	"slices"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	intList := []int{1, 3, 5, 9, 4, 2, 0}

	// asc
	slices.Sort(intList)
	fmt.Println(intList)
	// desc
	slices.Reverse(intList)
	fmt.Println(intList)

	stringList := []string{"a", "d", "z", "b", "c", "y"}

	// asc
	slices.Sort(stringList)
	fmt.Println(stringList)
	// desc
	slices.Reverse(stringList)
	fmt.Println(stringList)

	collection := []Person{
		{"Li L", 8},
		{"Json C", 3},
		{"Zack W", 15},
		{"Yi M", 2},
	}

	// asc
	slices.SortFunc(collection, func(a, b Person) int {
		return cmp.Compare(a.Age, b.Age)
	})
	fmt.Println(collection)

	// desc
	slices.SortFunc(collection, func(a, b Person) int {
		return cmp.Compare(b.Age, a.Age)
	})
	fmt.Println(collection)
}
