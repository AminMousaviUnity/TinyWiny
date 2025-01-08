package storage

import (
	"fmt"
	"sync"
)

// In-memory storage
var (
	mu      sync.RWMutex
	BaseURL = "http://localhost:8888/"
)

// GenerateShortURL generates a short URL based on a counter
var counter int64 = 1

// GenerateShortURL generates a new short URL based on the counter
func GenerateShortURL() string {
	mu.Lock()
	defer mu.Unlock()
	short := fmt.Sprintf("%d", counter)
	counter++
	return short
}
