package main

import (
	"context"
	"fmt"
	"time"
)

// context.Context: carries deadlines, cancellation signals, and request-scoped values
// Rule: always pass ctx as the FIRST argument to functions that do I/O or call other services

// --- 1. WithCancel — manual cancel ---

func streamData(ctx context.Context, ch chan<- int) {
	defer close(ch)
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("streamData cancelled:", ctx.Err())
			return
		case ch <- i:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// --- 2. WithTimeout — auto cancel after duration ---

func fetchUser(ctx context.Context, id int) (string, error) {
	done := make(chan string, 1)
	go func() {
		time.Sleep(200 * time.Millisecond) // simulate DB query
		done <- fmt.Sprintf("User#%d", id)
	}()

	select {
	case user := <-done:
		return user, nil
	case <-ctx.Done():
		return "", fmt.Errorf("fetchUser: %w", ctx.Err())
	}
}

// --- 3. WithValue — pass request-scoped data (e.g. request ID) ---

type contextKey string

const requestIDKey contextKey = "requestID"

func handleRequest(ctx context.Context) {
	reqID := ctx.Value(requestIDKey)
	fmt.Println("handling request:", reqID)
	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	reqID := ctx.Value(requestIDKey)
	fmt.Println("processing request:", reqID)
}

func main() {
	fmt.Println("=== 1. WithCancel ===")
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan int, 5)
	go streamData(ctx, ch)

	for i := 0; i < 4; i++ {
		fmt.Println("received:", <-ch)
	}
	cancel() // stop the goroutine
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n=== 2. WithTimeout ===")
	// fast enough
	ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel2()
	user, err := fetchUser(ctx2, 1)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("got:", user)
	}

	// too slow
	ctx3, cancel3 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel3()
	user, err = fetchUser(ctx3, 2)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("got:", user)
	}

	fmt.Println("\n=== 3. WithValue ===")
	ctx4 := context.WithValue(context.Background(), requestIDKey, "req-abc-123")
	handleRequest(ctx4)
}
