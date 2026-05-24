package main

import (
	"context"
	"time"
)

type Batcher struct {
	queue     chan string
	batchSize int
	timeout   time.Duration
}

func NewBatcher(batchSize int, timeout time.Duration) *Batcher {
	return &Batcher{
		queue:     make(chan string, 10000),
		batchSize: batchSize,
		timeout:   timeout,
	}
}

func (b *Batcher) Push(item string) {
	b.queue <- item
}

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
