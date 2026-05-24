package domain

import "errors"

var ErrTodoNotFound = errors.New("todo not found")

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoRepository interface {
	Create(title string) (Todo, error)
	FindByID(id int) (Todo, error)
	FindAll() ([]Todo, error)
}
