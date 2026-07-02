package http

import (
	"net/http"
	"strconv"

	"go-study/exercises/level-4-gin-production/ex13-clean-architecture/internal/usecase"

	"github.com/gin-gonic/gin"
)

// TodoHandler là HTTP Adapter (Tầng Giao tiếp) chịu trách nhiệm nhận request, trả response.
//
// 🧠 TẦNG DELIVERY (TẦNG GIAO TIẾP/TRANSPORT LAYER):
// - Nằm ở lớp ngoài cùng của Clean Architecture. Nó chỉ phụ thuộc vào `usecase.TodoUsecase`.
// - Tầng này ĐÓNG VAI TRÒ "Phiên dịch viên" (Translator):
//   - Chuyển dữ liệu HTTP đầu vào (JSON body, Path params, Query params) thành các biến thuần túy trong Go.
//   - Gọi hàm xử lý của tầng Usecase.
//   - Nhận kết quả và "dịch" ngược lại thành định dạng HTTP Response chuẩn (JSON, XML, Status Codes).
// - Quyết định cực kỳ quan trọng: TUYỆT ĐỐI không để đối tượng của Gin (`*gin.Context`) đi vào tầng Usecase hay Domain.
//   Nếu ta truyền `*gin.Context` vào Usecase, ta đã phá vỡ hoàn toàn nguyên lý độc lập của Clean Architecture,
//   khiến nghiệp vụ bị trói chặt vào Framework Gin và cực khó để viết unit tests hay chuyển đổi framework.
type TodoHandler struct {
	usecase *usecase.TodoUsecase
}

// NewTodoHandler là constructor khởi tạo thực thể TodoHandler mới.
func NewTodoHandler(u *usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{usecase: u}
}

// Create xử lý HTTP POST /todos.
func (h *TodoHandler) Create(c *gin.Context) {
	// Định nghĩa struct cục bộ (Local request DTO) cho riêng handler này.
	// Rất sạch sẽ, không làm ô nhiễm namespace toàn cục.
	var req struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Gọi lõi nghiệp vụ.
	todo, err := h.usecase.CreateTodo(req.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": todo})
}

// Get xử lý HTTP GET /todos/:id.
func (h *TodoHandler) Get(c *gin.Context) {
	// Nhận và chuẩn hóa dữ liệu từ URL path.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid ID"})
		return
	}

	// Gọi Usecase tìm kiếm dữ liệu.
	todo, err := h.usecase.GetTodo(id)
	if err != nil {
		// Trả về lỗi 404 phù hợp với quy chuẩn RESTful API.
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": todo})
}

// List xử lý HTTP GET /todos.
func (h *TodoHandler) List(c *gin.Context) {
	todos, err := h.usecase.ListTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": todos})
}
