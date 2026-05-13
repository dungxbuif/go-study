package main

import (
	"bufio"
	"fmt"
	"net"
)

func handler(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		conn.Write([]byte("Received: "))
		conn.Write([]byte(message))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	fmt.Println("TCP server listening on port 3000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handler(conn)
	}
}
