package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// FetchData thực hiện một HTTP GET request với context được đính kèm.
//
// 💡 CƠ CHẾ DƯỚI NẮP CAPO (Go http.Client & Context):
// 1. Khi ta truyền `ctx` vào `http.NewRequestWithContext`, Go sẽ ràng buộc vòng đời của request này vào context.
// 2. Bên dưới http package, khi thực hiện `client.Do(req)`, một goroutine khác sẽ giám sát kênh `ctx.Done()`.
// 3. Nếu context bị cancel hoặc timeout trước khi request xong, kênh `ctx.Done()` sẽ đóng (close).
// 4. Goroutine giám sát phát hiện `Done()` đóng và sẽ đóng kết nối TCP socket underlying ngay lập tức (connection reset/abort),
//    ngăn chặn lãng phí tài nguyên mạng và CPU. Điều này tương tự như việc hủy bỏ socket connection trong Node.js.
func FetchData(ctx context.Context, client *http.Client, url string) (string, error) {
	// Gắn context vào request. Nếu không truyền context, request có thể bị treo vô hạn nếu server không phản hồi.
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		// Nếu bị timeout, err ở đây sẽ trả về lỗi chứa: "context deadline exceeded"
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// FetchMultiple thực hiện tải dữ liệu từ nhiều URL song song (concurrently) bằng Goroutines.
// Toàn bộ quá trình fetch này bị giới hạn bởi một khoảng thời gian timeout chung.
//
// 🧠 SO SÁNH GOLANG CONTEXT VS NODE.JS ABORTCONTROLLER:
// - Go Context:
//   - Là một interface định nghĩa 4 phương thức: `Deadline()`, `Done()`, `Err()`, và `Value()`.
//   - `context.WithTimeout` tạo ra một child context và một hàm `cancel`. Dưới nắp capo, Go sử dụng một timer (`time.AfterFunc`)
//     để tự động gọi `cancel` sau thời gian `timeout` trôi qua.
//   - Hàm `cancel` có nhiệm vụ đóng một channel nội bộ kiểu `chan struct{}` (trả về bởi `ctx.Done()`).
//   - Việc đóng channel này hoạt động như một cơ chế broadcast: TẤT CẢ các bên đang lắng nghe `<-ctx.Done()` sẽ nhận được tín hiệu tức thì.
//   - BẮT BUỘC phải gọi `defer cancel()`. Nếu các goroutines hoàn thành trước khi hết hạn, gọi `cancel()` sẽ dọn dẹp timer nội bộ.
//     Nếu không, timer sẽ tiếp tục nằm trong bộ nhớ cho đến khi hết hạn mới được Garbage Collection giải phóng, gây memory leak!
// - Node.js AbortController:
//   - Sử dụng `new AbortController()` và truyền `signal` vào Axios hoặc Fetch.
//   - Khi hết hạn, ta gọi `controller.abort()`. Hàm này sẽ kích hoạt sự kiện 'abort' trên đối tượng `AbortSignal`.
//   - Phải tự dọn dẹp `clearTimeout(timeoutId)` để tránh rò rỉ timer của Node.js Event Loop.
func FetchMultiple(urls []string, timeout time.Duration) ([]string, error) {
	// Khởi tạo một context với timeout.
	// context.Background() đóng vai trò là "gốc" (Root Context) không thể bị hủy.
	// ctx là context con kế thừa từ Background nhưng có thêm cơ chế tự hủy (deadline).
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	
	// Luôn defer cancel() để giải phóng các tài nguyên (như timer) ngay khi hàm kết thúc,
	// tránh rò rỉ bộ nhớ (leak timer/resources) kể cả khi tiến trình kết thúc sớm hơn timeout.
	defer cancel()

	client := &http.Client{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []string
	var fetchErr error

	for _, url := range urls {
		wg.Add(1)
		// Tạo một goroutine cho mỗi URL để fetch song song (Tương tự Promise.all)
		go func(u string) {
			defer wg.Done()
			
			// Truyền context xuống cho FetchData. Nếu context này bị timeout,
			// FetchData sẽ dừng ngay lập tức và trả về lỗi.
			data, err := FetchData(ctx, client, u)
			
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				// Lưu lỗi đầu tiên xảy ra
				fetchErr = err
			} else {
				results = append(results, data)
			}
		}(url)
	}

	// Đợi tất cả các Goroutines hoàn thành.
	// CHÚ Ý: wg.Wait() sẽ đợi cho đến khi TẤT CẢ goroutines trả về. Nếu một request bị nghẽn mạng quá lâu,
	// nhưng nhờ có `ctx` truyền vào, kết nối sẽ bị ngắt sau đúng thời gian `timeout` đã định nghĩa.
	// Do đó, các Goroutines được đảm bảo sẽ kết thúc nhanh chóng và không bị treo vô hạn.
	wg.Wait()
	
	if fetchErr != nil {
		return nil, fetchErr
	}
	return results, nil
}

func main() {
	urls := []string{
		"https://httpbin.org/delay/1",
		"https://httpbin.org/status/200",
	}

	fmt.Println("=== Fetching with 3s Timeout ===")
	// Với 3s timeout, request delay 1s chắc chắn sẽ thành công
	results, err := FetchMultiple(urls, 3*time.Second)
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
	} else {
		fmt.Printf("Successfully fetched %d responses\n", len(results))
	}

	fmt.Println("\n=== Fetching with 500ms Timeout (expect fail) ===")
	// Với 500ms timeout, request delay 1s chắc chắn sẽ thất bại vì vượt quá deadline.
	_, err = FetchMultiple(urls, 500*time.Millisecond)
	if err != nil {
		fmt.Printf("Expected error occurred: %v\n", err) // Sẽ in ra lỗi context deadline exceeded
	}
}
