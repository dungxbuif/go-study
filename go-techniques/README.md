# 🛠️ 14 Kỹ Thuật Lập Trình Go Nâng Cao (Advanced Go Techniques)

Tài liệu này tổng hợp 14 kỹ thuật lập trình nâng cao và thực tiễn thiết kế hệ thống trong Go. Mỗi kỹ thuật đều đi kèm lời giải thích lý thuyết trực quan bằng tiếng Việt, mã nguồn giải mẫu hoàn chỉnh (cho phép sao chép/nhân bản trực tiếp) và liên kết đến các file chạy thử.

---

## 🗺️ Bản Đồ Wiki & Hệ Sinh Thái Học Tập (Wiki Navigation Hub)

|                  Trang Chủ                   |                 So Sánh Core                 |                              So Sánh Framework                               |               Kỹ Thuật Nâng Cao                |                  Lộ Trình Thực Hành                  |
| :------------------------------------------: | :------------------------------------------: | :--------------------------------------------------------------------------: | :--------------------------------------------: | :--------------------------------------------------: |
| 🏠 **[Trang Chủ (Wiki Root)](../README.md)** | 📊 **[Go vs Node.js Core](../GO_NODEJS.md)** | 🚀 **[Echo vs Gin vs NestJS vs Express](../framework-comparison/README.md)** | 🛠️ **[14 Kỹ Thuật Go Luyện Tập](./README.md)** | 🎯 **[20 Bài Tập Tự Luyện](../exercises/README.md)** |

---

## 📂 Danh Sách 14 Kỹ Thuật Go Nâng Cao

