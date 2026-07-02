package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger là middleware đo lường thời gian xử lý request (Latency) và ghi log.
//
// 🧠 MÔ HÌNH CỦ HÀNH (Onion Model - Gin vs Express):
// - Express: Thường xử lý dạng tuyến tính. Khi cần đo latency, Express không có cơ chế dừng luồng đợi handler chạy xong
//   một cách trực quan, nên thường phải đăng ký callback qua sự kiện `res.on('finish', ...)`.
// - Gin Context: Sở hữu một danh sách các handlers (`HandlersChain`, thực chất là một `[]HandlerFunc`) và một chỉ số `index` đại diện
//   cho middleware hiện tại.
//   - Khi gọi `c.Next()`: Gin sẽ tăng `index` và kích hoạt handler tiếp theo trong chuỗi.
//   - Khi toàn bộ các handler phía sau đã chạy xong, luồng điều khiển sẽ "quay ngược trở lại" (Backtracking) để chạy tiếp phần code nằm
//     PHÍA SAU dòng lệnh `c.Next()`.
//   - Điều này tạo ra mô hình củ hành hoàn hảo: Phần trước `c.Next()` xử lý trước (PRE-processing), phần sau `c.Next()` xử lý sau (POST-processing).
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// PRE-processing: Ghi nhận mốc thời gian bắt đầu khi request vừa bước vào middleware này
		start := time.Now()

		// Chuyển quyền điều khiển cho middleware tiếp theo hoặc Router Handler cuối cùng
		c.Next()

		// POST-processing: Chạy sau khi toàn bộ handler phía sau đã xử lý xong và ghi nhận response.
		// Tính toán thời gian trễ (latency).
		latency := time.Since(start)
		log.Printf("[Logger] %s %s %d %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)
	}
}

// CORSMiddleware cấu hình headers cho phép chia sẻ tài nguyên giữa các nguồn gốc khác nhau (Cross-Origin Resource Sharing).
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Thiết lập CORS Headers trực tiếp vào HTTP Response Writer
		// TODO: Set các header Access-Control-Allow-Origin, Methods, Headers
		//
		// 🧠 CƠ CHẾ HỦY REQUEST (c.Abort() vs Return Gotcha):
		// - Nếu Method là OPTIONS (Preflight Request), ta không muốn chuyển tiếp request này vào các business logic handlers.
		// - Ta gọi `c.AbortWithStatus(200)` hoặc `c.Abort()`.
		// - ⚠️ LƯU Ý QUAN TRỌNG: Gọi `c.Abort()` KHÔNG DỪNG ngay lập tức việc thực thi hàm hiện tại! Nó chỉ thiết lập chỉ số `index`
		//   của chuỗi handlers lên một số cực đại (`abortIndex = 63`) để báo hiệu cho Gin không gọi thêm bất cứ handler nào phía sau.
		// - Do đó, sau khi gọi `c.Abort()`, ta vẫn PHẢI thêm lệnh `return` thủ công nếu muốn dừng ngay lập tức hàm hiện tại,
		//   tránh việc code phía dưới tiếp tục chạy vô ích.
		//
		// Nếu Method là OPTIONS -> c.AbortWithStatus(200) và return!
		c.Next()
	}
}

// AuthMiddleware kiểm tra token/key của client trước khi cho phép đi vào vùng an toàn.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Lấy header "X-API-Key".
		// Nếu rỗng hoặc khác "secret-key" -> gọi c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"}) và return.
		c.Next()
	}
}

func SetupMiddlewareRouter() *gin.Engine {
	// 💡 gin.New() vs gin.Default():
	// - gin.Default() mặc định đính kèm sẵn Logger và Recovery middlewares.
	// - gin.New() khởi tạo một router "trống trơn" (barebone engine), không có bất cứ middleware mặc định nào.
	//   Điều này cực kỳ hữu ích khi ta muốn tùy chỉnh hoàn toàn hệ thống logging hoặc xử lý lỗi tùy biến của doanh nghiệp
	//   mà không muốn bị trùng lặp với middlewares mặc định của Gin.
	r := gin.New()

	// Đăng ký middlewares toàn cục (Global Middleware) cho router.
	// Bất cứ request nào đi vào hệ thống đều phải đi qua RequestLogger và CORSMiddleware đầu tiên.
	r.Use(RequestLogger())
	r.Use(CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	// Đăng ký middleware cục bộ (Route Group Middleware) chỉ áp dụng cho nhóm tuyến đường /api/v1.
	api.Use(AuthMiddleware())
	{
		api.GET("/data", func(c *gin.Context) {
			c.JSON(200, gin.H{"data": "Sensitive data"})
		})
	}

	return r
}

func main() {
	r := SetupMiddlewareRouter()
	r.Run(":8080")
}
