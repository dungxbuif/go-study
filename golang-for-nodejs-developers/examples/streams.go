package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	inStream := new(bytes.Buffer)
	w := bufio.NewWriter(inStream)
	_, err := w.Write([]byte("foo"))
	if err != nil {
		panic(err)
	}
	_, err = w.Write([]byte("bar"))
	if err != nil {
		panic(err)
	}
	err = w.Flush()
	if err != nil {
		panic(err)
	}

	inStream.WriteTo(os.Stdout)
	fmt.Println()

	outStream := new(bytes.Buffer)
	outStream.Write([]byte("abc\n"))
	outStream.Write([]byte("xyz\n"))
	piper, pipew := io.Pipe()

	go func() {
		sc := bufio.NewScanner(piper)
		for sc.Scan() {
			fmt.Println("received: " + sc.Text())
		}
		if err := sc.Err(); err != nil {
			panic(err)
		}
	}()

	go func() {
		defer pipew.Close()
		io.Copy(pipew, outStream)
	}()

	// Give it a bit of time to complete
	time.Sleep(100 * time.Millisecond)
}
