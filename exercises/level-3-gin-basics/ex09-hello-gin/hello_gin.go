package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetupRouter cấu hình và khởi tạo các router của Gin.
//
// 🧠 SỰ KHÁC BIỆT CỐT LÕI VỀ KIẾN TRÚC ROUTING (Gin vs Express):
// - Express (Node.js): Sử dụng cấu trúc danh sách tuyến tính (linear array) kết hợp Regular Expressions.
//   Mỗi request đến sẽ duyệt tuần tự qua từng route cho tới khi khớp. Nếu ứng dụng có hàng trăm routes,
//   hiệu năng route-matching sẽ bị giảm đáng kể (độ phức tạp O(N)).
// - Gin (Go): Sử dụng cấu trúc cây Radix Tree (một dạng Trie tối ưu). Việc tìm kiếm route chỉ phụ thuộc
//   vào độ dài của URL path (độ phức tạp O(K) với K là độ dài chuỗi), hoàn toàn độc lập với tổng số lượng routes.
//   Điều này giúp Gin đạt tốc độ định tuyến siêu tốc, gần như không tốn chi phí CPU.
func SetupRouter() *gin.Engine {
	// gin.Default() khởi tạo một *gin.Engine đã tích hợp sẵn 2 middleware mặc định dưới nắp capo:
	// 1. Logger: Ghi lại thông tin chi tiết của mỗi request (Method, Path, Status Code, Latency...).
	// 2. Recovery: Bắt (recover) mọi trường hợp xảy ra 'panic' trong quá trình xử lý request. 
	//    Thay vì làm sập toàn bộ tiến trình HTTP server (như Node.js nếu gặp uncaught exception mà không có xử lý đặc biệt),
	//    Recovery middleware của Gin sẽ tự động log stack trace và trả về HTTP 500 Internal Server Error, giúp ứng dụng duy trì hoạt động bền bỉ.
	r := gin.Default()

	// r.GET("/", ...) đăng ký route GET tại root.
	// gin.Context đại diện cho ngữ cảnh của request hiện tại. Nó hợp nhất cả Request và Response (tương đương req, res của Express).
	// c.JSON(200, gin.H{...}) thực hiện serialize map sang JSON.
	// gin.H là viết tắt của: map[string]interface{}. Dưới nắp capo, map trong Go được phân bổ trên Heap và là kiểu tham chiếu (reference type).
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, Go!"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// GET /users/:id biểu diễn path parameter (Route Parameters).
	// c.Param("id") lấy giá trị động từ URL. Dưới nắp capo, Gin đã phân tích cú pháp và lưu các tham số này
	// trong một mảng slices để tránh việc phân bổ bộ nhớ dư thừa trong quá trình xử lý chuỗi.
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(200, gin.H{"user_id": id})
	})

	// GET /search?q=keyword&page=2 biểu diễn query parameter.
	// c.DefaultQuery("q", "") lấy giá trị query string "q", nếu trống sẽ trả về mặc định "".
	// Vì tất cả query params trả về từ Gin đều là kiểu String, ta bắt buộc phải parse thủ công (Ví dụ: strconv.Atoi)
	// để chuyển sang kiểu số (Static typing), khác với Node.js tự động chuyển đổi linh hoạt.
	r.GET("/search", func(c *gin.Context) {
		q := c.DefaultQuery("q", "")
		pageStr := c.DefaultQuery("page", "1")
		
		// Chuyển chuỗi sang số nguyên. Dưới nắp capo, strconv.Atoi tối ưu hiệu năng
		// bằng cách lướt qua mảng byte trực tiếp, hạn chế cấp phát bộ nhớ.
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
		c.JSON(200, gin.H{"query": q, "page": page})
	})

	// Router Group giúp gom nhóm các API có chung tiền tố (Prefix).
	// Tương tự Express Router, r.Group("/api/v1") cho phép phân chia module,
	// áp dụng middleware chung cho toàn bộ group đó cực kỳ sạch sẽ và khoa học.
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	}

	return r
}

func main() {
	r := SetupRouter()
	
	// r.Run(":8080") lắng nghe và phục vụ HTTP requests trên cổng 8080.
	// Dưới nắp capo, hàm này gọi http.ListenAndServe(":8080", r).
	// Go Net HTTP package sẽ tạo ra 1 Goroutine RIÊNG BIỆT cho MỖI request kết nối đến.
	// Do đó, máy chủ có thể tận dụng toàn bộ số nhân CPU vật lý mà không gặp hiện tượng nghẽn luồng đơn như Node.js Event Loop!
	r.Run(":8080")
}
