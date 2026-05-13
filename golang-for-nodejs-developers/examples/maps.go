package main

import (
	"fmt"
	"maps"
)

func main() {
	map1 := make(map[string]string)

	map1["foo"] = "bar"

	item, found := map1["foo"]
	fmt.Println(found)
	fmt.Println(item)

	delete(map1, "foo")

	item, found = map1["foo"]
	fmt.Println(found)
	fmt.Println(item)

	map2 := map[string]int{
		"foo": 100,
		"bar": 200,
		"baz": 300,
	}

	// Use maps package to clone
	map3 := maps.Clone(map2)

	for key, value := range map3 {
		fmt.Println(key, value)
	}

	// Use maps package to delete by condition
	maps.DeleteFunc(map3, func(k string, v int) bool {
		return v > 250
	})

	fmt.Println("After DeleteFunc:", map3)
}
