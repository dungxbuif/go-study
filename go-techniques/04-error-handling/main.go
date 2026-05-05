package main

import (
	"errors"
	"fmt"
)

// Go treats errors as values — no exceptions
// Key ideas:
//   1. Sentinel errors  — compare with ==
//   2. Custom error types — carry extra context
//   3. Error wrapping (fmt.Errorf + %w) — preserve chain
//   4. errors.Is / errors.As — unwrap and inspect

// --- 1. Sentinel errors ---

var (
	ErrNotFound   = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

// --- 2. Custom error type ---

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %q: %s", e.Field, e.Message)
}

// --- 3. Wrapping errors ---

type DB struct {
	data map[int]string
}

func (db *DB) FindUser(id int) (string, error) {
	user, ok := db.data[id]
	if !ok {
		return "", fmt.Errorf("FindUser(%d): %w", id, ErrNotFound)
	}
	return user, nil
}

func (db *DB) GetProfile(id int) (string, error) {
	user, err := db.FindUser(id)
	if err != nil {
		return "", fmt.Errorf("GetProfile: %w", err) // wrap again
	}
	return "Profile of " + user, nil
}

func validate(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Message: "must be non-negative"}
	}
	return nil
}

func main() {
	db := &DB{data: map[int]string{1: "Alice", 2: "Bob"}}

	fmt.Println("=== 1. Sentinel Error ===")
	_, err := db.GetProfile(99)
	fmt.Println(err)                          // full chain
	fmt.Println(errors.Is(err, ErrNotFound))  // true — unwraps through chain

	fmt.Println("\n=== 2. Custom Error Type ===")
	err = validate(-1)
	fmt.Println(err)

	var ve *ValidationError
	if errors.As(err, &ve) { // unwrap to concrete type
		fmt.Printf("field=%q msg=%q\n", ve.Field, ve.Message)
	}

	fmt.Println("\n=== 3. Multiple Returns — happy path ===")
	profile, err := db.GetProfile(1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(profile)
}
