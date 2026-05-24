package usecase

import (
	"errors"
	"go-study/exercises/level-4-gin-production/ex13-clean-architecture/internal/domain"
)

type TodoUsecase struct {
	repo domain.TodoRepository
}

func NewTodoUsecase(repo domain.TodoRepository) *TodoUsecase {
	return &TodoUsecase{repo: repo}
}

func (u *TodoUsecase) CreateTodo(title string) (domain.Todo, error) {
	if title == "" {
		return domain.Todo{}, errors.New("title cannot be empty")
	}
	return u.repo.Create(title)
}

func (u *TodoUsecase) GetTodo(id int) (domain.Todo, error) {
	return u.repo.FindByID(id)
}

func (u *TodoUsecase) ListTodos() ([]domain.Todo, error) {
	return u.repo.FindAll()
}
