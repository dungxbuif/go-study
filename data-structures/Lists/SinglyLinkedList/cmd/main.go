package main

import (
	"fmt"
	"singlylinkedlist"
)

func main() {
	l := singlylinkedlist.New()

	fmt.Println("--- Thêm phần tử ---")
	l.PushBack("Node 2")
	l.PushBack("Node 3")
	l.PushFront("Node 1")
	printList(l)

	fmt.Println("\n--- Di chuyển & Xóa ---")
	mark := l.Front().Next() // Node 2
	l.InsertAfter("Node 2.5", mark)
	l.MoveToFront(l.Back()) // Move Node 3 to front
	l.Remove(mark)         // Remove Node 2
	printList(l)

	fmt.Printf("\nTổng số phần tử: %d\n", l.Len())
}

func printList(l *singlylinkedlist.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v", e.Value)
		if e.Next() != nil {
			fmt.Print(" -> ")
		}
	}
	fmt.Println()
}
