package main

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 🧠 CƠ CHẾ GOM METADATA TRONG HTTP/2 (gRPC Metadata & HTTP/2 under the hood):
// - gRPC không chạy trên HTTP/1.1 thông thường mà sử dụng **HTTP/2** làm giao thức truyền tải lõi.
// - Khác với HTTP/1.1 truyền Headers dạng Text thô cồng kềnh, HTTP/2 truyền tải dưới dạng các khung dữ liệu nhị phân (Binary Frames)
//   được nén chặt bằng thuật toán HPACK.
// - Giao thức gRPC truyền các siêu dữ liệu xác thực (ví dụ: token) thông qua HTTP/2 Headers, gọi là **Metadata**.
// - Trong Go, gRPC server tự động trích xuất các nhãn này từ mạng và nén chúng vào `context.Context`.
// - Ta sử dụng `metadata.FromIncomingContext(ctx)` để lấy ra đối tượng map `metadata.MD` trong bộ nhớ cục bộ.
func UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Printf("[Interceptor] Intercepting gRPC method: %s\n", info.FullMethod)

	// Trích xuất metadata từ Context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is missing")
	}

	// Lấy token xác thực
	tokens := md.Get("authorization")
	if len(tokens) == 0 || tokens[0] != "Bearer valid-token-123" {
		// Trả về mã lỗi gRPC tiêu chuẩn (gRPC Status Codes) thay vì HTTP Status Codes.
		// Dưới nắp capo, mã lỗi `Unauthenticated` (mã 16) sẽ được chuyển dịch thành HTTP/2 RST_STREAM frame
		// gửi về cho client với mã lỗi phù hợp, tránh lộ thông tin hệ thống.
		return nil, status.Errorf(codes.Unauthenticated, "invalid or missing token")
	}

	// Hoàn tất xác thực, chuyển tiếp request đến RPC Handler thực tế
	return handler(ctx, req)
}

// Giả lập RPC Request và Response
type HelloRequest struct {
	Name string
}
type HelloResponse struct {
	Message string
}

// MockHelloHandler giả lập xử lý RPC thực tế sau khi đã qua cửa ải Interceptor
func MockHelloHandler(ctx context.Context, req interface{}) (interface{}, error) {
	request := req.(*HelloRequest)
	return &HelloResponse{
		Message: fmt.Sprintf("Hello %s, your RPC executed successfully!", request.Name),
	}, nil
}

func main() {
	// 🧠 UNARY INTERCEPTORS (gRPC Middleware vs REST Middleware):
	// - Trong REST API (Gin), middleware duyệt qua chuỗi hàm bằng cách gọi `c.Next()`.
	// - Trong gRPC, ta sử dụng **Unary Interceptor**. Signature của hàm interceptor nhận vào `handler grpc.UnaryHandler`.
	// - Interceptor bọc xung quanh Handler thực tế giống như cấu trúc Decorator.
	// - Để cho phép đi tiếp, ta chỉ cần gọi `handler(ctx, req)` và trả về kết quả.
	// - Mô hình này gọn gàng, mang tính hàm số cao (Functional), loại bỏ hoàn toàn các cấu trúc phức tạp hay stateful của HTTP.
	
	// Định nghĩa struct ServerInfo giả lập phương thức RPC
	info := &grpc.UnaryServerInfo{
		FullMethod: "/HelloService/SayHello",
	}

	req := &HelloRequest{Name: "Dung"}

	fmt.Println("--- 1. Testing Unauthenticated Request (expect failure) ---")
	// Tạo context trống không có metadata
	badCtx := context.Background()
	_, err := UnaryAuthInterceptor(badCtx, req, info, MockHelloHandler)
	if err != nil {
		fmt.Printf("[Server] Rejected request: %v\n", err)
	}

	fmt.Println("\n--- 2. Testing Authenticated Request (expect success) ---")
	// Khởi tạo context chứa metadata hợp lệ
	md := metadata.Pairs("authorization", "Bearer valid-token-123")
	goodCtx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := UnaryAuthInterceptor(goodCtx, req, info, MockHelloHandler)
	if err == nil {
		fmt.Printf("[Server] Response message: %s\n", resp.(*HelloResponse).Message)
	}
}
