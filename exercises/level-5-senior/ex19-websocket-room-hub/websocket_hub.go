package main

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type WSService struct {
	mu          sync.RWMutex
	connections map[string][]*websocket.Conn            // address -> list conns
	rooms       map[string]map[*websocket.Conn]struct{} // room -> set of conns
}

func NewWSService() *WSService {
	return &WSService{
		connections: make(map[string][]*websocket.Conn),
		rooms:       make(map[string]map[*websocket.Conn]struct{}),
	}
}

func (s *WSService) AddConnection(address string, conn *websocket.Conn) {
	// TODO: Dùng Lock ghi để thêm connection vào slice connections của address
}

func (s *WSService) RemoveConnection(address string, conn *websocket.Conn) {
	// TODO: Dùng Lock ghi để xóa connection ra khỏi connections map và tất cả các rooms
}

func (s *WSService) AddToRoom(room string, conn *websocket.Conn) {
	// TODO: Dùng Lock ghi để thêm connection vào phòng
}

func (s *WSService) BroadcastToRoom(room string, msg interface{}) {
	// TODO: Dùng RLock để đọc danh sách connections trong room,
	// và gửi JSON message song song cho toàn bộ connections đang mở.
}

func main() {
	// Demo WebSocket Hub tại đây
}
