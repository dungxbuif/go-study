package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO:
		// 1. Ghi lại thời gian bắt đầu
		// 2. Gọi c.Next() để chuyển sang handler tiếp theo
		// 3. Đo lường thời gian chạy và in ra log method, path, status, latency
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log.Printf("[Logger] %s %s %d %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Set các header Access-Control-Allow-Origin, Methods, Headers
		// Nếu Method là OPTIONS -> c.AbortWithStatus(200)
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Lấy header "X-API-Key".
		// Nếu rỗng hoặc khác "secret-key" -> c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		c.Next()
	}
}

func SetupMiddlewareRouter() *gin.Engine {
	r := gin.New() // gin.New() không chứa logger và recovery mặc định

	r.Use(RequestLogger())
	r.Use(CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
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
