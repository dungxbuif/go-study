package main

import (
	"errors"

	"gorm.io/gorm"
)

// User biểu diễn bảng người dùng trong Database.
//
// 🧠 EMBEDDED STRUCT & SOFT DELETE (gorm.Model under the hood):
// - `gorm.Model` là một struct cơ sở chứa 4 trường: `ID` (khóa chính), `CreatedAt`, `UpdatedAt`, và `DeletedAt`.
// - Khi nhúng `gorm.Model` vào struct User, Go kế thừa (composition) toàn bộ các trường này.
// - 💡 Cơ chế SOFT DELETE: Khi trường `DeletedAt` kiểu `gorm.DeletedAt` được khai báo, GORM sẽ thay đổi hành vi xóa mặc định.
//   Lệnh `db.Delete(&user)` sẽ KHÔNG chạy `DELETE FROM users ...` mà thay vào đó chạy lệnh `UPDATE users SET deleted_at = CURRENT_TIMESTAMP ...`.
//   Mọi truy vấn tìm kiếm tiếp theo của GORM sẽ tự động chèn thêm điều kiện `WHERE deleted_at IS NULL` dưới nắp capo để ẩn các bản ghi đã xóa,
//   giúp khôi phục dữ liệu dễ dàng và đảm bảo tính toàn vẹn tham chiếu.
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:50"`
	Posts    []Post `gorm:"foreignKey:AuthorID"` // Thiết lập quan hệ 1-N (One-to-Many)
}

// Post biểu diễn bảng bài viết.
type Post struct {
	gorm.Model
	Title    string `gorm:"size:200"`
	Content  string `gorm:"type:text"`
	AuthorID uint   // Khóa ngoại liên kết tới User.ID
	Author   User   // Thiết lập quan hệ 1-1 (BelongsTo) để thực hiện Preload
}

// SetupDB thiết lập kết nối và cấu hình Connection Pool.
//
// 🧠 CONNECTION POOLING (database/sql under the hood):
// - GORM thực chất là một lớp bọc (wrapper) phía trên gói thư viện tiêu chuẩn `database/sql` của Go.
// - Thư viện `database/sql` quản lý một **Connection Pool** (Hồ chứa kết nối) đa luồng cực kỳ mạnh mẽ.
// - Khi ta thực hiện truy vấn, database/sql sẽ mượn (acquire) một kết nối nhàn rỗi (idle connection) trong pool,
//   thực thi lệnh SQL, rồi lập tức trả kết nối đó về pool để các goroutines khác tái sử dụng.
// - Ta cấu hình Connection Pool qua 3 tham số sống còn:
//   1. `SetMaxOpenConns(N)`: Số lượng kết nối tối đa được mở đồng thời tới DB. Ngăn chặn việc làm sập DB do mở quá nhiều kết nối từ server.
//   2. `SetMaxIdleConns(M)`: Số lượng kết nối tối đa được giữ ở trạng thái nhàn rỗi. Giảm thời gian chờ thiết lập lại TCP Handshake.
//   3. `SetConnMaxLifetime(T)`: Thời gian tồn tại tối đa của một kết nối để tránh lỗi rò rỉ hoặc ngắt kết nối do tường lửa DB.
func SetupDB() (*gorm.DB, error) {
	// TODO: Khởi tạo kết nối SQLite in-memory bằng GORM
	// Tự động migrate: db.AutoMigrate(&User{}, &Post{})
	return nil, errors.New("not implemented")
}

// UpsertPost thực hiện chèn dữ liệu mới hoặc cập nhật nếu đã tồn tại bản ghi (UPSERT).
//
// 🧠 TÍNH NGUYÊN TỬ (SQL Upsert vs Select-then-Update Race Condition):
// - Nếu ta viết logic: (1) Kiểm tra xem Post đã tồn tại chưa -> (2) Nếu chưa thì Insert, nếu có thì Update.
//   Trong môi trường đa luồng (Goroutines), 2 luồng đồng thời chạy bước (1) đều thấy chưa có, dẫn đến cả 2 cùng chạy (2) Insert,
//   gây ra lỗi trùng lặp khóa (Unique Key Violation).
// - Giải pháp: Thực hiện câu lệnh SQL nguyên tử (Atomic Query) tại tầng DB bằng `ON CONFLICT DO UPDATE`.
//   Lúc này, nhân cơ sở dữ liệu sẽ đảm bảo việc khóa hàng (Row Lock) và xử lý tranh chấp hoàn hảo tại tầng nhân.
func UpsertPost(db *gorm.DB, title, content string, authorID uint) error {
	// TODO: Sử dụng raw SQL với mệnh đề ON CONFLICT DO UPDATE để thực hiện UPSERT
	return errors.New("not implemented")
}

// GetPostsWithPagination thực hiện lấy danh sách phân trang và tải trước quan hệ (Eager Loading).
//
// 🧠 TRÁNH LỖI N+1 QUERIES (Eager Loading Preload under the hood):
// - Nếu ta duyệt qua N posts và với mỗi post ta lại chạy `SELECT * FROM users WHERE id = post.author_id`,
//   ta sẽ tạo ra N+1 truy vấn tới Database, hủy diệt hoàn toàn hiệu năng hệ thống.
// - `db.Preload("Author")` của GORM giải quyết việc này thông qua **Eager Loading**:
//   1. GORM chạy 1 truy vấn lấy toàn bộ Posts: `SELECT * FROM posts LIMIT ... OFFSET ...`
//   2. GORM thu thập tất cả các `author_id` khác nhau từ mảng Posts vừa lấy.
//   3. GORM chỉ chạy DUY NHẤT 1 truy vấn gộp: `SELECT * FROM users WHERE id IN (author_id_1, author_id_2, ...)`
//   4. GORM tự động ráp (map) các đối tượng User vào thuộc tính `Author` tương ứng của từng struct Post trong bộ nhớ.
//   Tổng cộng ta luôn luôn chỉ tốn đúng 2 câu truy vấn SQL, hiệu năng cực kỳ tối ưu!
func GetPostsWithPagination(db *gorm.DB, page, limit int) ([]Post, error) {
	// TODO: Truy vấn danh sách Posts có phân trang (Limit, Offset) và Preload("Author")
	return nil, errors.New("not implemented")
}

func main() {
	// Demo GORM & SQL UPSERT tại đây
}
