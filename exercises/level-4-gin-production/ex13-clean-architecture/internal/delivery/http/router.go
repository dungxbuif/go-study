package http

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, h *TodoHandler) {
	r.POST("/todos", h.Create)
	r.GET("/todos/:id", h.Get)
	r.GET("/todos", h.List)
}
