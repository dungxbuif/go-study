package main

import "fmt"

type Obj struct {
	SomeProperties map[string]any
}

func NewObj() *Obj {
	return &Obj{
		SomeProperties: map[string]any{
			"foo": "bar",
		},
	}
}

func (o *Obj) SomeMethod(prop string) any {
	return o.SomeProperties[prop]
}

func main() {
	obj := NewObj()

	item := obj.SomeProperties["foo"]
	fmt.Println(item)

	item = obj.SomeMethod("foo")
	fmt.Println(item)
}
