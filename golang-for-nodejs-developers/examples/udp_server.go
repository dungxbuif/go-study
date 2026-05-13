package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 3000,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	fmt.Printf("server listening %s\n", conn.LocalAddr().String())

	for {
		message := make([]byte, 1024)
		rlen, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			continue
		}

		data := strings.TrimSpace(string(message[:rlen]))
		fmt.Printf("received: %s from %s\n", data, remote)
	}
}
