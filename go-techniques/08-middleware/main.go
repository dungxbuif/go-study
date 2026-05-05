package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Middleware pattern in Go: wrap http.Handler with another http.Handler
// Each middleware does its work then calls next.ServeHTTP(w, r)

type Middleware func(http.Handler) http.Handler

// Chain applies middlewares left-to-right
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// --- Logging middleware ---

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("[LOG] %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		fmt.Printf("[LOG] done in %s\n", time.Since(start))
	})
}

// --- Auth middleware ---

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer secret" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// --- RequestID middleware — inject value via context ---

type ctxKey string

const reqIDKey ctxKey = "requestID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := fmt.Sprintf("req-%d", time.Now().UnixNano())
		ctx := context.WithValue(r.Context(), reqIDKey, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// --- Business handler ---

func userHandler(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(reqIDKey)
	fmt.Fprintf(w, "Hello from userHandler (reqID: %v)\n", reqID)
}

func main() {
	mux := http.NewServeMux()

	// Apply middleware chain: RequestID → Logger → Auth → handler
	mux.Handle("/api/users", Chain(
		http.HandlerFunc(userHandler),
		RequestID,
		Logger,
		Auth,
	))

	fmt.Println("Server listening on :8080")
	fmt.Println("Try: curl -H 'Authorization: Bearer secret' http://localhost:8080/api/users")
	http.ListenAndServe(":8080", mux)
}
