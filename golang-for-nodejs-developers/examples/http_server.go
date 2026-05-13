package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello world"))
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("HTTP server listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
