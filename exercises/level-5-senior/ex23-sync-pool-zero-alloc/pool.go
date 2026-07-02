package main

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

// UserRequest biểu diễn dữ liệu request cần xử lý.
type UserRequest struct {
	ID        int
	Data      []byte
	Processed bool
}

// 🧠 CƠ CHẾ SYNC.POOL & TỐI ƯU HÓA TRÁNH RÒ RỈ GC (sync.Pool under the hood):
// - Trong các ứng dụng High-Throughput (như HTTP/TCP Server nhận hàng chục vạn requests mỗi giây), việc khởi tạo liên tục
//   các biến đệm (Buffer, Slice, Struct) sẽ tạo ra một lượng khổng lồ rác trên Heap.
// - Garbage Collector (GC) của Go sẽ phải quét dọn liên tục, gây ra hiện tượng nghẽn CPU và các đợt trễ dừng thế giới (Stop-The-World Latency Spikes).
// - **`sync.Pool`** sinh ra để giải quyết vấn đề này bằng cách tái sử dụng (recycle) các đối tượng đã được phân bổ:
//   - **Get()**: Lấy một đối tượng nhàn rỗi trong Pool. Nếu Pool trống, nó tự động chạy hàm `New` để cấp phát mới.
//   - **Put()**: Trả đối tượng về Pool sau khi dùng xong.
// - ⚠️ LƯU Ý SỐNG CÒN (Đặc tính dọn dẹp của sync.Pool):
//   1. Đối tượng trong `sync.Pool` có thể bị GC tự động thu hồi bất cứ lúc nào trong chu kỳ quét GC (thường là trước mỗi đợt quét).
//      Do đó, ta chỉ dùng `sync.Pool` cho các đối tượng TẠM THỜI (Short-lived), KHÔNG ĐƯỢC dùng cho các đối tượng lưu trữ trạng thái lâu dài.
//   2. BẮT BUỘC phải **Reset (Xóa sạch dữ liệu)** của đối tượng trước khi Put lại vào Pool.
//      Nếu không reset (ví dụ một buffer chứa chuỗi tin cũ), request tiếp theo lấy đối tượng đó ra sẽ bị dính dữ liệu cũ của request trước,
//      gây lỗi rò rỉ thông tin cực kỳ nghiêm trọng (Information Leak/Security Bug)!
var requestPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("[Pool] Allocating new UserRequest object on Heap.")
		return &UserRequest{
			Data: make([]byte, 1024), // Khởi tạo buffer đệm 1KB
		}
	},
}

func ProcessRequest(id int, payload []byte) {
	// Lấy đối tượng từ pool
	req := requestPool.Get().(*UserRequest)
	
	// Trả lại pool khi xử lý xong
	defer func() {
		// BẮT BUỘC: Reset dữ liệu để tránh nhiễm bẩn ô nhớ (dirty memory pollution) cho request sau
		req.ID = 0
		req.Processed = false
		// Xóa trắng buffer
		for i := range req.Data {
			req.Data[i] = 0
		}
		
		requestPool.Put(req)
	}()

	// Gán thông tin mới
	req.ID = id
	copy(req.Data, payload)
	req.Processed = true

	// Giả lập xử lý nghiệp vụ ghi log hoặc nén dữ liệu
	_ = bytes.HasPrefix(req.Data, []byte("API"))
}

func main() {
	payload := []byte("API-Payload-Data")

	fmt.Println("--- 1. First iteration (Pool is empty, expecting allocations) ---")
	for i := 1; i <= 3; i++ {
		ProcessRequest(i, payload)
	}

	// Đợi một khoảng ngắn (không có GC quét)
	time.Sleep(10 * time.Millisecond)

	fmt.Println("\n--- 2. Second iteration (Reusing objects from Pool, 0 allocation expected) ---")
	// Lúc này, các đối tượng đã được Put ngược lại Pool, ta lấy ra dùng lại mà không tốn thêm Heap allocation
	for i := 4; i <= 6; i++ {
		ProcessRequest(i, payload)
	}
}
