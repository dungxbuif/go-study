package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-study/exercises/level-4-gin-production/ex13-clean-architecture/internal/usecase"
)

type TodoHandler struct {
	usecase *usecase.TodoUsecase
}

func NewTodoHandler(u *usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{usecase: u}
}

func (h *TodoHandler) Create(c *gin.Context) {
	var req struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	todo, err := h.usecase.CreateTodo(req.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": todo})
}

func (h *TodoHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid ID"})
		return
	}

	todo, err := h.usecase.GetTodo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": todo})
}

func (h *TodoHandler) List(c *gin.Context) {
	todos, err := h.usecase.ListTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": todos})
}
