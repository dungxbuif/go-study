package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator định nghĩa bộ validate dữ liệu đầu vào sử dụng go-playground/validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate thực hiện kiểm tra tính hợp lệ của dữ liệu
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// UserDTO định nghĩa cấu trúc dữ liệu JSON nhận được từ Client
type UserDTO struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

func main() {
	e := echo.New()

	// 1. Cấu hình Custom Validator cho Echo
	e.Validator = &CustomValidator{validator: validator.New()}

	// 2. Custom Middleware Ghi Log (Onion Model)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			// Gọi next(c) để chuyển xử lý xuống handler tiếp theo
			err := next(c)
			
			// Mã chạy sau khi handler đã hoàn tất xử lý (Onion model)
			stop := time.Now()
			
			fmt.Printf("[%s] %d %s %s %s\n",
				c.Request().Method,
				c.Response().Status,
				c.Path(),
				stop.Sub(start),
				c.Request().RemoteAddr,
			)
			return err
		}
	})

	// 3. Centralized Error Handler (Xử lý lỗi tập trung)
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := err.Error()

		// Kiểm tra nếu lỗi thuộc kiểu HTTPError của Echo
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = fmt.Sprintf("%v", he.Message)
		}

		_ = c.JSON(code, map[string]interface{}{
			"success": false,
			"error":   message,
		})
	}

	// 4. API Endpoints
	e.POST("/users", func(c echo.Context) error {
		u := new(UserDTO)
		
		// Trích xuất dữ liệu JSON tự động
		if err := c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Dữ liệu JSON không hợp lệ")
		}
		
		// Thực hiện Validate dữ liệu thông qua Custom Validator đã đăng ký
		if err := c.Validate(u); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		// Trả về dữ liệu JSON kèm trạng thái 201 Created
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"success": true,
			"data":    u,
		})
	})

	fmt.Println("Echo server đang chạy trên cổng 8080...")
	e.Logger.Fatal(e.Start(":8080"))
}
