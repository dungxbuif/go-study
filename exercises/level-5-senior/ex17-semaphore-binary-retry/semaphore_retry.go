package main

import (
	"errors"
	"sync"
	"time"
)

// ProcessItem xử lý một phần tử đơn lẻ. Mô phỏng tác vụ I/O tốn thời gian.
func ProcessItem(item int) error {
	time.Sleep(50 * time.Millisecond)
	// Giả lập lỗi ở các phần tử cụ thể để test thuật toán phân tách lỗi
	if item == 37 || item == 82 {
		return errors.New("simulated item failure")
	}
	return nil
}

// ProcessChunk thực hiện xử lý một chunk (nhóm phần tử) song song, giới hạn bởi Semaphore.
//
// 🧠 SEMAPHORE BẰNG BUFFERED CHANNEL DƯỚI NẮP CAPO:
// - Trong Node.js, để khống chế số lượng Promise chạy đồng thời (Concurrency Limit), ta phải tự viết
//   hoặc dùng thư viện như `p-limit` với hàng đợi (Queue) phức tạp.
// - Trong Go, ta sử dụng **Buffered Channel** làm Semaphore cực kỳ tinh gọn:
//   `sem := make(chan struct{}, limit)`
//   - **Acquire (Mượn slot)**: Gửi dữ liệu vào channel: `sem <- struct{}{}`.
//     Nếu channel đầy buffer (đã đạt giới hạn giới hạn song song), Goroutine hiện tại sẽ bị chặn (Block)
//     bởi Go Scheduler và đưa vào hàng đợi chờ mà không tốn tài nguyên CPU (Sleep state).
//   - **Release (Trả slot)**: Nhận dữ liệu ra khỏi channel: `<-sem`.
//     Thao tác này giải phóng 1 ô trống trong buffer, lập tức đánh thức Goroutine đang chờ tiếp theo.
//   - **Kiểu `struct{}` (Empty Struct)**: Trong Go, `struct{}` là một kiểu dữ liệu đặc biệt chiếm **0 bytes** bộ nhớ.
//     Sử dụng `chan struct{}` giúp Semaphore hoạt động siêu nhẹ, hoàn toàn không tốn dung lượng RAM!
func ProcessChunk(chunk []int, sem chan struct{}) ([]int, []int) {
	// TODO:
	// 1. Duyệt qua mảng chunk. Với mỗi item, acquire slot trong sem channel.
	// 2. Kích hoạt Goroutine chạy ProcessItem(item).
	// 3. Sử dụng sync.WaitGroup để đợi toàn bộ các Goroutines trong chunk hoàn thành.
	// 4. Giải phóng (release) slot trong sem khi Goroutine hoàn tất.
	// 5. Trả về: (successItems, failedItems)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var success []int
	var failed []int

	// Viết logic tại đây...

	_ = wg
	_ = mu
	return success, failed
}

// ProcessChunkWithBinaryRetry thực hiện chia đôi thử lại (Binary Split Retry).
//
// 🧠 THUẬT TOÁN CO LẬP LỖI BẰNG CHIA ĐÔI ĐỆ QUY (Binary Split Retry):
// - Trong hệ thống xử lý Batch lớn (như Import dữ liệu, Gửi Mail hàng loạt), nếu ta gửi 100 items và DB trả về lỗi chung chung,
//   ta không biết item nào lỗi. Nếu bỏ qua cả 100 thì lãng phí, nếu thử lại từng cái một (O(N)) thì quá chậm.
// - Giải pháp: **Binary Split Retry (Divide and Conquer)**:
//   1. Chạy cả nhóm 100 items. Nếu thành công -> Kết thúc.
//   2. Nếu phát hiện có lỗi:
//      - Nếu chỉ có 1 item lỗi -> Dễ dàng cô lập và báo cáo lỗi vĩnh viễn (Permanent Failure) cho duy nhất item đó.
//      - Nếu nhiều hơn 1 item lỗi -> Chia đôi nhóm thành 2 nhóm con (mỗi nhóm 50 items) và đệ quy chạy song song.
//   3. Tiếp tục lặp lại quá trình chia đôi. Độ phức tạp tìm kiếm lỗi giảm xuống còn O(Log N),
//      cho phép cô lập chính xác các phần tử dữ liệu lỗi cực nhanh mà không làm tắc nghẽn hay quá tải hệ thống!
func ProcessChunkWithBinaryRetry(chunk []int, sem chan struct{}, permanentFailed *[]int, mu *sync.Mutex) {
	// TODO:
	// 1. Gọi ProcessChunk(chunk, sem)
	// 2. Nếu không có item nào lỗi -> kết thúc thành công.
	// 3. Nếu có đúng 1 item bị lỗi -> Ghi vào danh sách permanentFailed (cần dùng mutex bảo vệ).
	// 4. Nếu có nhiều hơn 1 item lỗi -> Chia đôi mảng failedItems (Left, Right)
	//    và gọi đệ quy song song ProcessChunkWithBinaryRetry bằng Goroutines.
}

func main() {
	// Demo Channel Semaphore + Binary Split Retry tại đây
}
