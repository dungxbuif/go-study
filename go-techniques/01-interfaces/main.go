package main

import (
	"fmt"
	"math"
)

// =============================================
// TRIẾT LÝ: "Accept interfaces, return structs"
// Interface trong Go là implicit — không cần khai báo "implements"
// Struct tự động thỏa interface nếu có đủ method
// =============================================

// --- 1. Small interface ---
// Go khuyến khích interface nhỏ, 1-2 methods
// Thay vì 1 interface lớn, tách thành nhiều interface nhỏ

// 1. Small interface — Shape chỉ có 2 methods, Circle và Rectangle implement ngầm định
// 2. Dependency injection — UserService chỉ biết Storage interface, swap MemoryStorage → RedisStorage không cần đổi business logic
// 3. Type switch — xử lý any type an toàn

type Stringer interface {
	String() string
}

type Saver interface {
	Save() error
}

// Kết hợp interface (composition)
type StringSaver interface {
	Stringer
	Saver
}

// --- 2. Structs implement interface ngầm định ---

type Circle struct {
	Radius float64
}

type Rectangle struct {
	Width, Height float64
}

// Cả hai implement Shape — không cần viết "implements Shape"
type Shape interface {
	Area() float64
	Perimeter() float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// "Accept interfaces" — hàm nhận interface, không nhận concrete type
func printShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f | Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// --- 3. Interface như một contract cho dependency injection ---
// Dễ swap implementation mà không đổi business logic

type Storage interface {
	Save(key, value string) error
	Get(key string) (string, error)
}

type MemoryStorage struct {
	data map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{data: make(map[string]string)}
}

func (m *MemoryStorage) Save(key, value string) error {
	m.data[key] = value
	return nil
}

func (m *MemoryStorage) Get(key string) (string, error) {
	v, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("key %q not found", key)
	}
	return v, nil
}

// UserService chỉ biết đến Storage interface — không biết MemoryStorage hay Redis hay Postgres
type UserService struct {
	storage Storage // "Accept interfaces"
}

func NewUserService(s Storage) *UserService { // "Return structs"
	return &UserService{storage: s}
}

func (u *UserService) CreateUser(name, email string) error {
	return u.storage.Save(name, email)
}

func (u *UserService) GetEmail(name string) (string, error) {
	return u.storage.Get(name)
}

// --- 4. Interface{} / any — type assertion và type switch ---

func describe(i any) {
	switch v := i.(type) {
	case int:
		fmt.Printf("int: %d\n", v)
	case string:
		fmt.Printf("string: %q\n", v)
	case Shape:
		fmt.Printf("Shape with area: %.2f\n", v.Area())
	default:
		fmt.Printf("unknown type: %T\n", v)
	}
}

func main() {
	fmt.Println("=== 1. Small Interface ===")
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 6},
	}
	for _, s := range shapes {
		printShapeInfo(s)
	}

	fmt.Println("\n=== 2. Dependency Injection via Interface ===")
	store := NewMemoryStorage()
	svc := NewUserService(store) // swap store thành RedisStorage mà không đổi UserService

	_ = svc.CreateUser("alice", "alice@example.com")
	email, err := svc.GetEmail("alice")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Found email:", email)
	}

	_, err = svc.GetEmail("bob")
	fmt.Println("Get bob:", err)

	fmt.Println("\n=== 3. Type Switch ===")
	describe(42)
	describe("hello")
	describe(Circle{Radius: 3})
	describe(true)
}
