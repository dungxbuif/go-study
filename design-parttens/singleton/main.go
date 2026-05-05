package db

import (
	"fmt"
	"sync"
	"time"
)

type Database struct {
    connection string
}

var (
    instance *Database
    once     sync.Once
)

func GetInstance() *Database {
    once.Do(func() {
        instance = &Database{
            connection: fmt.Sprintf("Connected at %s", time.Now()),
        }
    })
    return instance
}

func (d *Database) Query(sql string) {
    fmt.Printf("[DB] Running: %s\n", sql)
}
