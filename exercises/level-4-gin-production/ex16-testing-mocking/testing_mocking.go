package main

import "errors"

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoRepository interface {
	Create(title string) (Todo, error)
}

type TodoUsecase struct {
	repo TodoRepository
}

func NewTodoUsecase(repo TodoRepository) *TodoUsecase {
	return &TodoUsecase{repo: repo}
}

func (u *TodoUsecase) CreateTodo(title string) (Todo, error) {
	// TODO:
	// 1. Kiểm tra nếu title rỗng -> Trả về lỗi "title is required"
	// 2. Gọi u.repo.Create(title) để ghi vào DB
	return Todo{}, errors.New("not implemented")
}

func main() {
	// Demo testing & mocking tại đây
}
