package domain

import "errors"

// ErrTodoNotFound định nghĩa lỗi nghiệp vụ khi không tìm thấy Todo.
// Đặt lỗi này tại tầng Domain giúp tất cả các tầng khác (Usecase, Repository, Delivery)
// có thể so sánh và xử lý lỗi đồng nhất mà không bị phụ thuộc vào các thư viện bên ngoài.
var ErrTodoNotFound = errors.New("todo not found")

// Todo đại diện cho Domain Entity (Thực thể Nghiệp vụ).
//
// 🧠 TẦNG DOMAIN - TRÁI TIM CỦA HỆ THỐNG:
// - Là tầng lõi trong kiến trúc Clean Architecture.
// - BẮT BUỘC KHÔNG ĐƯỢC PHỤ THUỘC vào bất cứ framework hay thư viện bên ngoài nào (không GORM, không Gin, không Express).
// - Định nghĩa cấu trúc dữ liệu thuần túy phản ánh chính xác các khái niệm của bài toán nghiệp vụ.
type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// TodoRepository định nghĩa cổng giao tiếp (Interface/Port) để tương tác với cơ sở dữ liệu.
//
// 🧠 NGUYÊN LÝ ĐẢO NGƯỢC PHỤ THUỘC (Dependency Inversion Principle - DIP):
// - Thay vì tầng Usecase phụ thuộc trực tiếp vào Repository cụ thể (như GORM, MongoDB),
//   nó sẽ phụ thuộc vào Interface này đặt tại tầng Domain.
// - Tầng Repository cụ thể (nằm ở rìa ngoài) sẽ phải "implements" Interface này.
// - Việc này giúp cô lập hoàn toàn lõi nghiệp vụ khỏi sự biến động của thế giới bên ngoài (swap DB không ảnh hưởng logic).
type TodoRepository interface {
	Create(title string) (Todo, error)
	FindByID(id int) (Todo, error)
	FindAll() ([]Todo, error)
}
