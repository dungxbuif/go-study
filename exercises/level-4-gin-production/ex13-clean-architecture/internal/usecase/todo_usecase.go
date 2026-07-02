package usecase

import (
	"errors"
	"go-study/exercises/level-4-gin-production/ex13-clean-architecture/internal/domain"
)

// TodoUsecase quản lý và thực hiện các quy trình nghiệp vụ (Use Cases/Business Rules) của Todo.
//
// 🧠 TẦNG USECASE (TẦNG NGHIỆP VỤ):
// - Chứa các logic tính toán nghiệp vụ cốt lõi (ví dụ: title không được rỗng, phân quyền, kiểm tra điều kiện).
// - Chỉ phụ thuộc vào `domain.TodoRepository` dưới dạng một interface trừu tượng.
// - Tầng Usecase hoàn toàn không biết và không quan tâm dữ liệu được lưu ở đâu (Memory, PostgreSQL, File)
//   hay request đến từ giao diện nào (HTTP Gin, gRPC, CLI).
// - Nhờ sự độc lập này, code nghiệp vụ có thể được viết Unit Tests cực kỳ dễ dàng bằng cách mock interface `TodoRepository`.
type TodoUsecase struct {
	repo domain.TodoRepository
}

// NewTodoUsecase là Constructor tạo ra một thực thể TodoUsecase mới.
// Nhận vào interface `domain.TodoRepository`. Bất cứ struct nào implement đủ các phương thức của interface này đều có thể truyền vào đây.
func NewTodoUsecase(repo domain.TodoRepository) *TodoUsecase {
	return &TodoUsecase{repo: repo}
}

// CreateTodo kiểm tra nghiệp vụ và ra lệnh tạo mới Todo.
func (u *TodoUsecase) CreateTodo(title string) (domain.Todo, error) {
	// Logic nghiệp vụ (Business Validation):
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
