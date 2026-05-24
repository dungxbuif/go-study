package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateUserDTO struct {
	Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Age      int    `json:"age" binding:"required,gte=18,lte=120"`
	Phone    string `json:"phone" binding:"required,vietnamphone"`
}

func ValidateVietnamPhone(fl validator.FieldLevel) bool {
	// TODO: Kiểm tra số điện thoại bắt đầu bằng "0" hoặc "+84".
	// Trả về true nếu hợp lệ, ngược lại false.
	phone := fl.Field().String()
	_ = phone
	return false
}

func RegisterHandler(c *gin.Context) {
	// TODO:
	// 1. Bind JSON vào struct CreateUserDTO sử dụng c.ShouldBindJSON
	// 2. Trả về 400 nếu validation lỗi
	// 3. Trả về 201 với JSON { "success": true, "data": <user_info_không_chứa_password> }
}

func SetupValidationRouter() *gin.Engine {
	r := gin.Default()

	// Đăng ký custom validator tag "vietnamphone"
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
