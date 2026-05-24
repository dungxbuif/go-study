package main

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoAPI struct {
	mu    sync.RWMutex
	todos map[int]Todo
	next  int
}

func NewTodoAPI() *TodoAPI {
	return &TodoAPI{
		todos: make(map[int]Todo),
		next:  1,
	}
}

func (api *TodoAPI) CreateHandler(c *gin.Context) {
	// TODO:
	// 1. Nhận JSON body { "title": "..." } bằng ShouldBindJSON
	// 2. Lock ghi (Lock/Unlock) để thêm Todo mới vào map
	// 3. Trả về status 201 với JSON { "success": true, "data": todo }
}

func (api *TodoAPI) ListHandler(c *gin.Context) {
	// TODO:
	// 1. Đọc RLock để lấy danh sách Todos từ map
	// 2. Hỗ trợ query completed=true/false
	// 3. Trả về status 200 với JSON { "success": true, "data": []todo }
}

func (api *TodoAPI) GetHandler(c *gin.Context) {
	// TODO: Tìm todo theo id. Nếu không thấy -> Trả về 404 { "success": false, "error": "Todo not found" }
}

func (api *TodoAPI) UpdateHandler(c *gin.Context) {
	// TODO: Cập nhật title hoặc completed của todo theo id. Đảm bảo dùng Lock ghi.
}

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
