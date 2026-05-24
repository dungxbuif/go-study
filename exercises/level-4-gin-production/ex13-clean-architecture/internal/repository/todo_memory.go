package repository

import (
	"sync"

	"go-study/exercises/level-4-gin-production/ex13-clean-architecture/internal/domain"
)

type InMemoryTodoRepository struct {
	mu    sync.RWMutex
	todos map[int]domain.Todo
	next  int
}

func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[int]domain.Todo),
		next:  1,
	}
}

func (r *InMemoryTodoRepository) Create(title string) (domain.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	todo := domain.Todo{ID: r.next, Title: title, Completed: false}
	r.todos[r.next] = todo
	r.next++
	return todo, nil
}

func (r *InMemoryTodoRepository) FindByID(id int) (domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	todo, ok := r.todos[id]
	if !ok {
		return domain.Todo{}, domain.ErrTodoNotFound
	}
	return todo, nil
}

func (r *InMemoryTodoRepository) FindAll() ([]domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]domain.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		list = append(list, todo)
	}
	return list, nil
}
