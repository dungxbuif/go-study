package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprint(os.Stderr, "hello error\n")
}
