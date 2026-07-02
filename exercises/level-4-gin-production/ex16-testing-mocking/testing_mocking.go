package main

import "errors"

// Todo đại diện cho cấu trúc Todo.
type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// TodoRepository định nghĩa cổng lưu trữ để Usecase tương tác.
// Tách rời giao tiếp qua interface là bước đi sống còn để thực hiện Unit Test.
type TodoRepository interface {
	Create(title string) (Todo, error)
}

// TodoUsecase chứa logic nghiệp vụ.
type TodoUsecase struct {
	repo TodoRepository
}

// NewTodoUsecase là Constructor nhận interface TodoRepository.
func NewTodoUsecase(repo TodoRepository) *TodoUsecase {
	return &TodoUsecase{repo: repo}
}

// CreateTodo kiểm tra nghiệp vụ và ra lệnh lưu trữ.
//
// 🧠 PHƯƠNG PHÁP MOCKING TĨNH TRONG GOLANG (Interface-based Mocking vs JS Jest Mock):
// - Trong Node.js/Jest, ta có thể dùng `jest.mock('./repository')` để ghi đè (patch) động bất cứ module nào.
// - Trong Go, code được biên dịch tĩnh trực tiếp sang mã máy. Ta KHÔNG THỂ hack hay khống chế thời gian thực thi (runtime) kiểu JS/Jest.
// - Cách duy nhất để mock là sử dụng **Interface**.
// - Lúc viết Test (file `_test.go`), ta sẽ tự tạo ra một struct giả lập (ví dụ: `type MockTodoRepository struct`)
//   chứa các trường mô phỏng (ví dụ: `shouldFail bool`, `called bool`) và triển khai đầy đủ interface `TodoRepository`.
//   Sau đó ta truyền (inject) Mock struct này vào `NewTodoUsecase`. Cực kỳ tường minh và an toàn kiểu dữ liệu (Type-safe)!
func (u *TodoUsecase) CreateTodo(title string) (Todo, error) {
	// TODO:
	// 1. Kiểm tra nếu title rỗng -> Trả về lỗi "title is required"
	// 2. Gọi u.repo.Create(title) để ghi vào DB
	return Todo{}, errors.New("not implemented")
}

// 🧠 TABLE-DRIVEN TESTS & RACE DETECTOR (Kiểm thử theo bảng và máy phát hiện xung đột):
// - Go có một mô thức viết test kinh điển gọi là **Table-Driven Tests**.
// - Thay vì viết nhiều hàm test rời rạc (`testCreateSuccess`, `testCreateEmptyTitle`), ta định nghĩa một mảng các Struct vô danh (Anonymous Structs):
//   ```go
//   tests := []struct {
//       name          string
//       inputTitle    string
//       mockFail      bool
//       expectedErr   string
//   }{
//       {"Success Case", "Task 1", false, ""},
//       {"Empty Title", "", false, "title is required"},
//   }
//   ```
// - Ta lặp qua mảng này và chạy `t.Run(tc.name, func(t *testing.T) { ... })`. Code kiểm thử chỉ cần viết đúng 1 lần nhưng bao phủ được mọi trường hợp!
// - 💡 Cờ phát hiện Race Condition (`go test -race ./...`):
//   - Đây là công cụ "quỷ khóc thần sầu" được tích hợp sẵn vào trình biên dịch Go.
//   - Khi bật `-race`, Go compiler sẽ chèn thêm mã giám sát ô nhớ vào file nhị phân test.
//   - Nếu có 2 goroutines tranh chấp đọc-ghi đồng thời vào cùng một ô nhớ mà không có Lock, Race Detector sẽ phát hiện và in log cảnh báo chi tiết từng dòng code.
//   - Node.js không có và cũng không cần công cụ này do Event Loop đơn luồng, nhưng trong Go đây là vũ khí tối tân để bảo vệ hệ thống đa luồng.
func main() {
	// Demo testing & mocking tại đây
}
