package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateUserDTO đại diện cho Data Transfer Object để nhận thông tin đăng ký tài khoản.
//
// 🧠 CƠ CHẾ STRUCT TAGS & REFLECTION (Phản chiếu dưới nắp capo):
// - Go là một ngôn ngữ biên dịch tĩnh (statically typed). Sau khi biên dịch, tên các trường (struct fields) sẽ bị lược bỏ.
// - Để giữ lại siêu dữ liệu cho thời gian thực thi (Runtime), Go sử dụng cơ chế `Struct Tags` đặt trong dấu backticks ``.
// - Thư viện `go-playground/validator` (sử dụng mặc định trong Gin) sẽ dùng gói `reflect` (Reflection) của Go
//   để "soi" cấu trúc struct này tại thời điểm runtime, đọc các quy tắc xác thực như `required,min=3,max=20`.
// - ⚠️ Cảnh báo hiệu năng: Reflection trong Go khá tốn kém tài nguyên vì nó bypass hệ thống compiler tối ưu của Go.
//   Tuy nhiên, validator/v10 giải quyết việc này cực tốt bằng cách CACHE lại cấu trúc định nghĩa của Struct sau lần đầu tiên phản chiếu,
//   giúp các lần validate tiếp theo diễn ra cực nhanh!
type CreateUserDTO struct {
	Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Age      int    `json:"age" binding:"required,gte=18,lte=120"`
	Phone    string `json:"phone" binding:"required,vietnamphone"` // Sử dụng custom tag tự định nghĩa bên dưới
}

// ValidateVietnamPhone là hàm callback custom validator.
//
// 💡 CƠ CHẾ HOẠT ĐỘNG:
// - Hàm nhận vào một đối tượng interface `validator.FieldLevel`, cung cấp thông tin về trường đang được xác thực.
// - Dưới nắp capo, ta dùng `fl.Field().String()` để truy xuất giá trị chuỗi thực tế của trường và kiểm tra định dạng.
func ValidateVietnamPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// TODO: Kiểm tra số điện thoại bắt đầu bằng "0" hoặc "+84".
	// Trả về true nếu hợp lệ, ngược lại false.
	_ = phone
	return false
}

// RegisterHandler xử lý yêu cầu đăng ký người dùng mới.
//
// 🧠 AN TOÀN BẢO MẬT (Mass Assignment Protection):
// - Trong Node.js/Express, nếu ta destruct thẳng `const user = req.body`, hacker có thể gửi kèm các trường nhạy cảm như `{ "role": "admin" }`
//   và nếu ta lưu trực tiếp vào Database mà không filter, ứng dụng sẽ bị tấn công "Mass Assignment" (gán hàng loạt). Zod giải quyết điều này qua việc chỉ trả về dữ liệu đã validated.
// - Trong Go/Gin, việc ánh xạ dữ liệu trực tiếp vào một DTO cụ thể (`CreateUserDTO`) hoạt động như một lớp màng lọc bảo mật cứng rắn.
//   Mọi trường dư thừa ngoài định nghĩa của struct (ví dụ: `role`, `is_admin`) sẽ hoàn toàn bị bỏ qua khi parser phân tích JSON,
//   ngăn chặn triệt để lỗ hổng Mass Assignment ngay từ tầng đón nhận request!
func RegisterHandler(c *gin.Context) {
	// TODO:
	// 1. Bind JSON vào struct CreateUserDTO sử dụng c.ShouldBindJSON
	// 2. Trả về 400 nếu validation lỗi
	// 3. Trả về 201 với JSON { "success": true, "data": <user_info_không_chứa_password> }
}

func SetupValidationRouter() *gin.Engine {
	r := gin.Default()

	// Đăng ký custom validator tag "vietnamphone"
	//
	// 💡 ĐĂNG KÝ VỚI ENGINE CỦA GIN:
	// - Gin sử dụng một singleton instance của validator.Validate làm engine mặc định.
	// - Để đăng ký một tag mới, ta phải lấy instance đó thông qua `binding.Validator.Engine()`
	//   và ép kiểu nó sang `*validator.Validate` để gọi hàm `.RegisterValidation("tagname", callback)`.
	if v, ok := validator.New().(interface{}); ok {
		_ = v
		// Gợi ý: binding.Validator.Engine().(*validator.Validate).RegisterValidation("vietnamphone", ValidateVietnamPhone)
	}

	r.POST("/register", RegisterHandler)

	return r
}

func main() {
	r := SetupValidationRouter()
	r.Run(":8080")
}
