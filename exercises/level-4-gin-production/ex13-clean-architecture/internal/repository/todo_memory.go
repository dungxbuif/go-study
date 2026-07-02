package repository

import (
	"sync"

	"go-study/exercises/level-4-gin-production/ex13-clean-architecture/internal/domain"
)

// InMemoryTodoRepository thực thi interface `domain.TodoRepository` lưu trữ dữ liệu trong bộ nhớ (In-memory).
//
// 🧠 TẦNG REPOSITORY (TẦNG HẠ TẦNG/DATABASE):
// - Đây là tầng Adapter nằm ở rìa ngoài, chịu trách nhiệm lưu trữ và truy vấn dữ liệu từ các phương tiện vật lý.
// - 💡 KIỂU DỮ LIỆU NGẦM ĐỊNH (Implicit/Structural Interfaces in Go):
//   - Khác với TypeScript sử dụng từ khóa `implements ITodoRepository` rõ ràng để báo hiệu cho trình biên dịch,
//     Go áp dụng cơ chế "Duck Typing" (Nếu nó đi như một con vịt và kêu như một con vịt, nó là con vịt).
//   - Struct `InMemoryTodoRepository` chỉ cần định nghĩa đầy đủ 3 hàm: `Create`, `FindByID`, và `FindAll`
//     với đúng signature (tên hàm, tham số, giá trị trả về) là Go sẽ tự động coi struct này implement `domain.TodoRepository`.
//   - Ưu điểm: Tách biệt hoàn toàn, giảm thiểu sự phụ thuộc chéo (decoupling). Tầng Domain hoàn toàn không cần biết struct nào implement nó.
type InMemoryTodoRepository struct {
	mu    sync.RWMutex
	todos map[int]domain.Todo
	next  int
}

// NewInMemoryTodoRepository khởi tạo thực thể InMemoryTodoRepository mới.
func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[int]domain.Todo),
		next:  1,
	}
}

// Create chèn một bản ghi Todo mới vào map.
// Sử dụng Write Lock (`Lock()`) để ngăn chặn tranh chấp ghi khi nhiều Goroutines đồng thời tạo Todo.
func (r *InMemoryTodoRepository) Create(title string) (domain.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	todo := domain.Todo{ID: r.next, Title: title, Completed: false}
	r.todos[r.next] = todo
	r.next++
	return todo, nil
}

// FindByID truy vấn Todo theo ID.
// Sử dụng Read Lock (`RLock()`) cho phép tối ưu truy cập đọc song song cực tốt.
func (r *InMemoryTodoRepository) FindByID(id int) (domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	todo, ok := r.todos[id]
	if !ok {
		return domain.Todo{}, domain.ErrTodoNotFound
	}
	return todo, nil
}

// FindAll lấy ra toàn bộ danh sách Todo hiện có dưới dạng Slice.
func (r *InMemoryTodoRepository) FindAll() ([]domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Khởi tạo một slice mới với capacity tối ưu để tránh phân bổ lại bộ nhớ (re-allocation) liên tục khi append.
	list := make([]domain.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		list = append(list, todo)
	}
	return list, nil
}
