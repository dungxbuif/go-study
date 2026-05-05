package main

import (
	"doublylinkedlist"
	"fmt"
)

func main() {
	l := doublylinkedlist.New()

	// 1. Thêm phần tử
	fmt.Println("--- Thêm phần tử ---")
	l.PushBack("Node 2")
	l.PushFront("Node 1")
	l.PushBack("Node 3")
	printForward(l) // 1 -> 2 -> 3

	// 2. Thử nghiệm InsertBefore
	fmt.Println("\n--- InsertBefore ---")
	mark := l.Back() // Node 3
	l.InsertBefore("Node 2.5", mark)
	printForward(l) // 1 -> 2 -> 2.5 -> 3

	// 3. Duyệt ngược (Backward Traversal)
	fmt.Println("\n--- Duyệt ngược (Backward) ---")
	for e := l.Back(); e != nil; e = e.Prev() {
		fmt.Printf("%v", e.Value)
		if e.Prev() != nil {
			fmt.Print(" <- ")
		}
	}
	fmt.Println()

	// 4. MoveToFront & Remove
	fmt.Println("\n--- Move & Remove ---")
	l.MoveToFront(mark) // Move Node 3 to front
	l.Remove(l.Back().Prev()) // Remove Node 2.5 (vì Node 3 đã đi, Node 2.5 là áp chót)
	printForward(l)

	fmt.Printf("\nTổng số phần tử: %d\n", l.Len())
}

func printForward(l *doublylinkedlist.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v", e.Value)
		if e.Next() != nil {
			fmt.Print(" -> ")
		}
	}
	fmt.Println()
}
