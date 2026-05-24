package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func HashPassword(password string) (string, error) {
	// TODO: Dùng bcrypt.GenerateFromPassword để băm mật khẩu
	return "", errors.New("not implemented")
}

func CheckPasswordHash(password, hash string) bool {
	// TODO: Dùng bcrypt.CompareHashAndPassword để đối chiếu mật khẩu
	return false
}

func GenerateToken(username, role string) (string, error) {
	// TODO: Tạo JWT token có hạn 1h chứa claims username và role
	return "", errors.New("not implemented")
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Lấy header "Authorization", parse Bearer token.
		// Xác thực token và c.Set("claims", claims) trước khi gọi c.Next()
		c.Next()
	}
}

func RoleRequired(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Lấy claims từ context, nếu role không khớp -> c.AbortWithStatusJSON(403)
		c.Next()
	}
}

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
