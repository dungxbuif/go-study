package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Go sync primitives:
//   Mutex       — mutual exclusion lock
//   RWMutex     — multiple readers OR one writer
//   WaitGroup   — wait for a group of goroutines
//   Once        — run exactly once (singleton init)
//   atomic      — lock-free integer operations

// --- 1. Mutex — protect shared state ---

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v[key]++
}

func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v[key]
}

// --- 2. RWMutex — optimize read-heavy workloads ---

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func (c *Cache) Set(k, v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[k] = v
}

func (c *Cache) Get(k string) (string, bool) {
	c.mu.RLock() // multiple goroutines can read simultaneously
	defer c.mu.RUnlock()
	v, ok := c.data[k]
	return v, ok
}

// --- 3. sync.Once — safe singleton initialization ---

type config struct{ dsn string }

var (
	cfg  *config
	once sync.Once
)

func getConfig() *config {
	once.Do(func() {
		fmt.Println("initializing config (runs once)")
		cfg = &config{dsn: "postgres://localhost/mydb"}
	})
	return cfg
}

// --- 4. atomic — simple counter without mutex ---

type AtomicCounter struct {
	val int64
}

func (c *AtomicCounter) Inc()        { atomic.AddInt64(&c.val, 1) }
func (c *AtomicCounter) Load() int64 { return atomic.LoadInt64(&c.val) }

func main() {
	fmt.Println("=== 1. Mutex ===")
	counter := &SafeCounter{v: make(map[string]int)}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc("key")
		}()
	}
	wg.Wait()
	fmt.Println("count:", counter.Value("key")) // always 100

	fmt.Println("\n=== 2. RWMutex Cache ===")
	cache := &Cache{data: make(map[string]string)}
	cache.Set("lang", "Go")
	v, _ := cache.Get("lang")
	fmt.Println("cache get:", v)

	fmt.Println("\n=== 3. sync.Once ===")
	c1 := getConfig()
	c2 := getConfig() // "initializing" only prints once
	fmt.Println("same instance:", c1 == c2)

	fmt.Println("\n=== 4. atomic ===")
	ac := &AtomicCounter{}
	var wg2 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			ac.Inc()
		}()
	}
	wg2.Wait()
	fmt.Println("atomic count:", ac.Load()) // always 1000
}
