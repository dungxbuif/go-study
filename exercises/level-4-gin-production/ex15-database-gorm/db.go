package main

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:50"`
	Posts    []Post `gorm:"foreignKey:AuthorID"`
}

type Post struct {
	gorm.Model
	Title    string `gorm:"size:200"`
	Content  string `gorm:"type:text"`
	AuthorID uint
	Author   User
}

func SetupDB() (*gorm.DB, error) {
	// TODO: Khởi tạo kết nối SQLite in-memory bằng GORM
	// Tự động migrate: db.AutoMigrate(&User{}, &Post{})
	return nil, errors.New("not implemented")
}

func UpsertPost(db *gorm.DB, title, content string, authorID uint) error {
	// TODO: Sử dụng raw SQL với mệnh đề ON CONFLICT DO UPDATE để thực hiện UPSERT
	return errors.New("not implemented")
}

func GetPostsWithPagination(db *gorm.DB, page, limit int) ([]Post, error) {
	// TODO: Truy vấn danh sách Posts có phân trang (Limit, Offset) và Preload("Author")
	return nil, errors.New("not implemented")
}

func main() {
	// Demo GORM & SQL UPSERT tại đây
}
