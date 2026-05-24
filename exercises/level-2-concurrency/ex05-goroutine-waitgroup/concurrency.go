package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ============================================================================
// 🧠 GIẢI THÍCH VỀ ĐỒNG THỜI (CONCURRENCY) TRONG GO:
//
// 1. Goroutine (`go func()`):
//    - Là các luồng xử lý siêu nhẹ quản lý bởi Go Runtime (chỉ tốn khoảng 2KB RAM ban đầu,
//      trong khi OS Thread tốn khoảng 1-2MB).
//    - Bản chất: Giúp chạy các dòng code bất đồng bộ song song với luồng chính (main).
//
// 2. sync.WaitGroup:
//    - Dùng để đồng bộ hoá (đồng bộ kết thúc). Nó giống như một bộ đếm (counter).
//    - `wg.Add(1)`: Cộng 1 vào bộ đếm trước khi chạy Goroutine.
//    - `wg.Done()`: Trừ 1 khỏi bộ đếm khi Goroutine kết thúc. Dùng `defer` để đảm bảo luôn chạy.
//    - `wg.Wait()`: Khoá luồng chính lại, bắt nó đợi cho đến khi bộ đếm quay về bằng 0.
//
// 3. sync.Mutex (Mutual Exclusion):
//    - Mặc dù Event Loop của JS đơn luồng nên không lo race condition, Go chạy đa luồng thực sự.
//    - Slice `results` trong Go KHÔNG thread-safe. Nếu nhiều Goroutine cùng gọi `append` vào slice
//      cùng một thời điểm, dữ liệu sẽ bị ghi đè lung tung hoặc gây sụp đổ chương trình (panic).
//    - `mu.Lock()` & `mu.Unlock()`: Chỉ cho phép ĐÚNG 1 Goroutine được thực thi đoạn code ghi dữ liệu
//      ở giữa tại một thời điểm. Các Goroutine khác đến sau phải xếp hàng đợi.
// ============================================================================

type Result struct {
	URL      string
	Status   int
	Duration time.Duration
}

// CheckURL gửi request HTTP GET đơn lẻ và trả về thông tin kết quả.
func CheckURL(client *http.Client, url string) Result {
	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		return Result{URL: url, Status: 0, Duration: duration}
	}
	defer resp.Body.Close() // Đảm bảo đóng body sau khi đọc xong để tránh rò rỉ kết nối (connection leak)

	return Result{URL: url, Status: resp.StatusCode, Duration: duration}
}

// CheckAllURLs kiểm tra song song toàn bộ danh sách URLs nhận được.
func CheckAllURLs(urls []string) []Result {
	// Khởi tạo một HTTP Client dùng chung với thời gian hết hạn là 5 giây
	client := &http.Client{Timeout: 5 * time.Second}
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []Result

	for _, url := range urls {
		// 1. Tăng bộ đếm WaitGroup lên 1 trước khi khởi chạy Goroutine mới
		wg.Add(1)

		// 2. Sử dụng từ khóa `go` trước một hàm ẩn danh (anonymous function) để chạy nó bất đồng bộ.
		// Truyền trực tiếp biến `url` vào tham số `u string` của hàm.
		// Việc truyền tham số này giúp chụp lại giá trị (capture value) tại từng vòng lặp,
		// tránh việc các Goroutine dùng chung một biến lặp gây sai lệch dữ liệu.
		go func(u string) {
			// 3. Đăng ký giải phóng bộ đếm sau khi hàm chạy xong bằng `defer`.
			// Dù hàm có chạy lỗi hay bị hoảng loạn (panic) thì wg.Done() vẫn luôn được gọi.
			defer wg.Done()

			res := CheckURL(client, u)

			// 4. Sử dụng Mutex để bảo vệ thao tác append vào kết quả dùng chung
			mu.Lock()
			results = append(results, res) // Đoạn mã nguy hiểm (Critical Section)
			mu.Unlock()
		}(url)
	}

	// 5. Khoá hàm main lại ở đây.
	// Chỉ khi toàn bộ Goroutines gọi wg.Done() và bộ đếm về 0, hàm mới tiếp tục chạy xuống dưới.
	wg.Wait()
	
	return results
}

func main() {
	urls := []string{
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/201",
		"https://httpbin.org/status/404",
		"https://httpbin.org/status/500",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/2",
	}

	fmt.Println("=== Checking URLs in Parallel ===")
	start := time.Now()
	results := CheckAllURLs(urls)
	duration := time.Since(start)

	for _, res := range results {
		fmt.Printf("URL: %s | Status: %d | Time: %v\n", res.URL, res.Status, res.Duration)
	}
	fmt.Printf("\nTotal duration: %v\n", duration)
}
