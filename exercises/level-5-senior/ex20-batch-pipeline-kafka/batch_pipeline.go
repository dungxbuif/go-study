package main

import (
	"context"
	"time"
)

// Batcher định nghĩa bộ gom lô dữ liệu để tối ưu hóa tần suất ghi Database hoặc Kafka.
type Batcher struct {
	queue     chan string   // Hàng đợi lưu trữ tạm thời các item trước khi gom lô
	batchSize int           // Kích thước tối đa của mỗi lô (ví dụ: 100 items)
	timeout   time.Duration // Thời gian chờ tối đa trước khi tự động flush (ví dụ: 1 giây)
}

// NewBatcher khởi tạo một Batcher mới.
// Hàng đợi queue được cấu hình với buffer lớn (ví dụ: 10,000) giúp hấp thụ xung lực tải (Traffic Spike) cực tốt,
// đóng vai trò như một bộ đệm giảm chấn giữa Producers (nhanh) và Consumers/Database (chậm).
func NewBatcher(batchSize int, timeout time.Duration) *Batcher {
	return &Batcher{
		queue:     make(chan string, 10000),
		batchSize: batchSize,
		timeout:   timeout,
	}
}

// Push gửi một phần tử vào hàng đợi để chờ xử lý theo lô.
func (b *Batcher) Push(item string) {
	b.queue <- item
}

// Run khởi chạy vòng lặp xử lý sự kiện trung tâm của Batcher.
//
// 🧠 MÔ HÌNH GHÉP KÊNH SELECT GOM LÔ (Select Multiplexing under the hood):
// - Trong Node.js, việc gom lô thường dựa trên `setTimeout` để tự động flush sau một khoảng thời gian.
//   Mã nguồn JS phải tự quản lý trạng thái timer (xóa timer cũ, lập timer mới) thủ công bằng biến số, dễ gây bug rò rỉ timer.
// - Trong Go, ta sử dụng vòng lặp vô hạn kết hợp lệnh **`select`** cực kỳ trực quan và mạnh mẽ ở tầng nhân:
//   - **Case 1: `item := <-b.queue`**: Khi có dữ liệu mới đi vào channel queue. Ta thêm dữ liệu vào slice `batch`.
//     Nếu kích thước slice đạt tới `batchSize`, ta kích hoạt hàm xử lý `processBatchFn` ngay lập tức để giải phóng RAM,
//     đồng thời reset lại Timer.
//   - **Case 2: `<-timer.C`**: Khi hết thời gian chờ tối đa (timeout) mà hàng đợi không gom đủ `batchSize` items.
//     Ta lập tức flush các items đang có lẻ tẻ trong bộ nhớ để bảo đảm độ trễ dữ liệu thấp (Low Latency).
//   - **Case 3: `<-ctx.Done()`**: Cơ chế **Graceful Shutdown** (Tắt ứng dụng an toàn). Khi hệ thống nhận tín hiệu
//     tắt ứng dụng (SIGTERM), Context sẽ hủy bỏ và kích hoạt case này. Ta sẽ thoát vòng lặp, chạy lệnh Flush cuối cùng để
//     ghi nốt các items đang nằm dở dang trong RAM xuống Database, bảo đảm tuyệt đối KHÔNG MẤT MÁT DỮ LIỆU!
func (b *Batcher) Run(ctx context.Context, processBatchFn func([]string)) {
	// TODO:
	// 1. Khởi tạo mảng batch trống với sức chứa batchSize.
	// 2. Tạo một time.Timer hoặc time.Ticker dựa trên b.timeout.
	// 3. Sử dụng vòng lặp vô hạn kết hợp select:
	//    - Case 1: Nhận item từ b.queue -> thêm vào batch. Nếu batch đạt b.batchSize -> gọi processBatchFn, reset batch và timer.
	//    - Case 2: Nhận tín hiệu từ timer -> nếu batch không trống -> gọi processBatchFn, reset batch và timer.
	//    - Case 3: Nhận tín hiệu dừng từ ctx.Done() -> dừng an toàn.
}

func main() {
	// Demo Batcher tại đây
}
