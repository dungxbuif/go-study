package main

import "fmt"

// Embedding = composition over inheritance
// Go does not have class inheritance — use embedding to reuse and extend behavior

type Animal struct {
	Name string
}

func (a Animal) Speak() string { return a.Name + " makes a sound" }
func (a Animal) Breathe()      { fmt.Println(a.Name, "is breathing") }

type Dog struct {
	Animal        // embedded — all Animal methods and fields are promoted to Dog
	Breed  string
}

// Override Animal.Speak
func (d Dog) Speak() string { return d.Name + " says: Woof!" }

// --- Promoted fields ---

type Address struct {
	Street string
	City   string
}

type Person struct {
	Address        // Street and City are promoted directly onto Person
	Name    string
}

// --- Embedding pointer to share mutable state ---

type Logger struct {
	prefix string
}

func (l *Logger) Log(msg string) {
	fmt.Printf("[%s] %s\n", l.prefix, msg)
}

type UserHandler struct {
	*Logger
	count int
}

func (u *UserHandler) Handle(path string) {
	u.count++
	u.Log(fmt.Sprintf("handling %s (total: %d)", path, u.count))
}

func main() {
	fmt.Println("=== 1. Method Promotion ===")
	dog := Dog{Animal: Animal{Name: "Rex"}, Breed: "Labrador"}
	fmt.Println(dog.Speak())        // Dog's override
	dog.Breathe()                   // promoted from Animal
	fmt.Println(dog.Animal.Speak()) // call Animal's original directly

	fmt.Println("\n=== 2. Field Promotion ===")
	p := Person{
		Name:    "Alice",
		Address: Address{Street: "123 Main St", City: "Hanoi"},
	}
	fmt.Println(p.City)           // promoted — access directly
	fmt.Println(p.Address.Street) // also valid full path

	fmt.Println("\n=== 3. Embedding Pointer ===")
	h := &UserHandler{Logger: &Logger{prefix: "USER"}}
	h.Handle("/api/users")
	h.Handle("/api/users/1")
}
