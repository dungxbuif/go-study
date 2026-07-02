package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// State đại diện cho các trạng thái của Circuit Breaker.
type State int

const (
	Closed State = iota // Mạch đóng (Bình thường, cho phép request đi qua)
	Open                // Mạch hở (Có lỗi, ngăn chặn request, gọi ngay fallback)
	HalfOpen            // Mạch hé mở (Thử nghiệm gửi một vài request để kiểm tra hệ thống phục hồi chưa)
)

func (s State) String() string {
	switch s {
	case Closed:
		return "CLOSED"
	case Open:
		return "OPEN"
	case HalfOpen:
		return "HALF-OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreaker bảo vệ hệ thống khỏi các thảm họa sập tầng dây chuyền (Cascading Failures).
//
// 🧠 CƠ CHẾ CIRCUIT BREAKER (Cắt mạch bảo vệ hệ thống):
// - Trong hệ thống Microservices, nếu một dịch vụ con (Downstream) bị sập hoặc phản hồi cực chậm,
//   dịch vụ gọi nó (Upstream) sẽ bị nghẽn các goroutines chờ đợi, dẫn đến cạn kiệt tài nguyên (CPU, socket descriptors) và sập toàn bộ hệ thống.
// - Circuit Breaker hoạt động giống như một chiếc Aptomat điện:
//   1. **CLOSED**: Mọi thứ bình thường. Ta theo dõi tỉ lệ lỗi. Nếu lỗi vượt ngưỡng (`FailureThreshold`), mạch ngắt chuyển sang **OPEN**.
//   2. **OPEN**: Từ chối kết nối ngay lập tức (Fail-Fast) mà không thèm gọi xuống downstream, chuyển tiếp trực tiếp sang hàm **Fallback** để cứu hộ.
//      Điều này giúp Downstream có thời gian "thở" để tự phục hồi và bảo vệ Upstream không bị nghẽn luồng chạy.
//   3. **HALF-OPEN**: Sau một khoảng thời gian chờ (`CooldownWindow`), mạch chuyển sang Half-Open, cho phép một vài request đi qua thử.
//      - Nếu thành công: Mạch đóng lại (**CLOSED**), hệ thống phục hồi bình thường.
//      - Nếu thất bại: Mạch ngắt tiếp (**OPEN**), bắt đầu lại chu kỳ đếm ngược.
type CircuitBreaker struct {
	mu               sync.Mutex
	state            State
	failureCount     int
	failureThreshold int
	cooldownWindow   time.Duration
	lastStateChange  time.Time
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            Closed,
		failureThreshold: threshold,
		cooldownWindow:   cooldown,
		lastStateChange:  time.Now(),
	}
}

// Execute thực thi hàm nghiệp vụ được bảo vệ bởi Circuit Breaker.
//
// 🧠 ĐA LUỒNG & THỜI GIAN THỰC (Thread-Safe State Transition):
// - Circuit Breaker được gọi từ nhiều Goroutines đồng thời, do đó ta dùng `sync.Mutex` để bảo vệ các thuộc tính trạng thái.
// - Dưới nắp capo, ta kiểm tra xem nếu trạng thái là OPEN và đã vượt quá `cooldownWindow`,
//   ta chủ động chuyển mạch sang HALF-OPEN ngay lập tức (Lazy Transition) để cho phép request này đi qua thử nghiệm.
func (cb *CircuitBreaker) Execute(action func() error, fallback func() error) error {
	cb.mu.Lock()
	
	// lazy state transition from Open to HalfOpen
	if cb.state == Open && time.Since(cb.lastStateChange) > cb.cooldownWindow {
		fmt.Println("[CircuitBreaker] Cooldown elapsed. Transitioning to HALF-OPEN")
		cb.state = HalfOpen
		cb.lastStateChange = time.Now()
	}

	currentState := cb.state
	cb.mu.Unlock()

	// Nếu mạch đang mở (OPEN), chạy ngay fallback cứu hộ mà không gọi hàm action bị lỗi
	if currentState == Open {
		fmt.Println("[CircuitBreaker] Circuit is OPEN. Triggering fallback fast-fail.")
		return fallback()
	}

	// Thực thi hành động thực tế (có thể gây lỗi hoặc chậm)
	err := action()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.onFailure()
		return fallback() // Nếu hành động thực tế lỗi, trả về kết quả fallback
	}

	cb.onSuccess()
	return nil
}

func (cb *CircuitBreaker) onFailure() {
	cb.failureCount++
	fmt.Printf("[CircuitBreaker] Failure recorded. Count: %d, State: %s\n", cb.failureCount, cb.state)

	if cb.state == Closed && cb.failureCount >= cb.failureThreshold {
		fmt.Println("[CircuitBreaker] Threshold exceeded. Tripping circuit to OPEN!")
		cb.state = Open
		cb.lastStateChange = time.Now()
	} else if cb.state == HalfOpen {
		fmt.Println("[CircuitBreaker] Test request failed in HALF-OPEN. Tripping circuit back to OPEN!")
		cb.state = Open
		cb.lastStateChange = time.Now()
	}
}

func (cb *CircuitBreaker) onSuccess() {
	fmt.Printf("[CircuitBreaker] Success recorded. State: %s\n", cb.state)
	if cb.state == HalfOpen {
		fmt.Println("[CircuitBreaker] Test request succeeded in HALF-OPEN. Closing circuit back to CLOSED!")
		cb.state = Closed
		cb.failureCount = 0
	} else if cb.state == Closed {
		cb.failureCount = 0 // Reset bộ đếm lỗi khi thành công
	}
}

func main() {
	cb := NewCircuitBreaker(3, 2*time.Second)

	// Mô phỏng hàm gọi API bên thứ ba có nguy cơ sập
	unstableAPI := func(shouldFail bool) func() error {
		return func() error {
			if shouldFail {
				return errors.New("third-party API error")
			}
			return nil
		}
	}

	// Hàm Fallback cung cấp dữ liệu mặc định/từ Cache dự phòng
	fallbackData := func() error {
		fmt.Println("[Fallback] Returning cached static mock data.")
		return nil
	}

	fmt.Println("--- 1. Testing normal behavior (CLOSED) ---")
	cb.Execute(unstableAPI(false), fallbackData)

	fmt.Println("\n--- 2. Inducing errors to trip the circuit ---")
	for i := 0; i < 4; i++ {
		cb.Execute(unstableAPI(true), fallbackData)
	}

	fmt.Println("\n--- 3. Circuit is now OPEN. Requests should fail-fast instantly without calling unstable API ---")
	cb.Execute(unstableAPI(false), fallbackData)

	fmt.Println("\n--- 4. Waiting for Cooldown window to expire (2 seconds) ---")
	time.Sleep(2100 * time.Millisecond)

	fmt.Println("\n--- 5. Making a request now. Mạch sẽ tự động chuyển Half-Open và thành công đóng mạch ---")
	cb.Execute(unstableAPI(false), fallbackData)

	fmt.Println("\n--- 6. Mạch đã đóng lại (CLOSED). Kiểm chứng hoạt động bình thường ---")
	cb.Execute(unstableAPI(false), fallbackData)
}
