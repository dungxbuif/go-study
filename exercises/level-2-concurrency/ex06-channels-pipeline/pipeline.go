package main

import "fmt"

// ============================================================================
// 💡 CƠ CHẾ HOẠT ĐỘNG "DƯỚI NẮP CAPO" CỦA GO CHANNEL PIPELINE
// ============================================================================
//
// 1. Chạy Đồng Thời Thực Sự (True Concurrency):
//    - Khi gọi `Generator`, `Filter`, và `Squarer`, mỗi Stage sẽ sinh ra một **Goroutine** độc lập
//      nhờ từ khóa `go func()`.
//    - GMP Scheduler của Go tự động phân phối các Goroutines này chạy song song trên nhiều luồng
//      hệ điều hành (OS Threads) và tận dụng tối đa đa nhân CPU.
//
// 2. Kênh Truyền Tin Đồng Bộ (Unbuffered Channels):
//    - Ở đây, chúng ta sử dụng Unbuffered Channels (`make(chan int)` không chỉ định kích thước buffer).
//    - **Đặc tính cốt lõi:** Lệnh gửi dữ liệu (`out <- n`) BẮT BUỘC phải đợi đến khi có lệnh nhận dữ liệu
//      (`<-in` hoặc `for n := range in` ở Stage tiếp theo) sẵn sàng, và ngược lại.
//    - Nếu không có ai nhận, Goroutine gửi sẽ bị **chặn lại (Block/Suspend)**. Trình điều phối Go sẽ
//      "cất" Goroutine này đi (Parked) để không tốn CPU, và đánh thức nó dậy (Unparked) khi có người nhận.
//
// 3. Cơ Chế Tự Động Kiểm Soát Tốc Độ (Automatic Backpressure):
//    - Đây là ưu điểm vượt trội của Go so với Node.js.
//    - Nếu Stage cuối (ví dụ main ghi DB) chạy rất chậm, nó sẽ block việc nhận dữ liệu từ `Squarer`.
//      `Squarer` bị block lệnh nhận sẽ lập tức block lệnh gửi của nó sang `Filter`. Cứ thế ngược dòng,
//      toàn bộ hệ thống tự động chậm lại theo Stage chậm nhất.
//    - **Hệ quả:** RAM không bao giờ bị phình to vì không có mảng đệm nào tích lũy dữ liệu rác!
//
// 4. Hiệu Ứng Domino Khi Đóng Kênh (`close(out)`):
//    - Khi `Generator` chạy hết mảng, nó gọi `close(out)`.
//    - Stage `Filter` đang chờ ở `for n := range in` nhận thấy channel bị đóng ➡️ thoát vòng lặp ➡️ gọi `close(out)` của nó.
//    - Stage `Squarer` nhận tín hiệu đóng ➡️ thoát vòng lặp ➡️ gọi `close(out)` của nó.
//    - Hàm `main` thoát vòng lặp `range` và kết thúc chương trình. Lỗi rò rỉ luồng (Goroutine Leak) hoàn toàn bằng 0!
// ============================================================================

// Generator: Stage 1 - Sinh dữ liệu từ slice chuyển vào channel
// Trả về một kênh chỉ đọc (receive-only channel: `<-chan int`)
func Generator(nums []int) <-chan int {
	out := make(chan int)
	
	// Sinh ra Goroutine độc lập chạy ngầm
	go func() {
		for _, n := range nums {
			// Gửi từng số `n` vào channel.
			// Dòng này sẽ BỊ CHẶN cho đến khi Stage 2 (Filter) bắt đầu đọc giá trị này.
			out <- n
		}
		// Đóng channel để báo hiệu cho Stage tiếp theo rằng: "Không còn dữ liệu nào nữa!"
		close(out)
	}()
	
	// Trả về channel ngay lập tức (không đợi Goroutine chạy xong - Non-blocking return)
	return out
}

// Filter: Stage 2 - Lọc dữ liệu chẵn
// Nhận vào một channel chỉ đọc `in`, trả về một channel chỉ đọc `out`
func Filter(in <-chan int) <-chan int {
	out := make(chan int)
	
	go func() {
		// Vòng lặp `range in` sẽ liên tục đọc từ channel `in`.
		// Vòng lặp này tự động kết thúc khi channel `in` bị đóng ở Stage 1.
		for n := range in {
			if n%2 == 0 {
				// Đẩy số chẵn sang Stage 3 (Squarer).
				// Tiếp tục bị chặn nếu Stage 3 chưa sẵn sàng đọc.
				out <- n
			}
		}
		// Đóng channel đầu ra để báo hiệu cho Stage 3
		close(out)
	}()
	
	return out
}

// Squarer: Stage 3 - Bình phương các số nhận được
func Squarer(in <-chan int) <-chan int {
	out := make(chan int)
	
	go func() {
		for n := range in {
			// Đẩy giá trị đã bình phương ra ngoài
			out <- n * n
		}
		close(out)
	}()
	
	return out
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Thiếp lập đường ống (Pipeline Wiring):
	// 1. Generator gửi dữ liệu vào channel `c`
	c := Generator(nums)
	// 2. Filter đọc từ `c`, lọc chẵn gửi vào channel tiếp theo.
	// 3. Squarer đọc từ channel của Filter, bình phương và trả về channel cuối `out`.
	out := Squarer(Filter(c))

	fmt.Println("=== Pipeline output ===")
	// Vòng lặp đọc dữ liệu cuối cùng từ channel `out`
	// Vòng lặp này sẽ chặn hàm `main` cho đến khi có dữ liệu chảy ra từ pipeline.
	// Tự động thoát khi tất cả các Stage đã đóng channel tuần tự theo hiệu ứng Domino.
	for val := range out {
		fmt.Println(val)
	}
}
