package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

// WSService định nghĩa cấu trúc Hub để quản lý các kết nối WebSocket.
//
// 🧠 CƠ CHẾ UPGRADE KẾT NỐI (HTTP Hijacking under the hood):
// - Giao thức WebSocket ban đầu khởi xuất từ một request HTTP GET bình thường kèm theo header `Upgrade: websocket`.
// - Dưới nắp capo, thư viện `gorilla/websocket` (hoặc gin's upgrader) sẽ thực hiện **Connection Hijacking** (Cướp kết nối):
//   1. Trình nâng cấp (Upgrader) sẽ chiếm quyền kiểm soát socket TCP thô (`net.Conn`) trực tiếp từ HTTP Server của Go.
//   2. Nó trả về mã trạng thái HTTP `101 Switching Protocols` cho Client.
//   3. Thay vì đóng socket khi kết thúc hàm handler như HTTP thông thường, Go giữ kết nối TCP socket này mở vĩnh viễn,
//      cho phép truyền và nhận dữ liệu bất đồng bộ 2 chiều (Full-Duplex).
type WSService struct {
	mu          sync.RWMutex
	connections map[string][]*websocket.Conn            // address/user -> danh sách các conns đang mở
	rooms       map[string]map[*websocket.Conn]struct{} // room -> set (tập hợp) các conns đang tham gia phòng
}

func NewWSService() *WSService {
	return &WSService{
		connections: make(map[string][]*websocket.Conn),
		rooms:       make(map[string]map[*websocket.Conn]struct{}),
	}
}

// AddConnection thêm kết nối WebSocket mới vào Hub.
// Thao tác này ghi đè lên map, do đó bắt buộc phải sử dụng Write Lock (`Lock()`).
// Nếu 2 Goroutines của 2 users kết nối đồng thời và ta không Lock, map connections sẽ bị hỏng cấu trúc ô nhớ, gây panic sập server.
func (s *WSService) AddConnection(address string, conn *websocket.Conn) {
	// TODO: Dùng Lock ghi để thêm connection vào slice connections của address
}

// RemoveConnection ngắt kết nối và dọn dẹp bộ nhớ của Connection.
// ⚠️ LƯU Ý RÒ RỈ GOROUTINE (Goroutine/Memory Leak):
// - Khi một client tắt trình duyệt hoặc mất sóng đột ngột, TCP connection sẽ bị ngắt (Broken Pipe).
// - Nếu ta không có cơ chế bắt sự kiện lỗi lúc đọc/ghi và không dọn dẹp connection khỏi maps,
//   đối tượng con trỏ `*websocket.Conn` sẽ nằm lì trong RAM vĩnh viễn và không bao giờ được GC thu hồi.
// - Tệ hơn, goroutine chuyên trách đọc socket đó sẽ bị treo vô hạn (blocking read), gây rò rỉ hàng vạn Goroutines trong production!
func (s *WSService) RemoveConnection(address string, conn *websocket.Conn) {
	// TODO: Dùng Lock ghi để xóa connection ra khỏi connections map và tất cả các rooms
}

func (s *WSService) AddToRoom(room string, conn *websocket.Conn) {
	// TODO: Dùng Lock ghi để thêm connection vào phòng
}

// BroadcastToRoom thực hiện phát sóng tin nhắn đến toàn bộ thành viên trong phòng.
//
// 🧠 TỐI ƯU HÓA PHÁT TIN NHẮN (Read Lock & Concurrency):
// - Khi Broadcast, ta không thay đổi (ghi) cấu trúc của map `rooms`, ta chỉ DUYỆT qua danh sách connections đang có (đọc).
// - Do đó, sử dụng `RLock()` để cho phép nhiều luồng phát sóng tin nhắn song song xuyên suốt các phòng khác nhau
//   mà không gây nghẽn luồng.
// - 💡 Mẹo nâng cao: Để tránh việc một kết nối socket bị nghẽn mạng làm chậm trễ cả phòng, ta nên kích hoạt
//   việc gửi tin nhắn `conn.WriteJSON()` bên trong các Goroutines riêng lẻ cho từng kết nối trong phòng!
func (s *WSService) BroadcastToRoom(room string, msg interface{}) {
	// TODO: Dùng RLock để đọc danh sách connections trong room,
	// và gửi JSON message song song cho toàn bộ connections đang mở.
}

func main() {
	// Demo WebSocket Hub tại đây
}
