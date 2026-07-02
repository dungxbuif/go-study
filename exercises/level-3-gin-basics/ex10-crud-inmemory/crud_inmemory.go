package main

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Todo biểu diễn cấu trúc dữ liệu cho một mục công việc.
// Các tags như `json:"id"` đóng vai trò siêu dữ liệu (Metadata), hướng dẫn thư viện JSON serializer của Go
// (dưới nắp capo sử dụng cơ chế Reflection - phản chiếu) biết cách đổi tên trường khi chuyển đổi qua lại giữa struct và JSON string.
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// TodoAPI định nghĩa cấu trúc lưu trữ và điều phối nghiệp vụ CRUD trong bộ nhớ.
//
// 🧠 ĐA LUỒNG & AN TOÀN BỘ NHỚ (Go RWMutex vs Node.js Event Loop):
// - Trong Node.js, `todos = new Map()` là an toàn trước lỗi hỏng bộ nhớ (memory corruption) vì mọi hoạt động
//   CRUD chỉ diễn ra trên duy nhất một luồng (Event Loop). Không bao giờ có chuyện 2 dòng code ghi vào Map cùng một micro giây.
// - Trong Go, Gin phục vụ hàng ngàn requests song song bằng hàng ngàn Goroutines chạy trên nhiều nhân CPU thực.
//   Mặc định, `map` của Go KHÔNG an toàn đa luồng (not thread-safe). Nếu 2 Goroutines đồng thời ghi vào map, Go runtime sẽ
//   phát hiện sự cố và lập tức văng lỗi Fatal Panic: `fatal error: concurrent map writes` và dừng ứng dụng ngay lập tức!
// - Giải pháp: Sử dụng `sync.RWMutex` (Read-Write Mutex).
//   - `RLock()` / `RUnlock()`: Nhiều Goroutines có thể ĐỌC dữ liệu cùng lúc.
//   - `Lock()` / `Unlock()`: Chỉ DUY NHẤT một Goroutine được QUYỀN GHI dữ liệu. Khi đang ghi, không ai được phép đọc hay ghi khác.
type TodoAPI struct {
	mu    sync.RWMutex
	todos map[int]Todo
	next  int
}

// NewTodoAPI khởi tạo một thực thể TodoAPI mới nằm trên Heap và trả về con trỏ (Pointer).
// Dưới nắp capo, hàm make(map[int]Todo) phân bổ một cấu trúc bảng băm (hash table) rỗng trên Heap.
func NewTodoAPI() *TodoAPI {
	return &TodoAPI{
		todos: make(map[int]Todo),
		next:  1,
	}
}

// CreateHandler xử lý yêu cầu tạo Todo mới.
//
// 💡 CƠ CHẾ CỦA ShouldBindJSON:
// - Phương thức `c.ShouldBindJSON(&input)` nhận vào một con trỏ trỏ tới struct.
// - Dưới nắp capo, Gin đọc luồng dữ liệu (Stream) từ `c.Request.Body` (vốn là một `io.ReadCloser`),
//   sau đó sử dụng thư viện `encoding/json` để parse và ánh xạ (deserialize) các trường từ JSON vào struct thông qua Pointer.
// - Con trỏ (&) là bắt buộc để hàm binding có thể ghi đè trực tiếp giá trị vào ô nhớ của biến ngoài.
func (api *TodoAPI) CreateHandler(c *gin.Context) {
	// TODO:
	// 1. Nhận JSON body { "title": "..." } bằng ShouldBindJSON
	// 2. Lock ghi (Lock/Unlock) để thêm Todo mới vào map nhằm đảm bảo thread-safety.
	// 3. Trả về status 201 với JSON { "success": true, "data": todo }
}

// ListHandler xử lý yêu cầu lấy danh sách Todo (có bộ lọc completed).
//
// 💡 TỐI ƯU HÓA ĐỌC SONG SONG (Read Lock):
// - Tại đây ta chỉ đọc dữ liệu, do đó sử dụng `api.mu.RLock()` thay vì `Lock()`.
// - Điều này cho phép hàng ngàn requests đọc danh sách Todos chạy song song hoàn toàn không bị block nhau,
//   nâng cao thông lượng xử lý của hệ thống (High Throughput).
func (api *TodoAPI) ListHandler(c *gin.Context) {
	// TODO:
	// 1. Đọc RLock để lấy danh sách Todos từ map an toàn
	// 2. Hỗ trợ query completed=true/false
	// 3. Trả về status 200 với JSON { "success": true, "data": []todo }
}

// GetHandler tìm kiếm một Todo theo ID được truyền qua URL path param.
func (api *TodoAPI) GetHandler(c *gin.Context) {
	// TODO: Tìm todo theo id.
	// Sử dụng api.mu.RLock() để bảo vệ quá trình đọc map.
	// Nếu không thấy -> Trả về 404 { "success": false, "error": "Todo not found" }
}

// UpdateHandler cập nhật thông tin Todo.
// Yêu cầu sử dụng Write Lock (`Lock()`) vì thao tác này thay đổi trực tiếp ô nhớ của Todo trong map.
func (api *TodoAPI) UpdateHandler(c *gin.Context) {
	// TODO: Cập nhật title hoặc completed của todo theo id. Đảm bảo dùng Lock ghi.
}

// DeleteHandler xóa Todo ra khỏi bộ nhớ.
// Thao tác xóa (`delete(map, key)`) là một thao tác cấu trúc lại bảng băm bên dưới, 
// bắt buộc phải được bảo vệ nghiêm ngặt bằng Write Lock để tránh tranh chấp bộ nhớ.
func (api *TodoAPI) DeleteHandler(c *gin.Context) {
	// TODO: Xóa todo theo id khỏi map. Đảm bảo dùng Lock ghi.
}

func SetupTodoRouter(api *TodoAPI) *gin.Engine {
	r := gin.Default()

	r.POST("/todos", api.CreateHandler)
	r.GET("/todos", api.ListHandler)
	r.GET("/todos/:id", api.GetHandler)
	r.PUT("/todos/:id", api.UpdateHandler)
	r.DELETE("/todos/:id", api.DeleteHandler)

	return r
}

func main() {
	api := NewTodoAPI()
	r := SetupTodoRouter(api)
	r.Run(":8080")
}