1. [01 - Interfaces (Triết lý Thiết Kế & DI)](#01---interfaces-triết-lý-thiết-kế--di)
2. [02 - Struct Embedding (Composition over Inheritance)](#02---struct-embedding-composition-over-inheritance)
3. [03 - Functional Options Pattern](#03---functional-options-pattern)
4. [04 - Error Handling (Wrapping, Is/As & Custom Errors)](#04---error-handling-wrapping-isas--custom-errors)
5. [05 - Goroutines & Channels (Đồng Thời Cơ Bản)](#05---goroutines--channels-đồng-thời-cơ-bản)
6. [06 - Select & Timeout (Đa Kênh Bất Đồng Bộ)](#06---select--timeout-đa-kênh-bất-đồng-bộ)
7. [07 - Context (Quản Lý Vòng Đời & Timeout)](#07---context-quản-lý-vòng-đời--timeout)
8. [08 - Middleware Pattern (HTTP Chain)](#08---middleware-pattern-http-chain)
9. [09 - Closures & State Capture](#09---closures--state-capture)
10.   [10 - Generics (Type Parameters & Constraints)](#10---generics-type-parameters--constraints)
11.   [11 - Defer, Panic & Recover (Xử Lý Lỗi Hệ Thống)](#11---defer-panic--recover-xử-lý-lỗi-hệ-thống)
12.   [12 - Table-Driven Tests (Kiểm Thử Chuẩn Go)](#12---table-driven-tests-kiểm-thử-chuẩn-go)
13.   [13 - Sync Primitives (Mutex, WaitGroup, Once, Atomics)](#13---sync-primitives-mutex-waitgroup-once-atomics)
14.   [14 - Pipeline Pattern (Fan-In, Fan-Out)](#14---pipeline-pattern-fan-in-fan-out)

---

### 01 - Interfaces (Triết lý Thiết Kế & DI)

- **Khái niệm**: Go định nghĩa interface ngầm định (**Implicit Interface** / _Duck Typing_). Không cần khai báo `implements`. Triết lý cốt lõi của Go là: **"Accept interfaces, return structs"** giúp giảm thiểu ràng buộc và dễ dàng triển khai Dependency Injection.
- 👉 **File thực thi**: [01-interfaces/main.go](./01-interfaces/main.go)

```go
package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct { Radius float64 }
func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

type Rectangle struct { Width, Height float64 }
func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

func printShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f | Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

type Storage interface {
	Save(key, value string) error
	Get(key string) (string, error)
}

type MemoryStorage struct { data map[string]string }
func NewMemoryStorage() *MemoryStorage { return &MemoryStorage{data: make(map[string]string)} }
func (m *MemoryStorage) Save(key, value string) error { m.data[key] = value; return nil }
func (m *MemoryStorage) Get(key string) (string, error) {
	v, ok := m.data[key]; if !ok { return "", fmt.Errorf("key %q not found", key) }; return v, nil
}

type UserService struct { storage Storage }
func NewUserService(s Storage) *UserService { return &UserService{storage: s} }

func main() {
	shapes := []Shape{Circle{Radius: 5}, Rectangle{Width: 4, Height: 6}}
	for _, s := range shapes { printShapeInfo(s) }

	svc := NewUserService(NewMemoryStorage())
	_ = svc.storage.Save("alice", "alice@example.com")
}
```

---

### 02 - Struct Embedding (Composition over Inheritance)

- **Khái niệm**: Go không hỗ trợ kế thừa lớp (Class-based Inheritance) mà sử dụng cơ chế **Struct Embedding** để đóng gói và tái sử dụng code (Composition over Inheritance). Trường được nhúng không có định danh sẽ chuyển giao các method của nó lên struct cha một cách tự động.
- 👉 **File thực thi**: [02-embedding/main.go](./02-embedding/main.go)

```go
package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() { fmt.Printf("Hi, my name is %s!\n", p.Name) }

// Employee nhúng Person
type Employee struct {
	Person // Embedded struct (Anonymous field)
	ID     int
	Salary float64
}

func main() {
	emp := Employee{
		Person: Person{Name: "Alice", Age: 30},
		ID:     1001,
		Salary: 5000,
	}
	emp.Greet() // Gọi trực tiếp method của Person từ Employee
}
```

---

### 03 - Functional Options Pattern

- **Khái niệm**: Giải pháp tối ưu để cấu hình các struct có nhiều trường tùy chọn (optional parameters) mà không cần viết quá nhiều overload constructor hoặc truyền một struct Config rườm rà.
- 👉 **File thực thi**: [03-functional-options/main.go](./03-functional-options/main.go)

```go
package main

import (
	"fmt"
	"time"
)

type Server struct {
	host    string
	port    int
	timeout time.Duration
	tls     bool
}

type ServerOption func(*Server)

func WithPort(port int) ServerOption { return func(s *Server) { s.port = port } }
func WithTimeout(d time.Duration) ServerOption { return func(s *Server) { s.timeout = d } }
func WithTLS() ServerOption { return func(s *Server) { s.tls = true } }

func NewServer(opts ...ServerOption) *Server {
	s := &Server{host: "localhost", port: 8080, timeout: 30 * time.Second}
	for _, opt := range opts { opt(s) }
	return s
}

func main() {
	s := NewServer(WithPort(9090), WithTimeout(10*time.Second), WithTLS())
	fmt.Printf("Server configured: %+v\n", s)
}
```

---

### 04 - Error Handling (Wrapping, Is/As & Custom Errors)

- **Khái niệm**: Go coi lỗi là một giá trị thông thường. Go cung cấp cơ chế bọc lỗi (**Error Wrapping** qua `%w`), cho phép giữ nguyên vết lỗi ban đầu và kiểm tra kiểu lỗi qua `errors.Is()` (so sánh giá trị lỗi) và `errors.As()` (so sánh kiểu struct lỗi).
- 👉 **File thực thi**: [04-error-handling/main.go](./04-error-handling/main.go)

```go
package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("resource not found")

type ValidationError struct {
	Field, Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

func process(id int) error {
	if id < 0 {
		return ValidationError{Field: "ID", Message: "must be positive"}
	}
	if id > 100 {
		return fmt.Errorf("fetching record %d: %w", id, ErrNotFound) // Wrapping
	}
	return nil
}

func main() {
	err := process(120)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("Lỗi không tìm thấy dữ liệu!")
	}

	err2 := process(-5)
	var valErr ValidationError
	if errors.As(err2, &valErr) {
		fmt.Printf("Lỗi Validate: Field %s | Message: %s\n", valErr.Field, valErr.Message)
	}
}
```

---

### 05 - Goroutines & Channels (Đồng Thời Cơ Bản)

- **Khái niệm**: Đơn vị xử lý đồng thời siêu nhẹ (Goroutines) và công cụ truyền tin an toàn giữa các luồng (Channels). Hỗ trợ cơ chế phi chặn (non-blocking) và chia sẻ thông tin mà không cần khóa bộ nhớ thủ công.
- 👉 **File thực thi**: [05-goroutines-channels/main.go](./05-goroutines-channels/main.go)

```go
package main

import "fmt"

func main() {
	ch := make(chan int, 3) // Buffered Channel

	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch) // Đóng channel khi hoàn tất
	}()

	for val := range ch {
		fmt.Println("Received:", val)
	}
}
```

---

### 06 - Select & Timeout (Đa Kênh Bất Đồng Bộ)

- **Khái niệm**: Sử dụng từ khóa `select` để lắng nghe đồng thời trên nhiều channel. Phổ biến nhất là kết hợp `time.After` để thiết lập thời gian chờ tối đa (timeout) tránh treo ứng dụng.
- 👉 **File thực thi**: [06-select-timeout/main.go](./06-select-timeout/main.go)

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "Dữ liệu hoàn tất"
	}()

	select {
	case res := <-ch:
		fmt.Println(res)
	case <-time.After(1 * time.Second): // Hủy sau 1 giây
		fmt.Println("Lỗi: Quá thời gian chờ (Timeout)!")
	}
}
```

---

### 07 - Context (Quản Lý Vòng Đời & Timeout)

- **Khái niệm**: `context.Context` dùng để quản lý thời gian sống (lifecycle), lan truyền tín hiệu hủy bỏ (cancellation) và lưu trữ metadata xuyên suốt các luồng xử lý API/Database.
- 👉 **File thực thi**: [07-context/main.go](./07-context/main.go)

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func performTask(ctx context.Context) {
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Công việc hoàn thành!")
	case <-ctx.Done():
		fmt.Println("Tác vụ bị hủy bỏ:", ctx.Err())
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Giải phóng tài nguyên

	go performTask(ctx)
	time.Sleep(2 * time.Second)
}
```

---

### 08 - Middleware Pattern (HTTP Chain)

- **Khái niệm**: Go triển khai middleware dạng lồng ghép củ hành (Onion Model) dựa trên kiểu hàm `http.Handler` và `http.HandlerFunc`. Mã nguồn trước `next.ServeHTTP` chạy ở chiều đi vào, mã nguồn sau nó chạy ở chiều đi ra.
- 👉 **File thực thi**: [08-middleware/main.go](./08-middleware/main.go)

```go
package main

import (
	"fmt"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // Đi tiếp sang handler sau
		fmt.Printf("[%s] %s - %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello, Go Middleware!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: LoggingMiddleware(mux),
	}
	fmt.Println("Server đang chạy trên cổng 8080...")
	_ = server.ListenAndServe()
}
```

---

### 09 - Closures & State Capture

- **Khái niệm**: Closures (hàm khép kín) cho phép một hàm ẩn danh truy cập và "chụp" (capture) trạng thái biến ở bên ngoài phạm vi của nó. Dùng để tạo các hàm phát sinh tự động, cấu hình động hoặc lưu cache.
- 👉 **File thực thi**: [09-closures/main.go](./09-closures/main.go)

```go
package main

import "fmt"

func Generator() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	next := Generator()
	fmt.Println(next()) // 1
	fmt.Println(next()) // 2
	fmt.Println(next()) // 3
}
```

---

### 10 - Generics (Type Parameters & Constraints)

- **Khái niệm**: Từ phiên bản **Go 1.18+**, Go hỗ trợ lập trình tổng quát (**Generics**) sử dụng các tham số kiểu dữ liệu (`[T any]`) và ràng buộc kiểu (`any`, `comparable`, interface constraints) giúp tái sử dụng mã nguồn an toàn lúc biên dịch.
- 👉 **File thực thi**: [10-generics/main.go](./10-generics/main.go)

```go
package main

import "fmt"

// Generic Map keys to slices
func GetKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	m := map[string]int{"one": 1, "two": 2}
	keys := GetKeys(m)
	fmt.Println("Keys:", keys)
}
```

---

### 11 - Defer, Panic & Recover (Xử Lý Lỗi Hệ Thống)

- **Khái niệm**: `defer` trì hoãn việc gọi hàm đến khi kết thúc hàm chứa nó (dọn dẹp tài nguyên LIFO). `panic` dừng khẩn cấp hệ thống, và `recover` bắt giữ `panic` bên trong hàm `defer` để cứu server khỏi sập đột ngột.
- 👉 **File thực thi**: [11-defer-panic-recover/main.go](./11-defer-panic-recover/main.go)

```go
package main

import "fmt"

func SafeExecution() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Đã giải cứu thành công hệ thống khỏi panic:", r)
		}
	}()
	fmt.Println("Bắt đầu thực hiện tác vụ...")
	panic("Lỗi chia cho 0 hoặc lỗi nil pointer nghiêm trọng!")
}

func main() {
	SafeExecution()
	fmt.Println("Ứng dụng vẫn tiếp tục hoạt động an toàn!")
}
```

---

### 12 - Table-Driven Tests (Kiểm Thử Chuẩn Go)

- **Khái niệm**: Phong cách kiểm thử chuẩn mực của cộng đồng Go. Nhóm toàn bộ các ca kiểm thử (test cases) thành một mảng dữ liệu (Table), duyệt qua mảng và gọi `t.Run()` để kiểm thử giúp code kiểm thử gọn gàng và dễ mở rộng.
- 👉 **File thực thi**: [12-table-driven-tests/main.go](./12-table-driven-tests/main.go)

```go
package main

import "testing"

func Add(a, b int) int { return a + b }

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -1, -1, -2},
		{"zero value", 0, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.want {
				t.Errorf("Add() = %d, want %d", got, tt.want)
			}
		})
	}
}
```

---

### 13 - Sync Primitives (Mutex, WaitGroup, Once, Atomics)

- **Khái niệm**: Đồng bộ luồng và bộ nhớ an toàn.
   - `sync.Mutex` / `sync.RWMutex`: Khóa đọc/ghi bảo vệ bộ nhớ khỏi race condition.
   - `sync.WaitGroup`: Đợi một nhóm Goroutines hoàn thành.
   - `sync.Once`: Đảm bảo một hàm (như khởi tạo DB) chỉ chạy duy nhất 1 lần trong đời ứng dụng.
   - `sync/atomic`: Phép toán nguyên tử cấp CPU cực nhẹ không cần lock bộ nhớ.
- 👉 **File thực thi**: [13-sync-primitives/main.go](./13-sync-primitives/main.go)

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v[key]++
}

func main() {
	var count atomic.Uint64 // Atomic Counter
	var wg sync.WaitGroup

	counter := SafeCounter{v: make(map[string]int)}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc("clicks")
			count.Add(1)
		}()
	}
	wg.Wait()
	fmt.Println("Mutex clicks:", counter.v["clicks"])
	fmt.Println("Atomic count:", count.Load())
}
```

---

### 14 - Pipeline Pattern (Fan-In, Fan-Out)

- **Khái niệm**: Mô hình xử lý luồng dữ liệu lớn tuần tự hoặc song song qua nhiều giai đoạn (Stages). Sử dụng **Fan-Out** (chia công việc ra cho nhiều Goroutines rảnh cùng làm) và **Fan-In** (thu thập toàn bộ kết quả từ nhiều channel về một kênh chính duy nhất).
- 👉 **File thực thi**: [14-pipeline-pattern/main.go](./14-pipeline-pattern/main.go)

```go
package main

import (
	"fmt"
	"sync"
)

// Giai đoạn 1: Sinh số nguyên ngẫu nhiên
func GenerateNumbers(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

// Giai đoạn 2: Nhân đôi giá trị nhận được
func DoubleNumber(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * 2:
			case <-done:
				return
			}
		}
	}()
	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)

	numsCh := GenerateNumbers(done, 1, 2, 3, 4)
	results := DoubleNumber(done, numsCh)

	for res := range results {
		fmt.Println("Pipeline result:", res)
	}
}
```
