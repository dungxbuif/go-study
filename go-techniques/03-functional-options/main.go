package main

import (
	"fmt"
	"time"
)

// Functional Options Pattern
// Solves the problem of structs with many optional config fields
// without requiring a separate Config struct or many constructor variants

type Server struct {
	host    string
	port    int
	timeout time.Duration
	maxConn int
	tls     bool
}

// ServerOption is a function that mutates a Server
type ServerOption func(*Server)

// Each option is a constructor function that returns a ServerOption
func WithHost(host string) ServerOption {
	return func(s *Server) { s.host = host }
}

func WithPort(port int) ServerOption {
	return func(s *Server) { s.port = port }
}

func WithTimeout(d time.Duration) ServerOption {
	return func(s *Server) { s.timeout = d }
}

func WithMaxConn(n int) ServerOption {
	return func(s *Server) { s.maxConn = n }
}

func WithTLS() ServerOption {
	return func(s *Server) { s.tls = true }
}

// NewServer applies defaults first, then overrides with provided options
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		host:    "localhost",
		port:    8080,
		timeout: 30 * time.Second,
		maxConn: 100,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) String() string {
	return fmt.Sprintf("Server{host:%s port:%d timeout:%s maxConn:%d tls:%v}",
		s.host, s.port, s.timeout, s.maxConn, s.tls)
}

func main() {
	// Use only defaults
	s1 := NewServer()
	fmt.Println("default:", s1)

	// Override only what you need
	s2 := NewServer(
		WithPort(9090),
		WithTimeout(10*time.Second),
		WithTLS(),
	)
	fmt.Println("custom: ", s2)

	// Production config
	s3 := NewServer(
		WithHost("0.0.0.0"),
		WithPort(443),
		WithMaxConn(1000),
		WithTLS(),
	)
	fmt.Println("prod:   ", s3)
}
