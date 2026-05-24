package main

import (
	"github.com/gin-gonic/gin"
	"go-study/exercises/level-4-gin-production/ex13-clean-architecture/internal/delivery/http"
	"go-study/exercises/level-4-gin-production/ex13-clean-architecture/internal/repository"
	"go-study/exercises/level-4-gin-production/ex13-clean-architecture/internal/usecase"
)

func main() {
	r := gin.Default()

	todoRepo := repository.NewInMemoryTodoRepository()
	todoUsecase := usecase.NewTodoUsecase(todoRepo)
	todoHandler := http.NewTodoHandler(todoUsecase)

	http.SetupRoutes(r, todoHandler)

	r.Run(":8080")
}
