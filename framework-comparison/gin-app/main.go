package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UserDTO định nghĩa dữ liệu đầu vào và các tag validation cho Gin
type UserDTO struct {
	Name  string `json:"name" binding:"required,min=2"`
	Email string `json:"email" binding:"required,email"`
}

func main() {
	// Khởi tạo router Gin trắng (không bao gồm Logger và Recovery mặc định)
	r := gin.New()

	// 1. Custom Middleware Ghi Log (Context Chaining)
	r.Use(func(c *gin.Context) {
		start := time.Now()
		
		// Tiếp tục chạy các middleware tiếp theo hoặc Handler
		c.Next() 
		
		stop := time.Now()

		// In log sau khi đã hoàn thành xử lý
		fmt.Printf("[%s] %d %s %s %s\n",
			c.Request.Method,
			c.Writer.Status(),
			c.Request.URL.Path,
			stop.Sub(start),
			c.ClientIP(),
		)
	})

	// 2. Middleware Panic Recovery giúp server không bị sập khi gặp lỗi runtime nghiêm trọng
	r.Use(gin.Recovery())

	// 3. API Endpoints
	r.POST("/users", func(c *gin.Context) {
		var u UserDTO
		
		// Bind và Validate JSON cùng lúc bằng ShouldBindJSON
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    u,
		})
	})

	fmt.Println("Gin server đang chạy trên cổng 8081...")
	_ = r.Run(":8081")
}
