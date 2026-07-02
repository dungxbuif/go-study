package main

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// User biểu diễn cấu trúc dữ liệu người dùng.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// HashPassword thực hiện băm mật khẩu bằng thuật toán Bcrypt.
//
// 🧠 BĂM MẬT KHẨU (Bcrypt under the hood):
// - Bcrypt là một thuật toán băm (hashing) có chủ ý chạy CHẬM (adaptive hash algorithm), phụ thuộc vào "Work Factor" (Cost).
// - Khác với MD5 hay SHA256 chạy cực nhanh và dễ bị tấn công brute-force bằng GPU, Bcrypt được thiết kế để tiêu tốn nhiều thời gian CPU.
// - Dưới nắp capo, Bcrypt tự động tạo Salt (muối) ngẫu nhiên cho mỗi mật khẩu, trộn salt vào mật khẩu rồi thực hiện hàng ngàn chu kỳ tính toán
//   mã hóa Blowfish. Do đó, cùng một mật khẩu "123456" băm 2 lần sẽ ra 2 kết quả hoàn toàn khác nhau.
// - Dùng `bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)` trong package `golang.org/x/crypto/bcrypt`.
func HashPassword(password string) (string, error) {
	// TODO: Dùng bcrypt.GenerateFromPassword để băm mật khẩu
	return "", errors.New("not implemented")
}

// CheckPasswordHash đối chiếu mật khẩu thô và chuỗi hash băm lưu trữ.
// Dưới nắp capo, bcrypt trích xuất Salt từ chuỗi hash, băm mật khẩu thô với salt đó, rồi thực hiện so sánh
// thời gian không đổi (Constant-time comparison) để ngăn chặn các cuộc tấn công Timing Attack.
func CheckPasswordHash(password, hash string) bool {
	// TODO: Dùng bcrypt.CompareHashAndPassword để đối chiếu mật khẩu
	return false
}

// GenerateToken tạo ra JWT (JSON Web Token) định danh người dùng.
//
// 🧠 CƠ CHẾ KÝ JWT (JWT Signing & Claims):
// - JWT gồm 3 phần: Header, Payload (Claims), và Signature.
// - Ta sử dụng thư viện `github.com/golang-jwt/jwt/v5` để định nghĩa Claims struct chứa các thông tin như `username`, `role` và `exp` (expiration).
// - Signature được sinh ra bằng cách ký (bằng khóa bí mật HMAC-SHA256) lên chuỗi kết hợp của Header và Payload được mã hóa Base64Url.
// - Kích thước của token càng lớn sẽ càng tốn băng thông, do đó chỉ nên lưu trữ các thông tin định danh tối thiểu (như UserID, Roles) trong JWT Claims.
func GenerateToken(username, role string) (string, error) {
	// TODO: Tạo JWT token có hạn 1h chứa claims username và role
	return "", errors.New("not implemented")
}

// JWTMiddleware là bộ lọc xác thực JWT Token được truyền qua Authorization Header.
//
// 🧠 CƠ CHẾ KHÓA DỮ LIỆU CONTEXT (Gin Context Keys & Type Assertion):
// - Khi xác thực token thành công, ta lưu trữ thông tin Claims vào `gin.Context` thông qua `c.Set("claims", claims)`.
// - Dưới nắp capo, `gin.Context` lưu trữ các giá trị này trong một `map[string]any` (Keys map).
// - Các middleware hoặc handler phía sau có thể lấy dữ liệu ra bằng `c.Get("claims")`.
// - ⚠️ LƯU Ý: Vì hàm `c.Get()` trả về kiểu dữ liệu `any` (interface{}), ta BẮT BUỘC phải thực hiện ép kiểu
//   (Type Assertion, ví dụ: `claims, ok := val.(*CustomClaims)`) để có thể sử dụng lại các thuộc tính tĩnh một cách an toàn!
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Lấy header "Authorization", parse Bearer token.
		// Xác thực token và c.Set("claims", claims) trước khi gọi c.Next()
		c.Next()
	}
}

// RoleRequired kiểm tra xem người dùng có quyền truy cập phù hợp dựa trên vai trò (Role).
func RoleRequired(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Lấy claims từ context, nếu role không khớp -> c.AbortWithStatusJSON(403) và return!
		c.Next()
	}
}

// 🧠 CƠ CHẾ ĐỌC BODY NHIỀU LẦN (ZK Proof/Mock API Gotcha):
// - Giả sử ta muốn xây dựng một middleware để mock ZK (Zero-Knowledge) Verification qua Request Body.
// - Trong Go, `c.Request.Body` thuộc kiểu `io.ReadCloser` (vốn là một Stream mạng).
// - ⚠️ CẢNH BÁO: Vì là Stream, nó chỉ cho phép ĐỌC MỘT LẦN DUY NHẤT. Khi một middleware đọc hết Body để kiểm tra,
//   con trỏ stream sẽ nhảy đến EOF (End of File). Handler tiếp theo gọi `ShouldBindJSON` sẽ bị lỗi trống dữ liệu (EOF error).
// - 🛠️ GIẢI PHÁP: Nếu bắt buộc phải đọc trước Request Body:
//   1. Đọc toàn bộ body thành mảng bytes bằng `io.ReadAll(c.Request.Body)`.
//   2. Tiến hành xử lý/xác thực bytes đó.
//   3. "Khôi phục" stream bằng cách gán một ReadCloser mới: `c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))`
//      trước khi chuyển tiếp qua `c.Next()`.
func SetupAuthRouter() *gin.Engine {
	r := gin.Default()

	// TODO:
	// POST /auth/register -> Lưu user đã băm password vào map
	// POST /auth/login -> Xác thực password và trả về JWT token
	// GET /auth/profile -> Yêu cầu JWTMiddleware, trả về thông tin user trong token

	return r
}

func main() {
	r := SetupAuthRouter()
	r.Run(":8080")
}
