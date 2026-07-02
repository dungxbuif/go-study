# 🎯 Go Backend Live Coding & Interview Requirements

Khi phỏng vấn vị trí Go Backend Developer, các nhà tuyển dụng thường tập trung vào khả năng thiết kế hệ thống có tính mở rộng cao, xử lý đồng thời (concurrency), và viết code idiomatic Go (sạch, dễ test, quản lý bộ nhớ tốt).

Dưới đây là danh sách tổng hợp các đề bài Live Coding và System Design phổ biến nhất dựa trên quá trình research.

## 1. Các Chủ Đề Cốt Lõi (Core Areas)

- **Concurrency & Goroutines:** Quản lý vòng đời goroutine, tránh rò rỉ (leak), sử dụng `sync.WaitGroup`, `sync.Mutex`, và Channels.
- **Data Structures:** Hiểu rõ cách Slice tự động tăng dung lượng, sự khác biệt giữa `make` và `new`, thao tác với map an toàn trong môi trường multi-thread (`sync.Map` hoặc Mutex + Map).
- **Thiết kế API & Middleware:** Tự viết REST API sử dụng package standard `net/http` hoặc các framework như Gin/Echo. Khả năng viết custom middleware (Logging, Auth, Rate Limit).
- **Error Handling:** Xử lý lỗi theo phong cách Go (trả về lỗi như một giá trị), wrap lỗi với `fmt.Errorf`, và thực hiện Graceful Shutdown.
- **Testing:** Viết unit tests (`_test.go`), sử dụng `httptest` để mock API, dùng Table-driven tests.

---

## 2. Danh Sách Đề Bài Live Coding Thực Chiến

### Bài 1: Rate Limiter Middleware

Bảo vệ API khỏi bị lạm dụng bằng cách giới hạn số lượng request từ một IP.

- **Yêu cầu cơ bản:** Cài đặt thuật toán Token Bucket hoặc Fixed Window trên memory (sử dụng Map kết hợp với Mutex).
- **Nâng cao:** Tích hợp Redis để chạy trên môi trường phân tán (Distributed System). Viết unit test cho middleware dùng `httptest`.
- **Kỹ năng kiểm tra:** Kiến trúc Middleware trong Go, Concurrency-safe data structures, HTTP context.

### Bài 2: URL Shortener (System Design & Code)

Bài toán kinh điển đánh giá khả năng thiết kế hệ thống bất đối xứng (Read/Write asymmetry).

- **Yêu cầu cơ bản:** Viết API `POST /shorten` nhận URL dài trả về code ngắn (thường dùng hash base62). API `GET /:code` redirect về URL gốc (HTTP 301/302).
- **Nâng cao:** Thêm Redis Cache (Read-through) để giảm tải Database. Cấu trúc project rành mạch theo Clean Architecture hoặc Layered (Handler -> Service -> Repository).
- **Kỹ năng kiểm tra:** Giao tiếp Database, Caching strategies, Interface design để dễ dàng thay đổi storage.

### Bài 3: Worker Pool & Task Queue

Thiết kế hệ thống xử lý hàng ngàn background jobs một cách an toàn.

- **Yêu cầu cơ bản:** Cho một lượng lớn jobs, cấu hình số lượng worker cố định (ví dụ 100 workers) để xử lý mà không làm quá tải CPU/RAM.
- **Nâng cao:** Có khả năng Graceful Shutdown thông qua `context.Context` khi nhận tín hiệu SIGINT/SIGTERM từ OS (tự động ngừng nhận job mới và đợi job cũ xong).
- **Kỹ năng kiểm tra:** Channels (Buffered vs Unbuffered), Goroutines lifecycle, Context cancellation, WaitGroup.

### Bài 4: Thread-Safe LRU Cache (In-Memory)

Viết một bộ đệm cache có giới hạn kích thước với chính sách loại bỏ "Ít được sử dụng gần đây nhất" (Least Recently Used).

- **Yêu cầu cơ bản:** Triển khai các hàm `Get(key)`, `Put(key, value)`. Khi dung lượng đầy, tự động xoá phần tử cũ nhất.
- **Nâng cao:** Đảm bảo Thread-safe trong môi trường chạy nhiều goroutines (sử dụng Mutex `sync.RWMutex`), tối ưu hiệu năng O(1) (kết hợp Doubly Linked List và Hash Map).
- **Kỹ năng kiểm tra:** Cấu trúc dữ liệu nâng cao, Pointers, Concurrency synchronization.

### Bài 5: Xử Lý File Dung Lượng Lớn (Data Pipeline)

Đọc và xử lý file CSV hoặc Log dung lượng lớn (VD: 10GB) với bộ nhớ RAM giới hạn (VD: 100MB).

- **Yêu cầu cơ bản:** Stream dữ liệu dùng `bufio.Scanner`, đẩy từng dòng vào channel, tập hợp các worker lấy ra xử lý và tổng hợp kết quả cuối cùng.
- **Nâng cao:** Áp dụng Fan-out, Fan-in pattern. Đảm bảo luồng dữ liệu không bị chặn (deadlocks) ở các bước.
- **Kỹ năng kiểm tra:** I/O stream, Concurrency pipeline, Memory management, Garbage Collection awareness.

---

## 3. Lời Khuyên "Ăn Điểm" Khi Live Coding Bằng Go

1. **Keep it Simple (MVP First):** Luôn ưu tiên viết bản chạy được (MVP) với In-Memory storage (Dùng Map) trước khi nhắc tới hay setup Redis/Kafka.
2. **Handle Errors Cẩn Thận:** Đừng bao giờ dùng `_` để bỏ qua lỗi trừ khi giải thích được lý do chính đáng. Interviewer đặc biệt quan tâm cách bạn xử lý các rẽ nhánh lỗi.
3. **Nằm lòng triết lý Go:** _"Do not communicate by sharing memory; instead, share memory by communicating"_. Hãy cân nhắc việc dùng Channel để truyền dữ liệu trước khi chọn cách lock vùng nhớ bằng Mutex.
4. **Viết Test:** Nếu hoàn thành sớm, hãy chủ động viết một file `_test.go` dùng `t.Run()` (Table-driven tests). Điều này thể hiện tư duy của một Senior/Professional Developer.
