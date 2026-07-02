package main

import (
	"context"
	"fmt"
	"net/http"
)

// contextKey là kiểu dữ liệu tùy chỉnh cho các keys của context để tránh đụng độ khóa (Key Collision) giữa các gói thư viện khác nhau.
type contextKey string

const traceIDKey contextKey = "trace_id"

// WithTraceID lưu trữ TraceID vào trong Context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	// Dưới nắp capo, context trong Go là bất biến (immutable).
	// Hàm context.WithValue tạo ra một child context mới bọc xung quanh parent context,
	// hình thành một cấu trúc cây tìm kiếm ngược (linked list/tree).
	// Khi tìm kiếm một key, Go sẽ đi ngược từ nút lá hiện tại lên nút gốc để tìm.
	return context.WithValue(ctx, traceIDKey, traceID)
}

// GetTraceID lấy TraceID từ Context ra.
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if val, ok := ctx.Value(traceIDKey).(string); ok {
		return val
	}
	return ""
}

// InjectTraceHeader ghi thông tin TraceID vào HTTP Headers của request gửi đi.
//
// 🧠 CƠ CHẾ KHÁCH QUAN TRONG PHÂN TÁN (Distributed Context Propagation under the hood):
// - Context của Go chỉ có giá trị trong bộ nhớ của MỘT tiến trình ứng dụng đơn lẻ (Local Memory Space).
// - Khi ứng dụng thực hiện cuộc gọi HTTP (gRPC) sang một service khác, context của Go sẽ KHÔNG tự động bay qua mạng.
// - Để giải quyết, ta sử dụng giao thức **W3C Trace Context Standard** (chuẩn hóa quốc tế):
//   - Ta biến đổi TraceID trong bộ nhớ thành một chuỗi Header HTTP, ví dụ: `traceparent: 00-[trace_id]-[span_id]-01`.
//   - Ghi Header này vào Outgoing HTTP Request.
//   - Ở phía bên nhận, server phụ trách xử lý sẽ đọc HTTP Header này, giải mã (parse) và tái tạo lại một
//     `context.Context` mới tinh trong RAM của nó chứa TraceID đó.
// - Cơ chế này kết nối các bản log rời rạc của hàng trăm microservices độc lập thành một chuỗi vòng đời request liền mạch!
func InjectTraceHeader(ctx context.Context, req *http.Request) {
	traceID := GetTraceID(ctx)
	if traceID != "" {
		// Ghi theo định dạng tiêu chuẩn W3C traceparent
		req.Header.Set("traceparent", fmt.Sprintf("00-%s-0000000000000000-01", traceID))
	}
}

// ExtractTraceHeader trích xuất TraceID từ HTTP Headers của request đi vào.
func ExtractTraceHeader(req *http.Request) string {
	traceParent := req.Header.Get("traceparent")
	if len(traceParent) >= 35 {
		// Format: 00-trace_id-span_id-flags (00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01)
		// Trích xuất phần trace_id
		return traceParent[3:35]
	}
	return ""
}

func main() {
	// Root Context đóng vai trò khởi nguồn
	rootCtx := context.Background()

	// Khởi tạo một Trace ID giả lập
	traceID := "4bf92f3577b34da6a3ce929d0e0e4736"
	ctxWithTrace := WithTraceID(rootCtx, traceID)

	fmt.Printf("[Service A] Starting request. Created TraceID: %s\n", traceID)

	// Tạo một request HTTP giả lập để gọi sang Service B
	req, _ := http.NewRequest("GET", "http://service-b/api/data", nil)

	// Inject Trace ID vào Outgoing Headers
	InjectTraceHeader(ctxWithTrace, req)
	fmt.Printf("[Service A] Injected traceparent header: %s\n", req.Header.Get("traceparent"))

	// --- Giả lập cuộc gọi qua mạng đến Service B ---
	fmt.Println("\n--- Simulating Network Transmission to Service B ---")

	// Service B nhận request và trích xuất Trace ID
	receivedTraceID := ExtractTraceHeader(req)
	serviceBCtx := WithTraceID(context.Background(), receivedTraceID)

	fmt.Printf("[Service B] Received request. Extracted TraceID: %s\n", GetTraceID(serviceBCtx))
	fmt.Println("[Service B] Querying Database... TraceID is automatically attached to DB queries and logs.")
}
