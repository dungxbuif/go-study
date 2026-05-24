package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	id        int
	name      string
	email     string
	createdAt time.Time
}

func NewUser(id int, name string, email string) *User {
	return &User{
		id:        id,
		name:      name,
		email:     email,
		createdAt: time.Now(),
	}
}

type UserService struct {
	users  []User
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users:  []User{},
		nextID: 1,
	}
}

func (s *UserService) Create(name, email string) User {
	u := User{
		id:        s.nextID,
		name:      name,
		email:     email,
		createdAt: time.Now(),
	}
	s.users = append(s.users, u)
	s.nextID++
	return u
}

func (s *UserService) FindByID(id int) (User, error) {
	for _, u := range s.users {
		if u.id == id {
			return u, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (s *UserService) FindByEmail(email string) (User, error) {
	for _, u := range s.users {
		if u.email == email {
			return u, nil
		}
	}
	return User{}, ErrUserNotFound
}

func (s UserService) Count() int {
	return len(s.users)
}

func main() {
	svc := NewUserService()

	svc.Create("Alice", "alice@example.com")
	svc.Create("Bob", "bob@example.com")
	svc.Create("Charlie", "charlie@example.com")

	fmt.Printf("Total users: %d\n", svc.Count())

	if u, err := svc.FindByID(2); err == nil {
		fmt.Printf("Found by ID: %+v\n", u)
	}

	if u, err := svc.FindByEmail("alice@example.com"); err == nil {
		fmt.Printf("Found by email: %+v\n", u)
	}

	if _, err := svc.FindByID(999); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}