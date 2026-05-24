package main

import (
	"errors"
	"sync"
	"time"
)

func ProcessItem(item int) error {
	time.Sleep(50 * time.Millisecond)
	if item == 37 || item == 82 {
		return errors.New("simulated item failure")
	}
	return nil
}

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
