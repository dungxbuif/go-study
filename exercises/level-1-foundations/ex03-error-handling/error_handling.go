package main

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s' — %s", e.Field, e.Message)
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id '%s' not found", e.Resource, e.ID)
}

type Config struct {
	Port int
}

func ParseConfig(filename string) (Config, error) {
	if filename == "nonexistent.json" {
		return Config{}, fmt.Errorf("failed to read file: %w", &NotFoundError{Resource: "config file", ID: filename})
	}
	if filename == "invalid.json" {
		return Config{}, fmt.Errorf("failed to parse content: %w", &ValidationError{Field: "content", Message: "invalid JSON format"})
	}
	if filename == "missing_port.json" {
		return Config{}, fmt.Errorf("failed to validate: %w", &ValidationError{Field: "port", Message: "is required"})
	}
	if filename == "app.json" {
		return Config{Port: 8080}, nil
	}
	return Config{}, errors.New("unknown file")
}

func main() {
	testFiles := []string{
		"nonexistent.json",
		"invalid.json",
		"missing_port.json",
		"app.json",
	}

	for _, file := range testFiles {
		config, err := ParseConfig(file)
		if err != nil {
			var nfErr *NotFoundError
			var valErr *ValidationError

			if errors.As(err, &nfErr) {
				fmt.Printf("[NOT FOUND] %s: %s\n", nfErr.Resource, nfErr.ID)
			} else if errors.As(err, &valErr) {
				fmt.Printf("[VALIDATION] field '%s': %s\n", valErr.Field, valErr.Message)
			} else {
				fmt.Printf("[UNKNOWN] %v\n", err)
			}
		} else {
			fmt.Printf("Config loaded successfully: Port = %d\n", config.Port)
		}
	}
}
