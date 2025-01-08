package storage

import (
	"fmt"
	"sync"
)

// In-memory storage
var (
	urlStore = make(map[string]string) // shortURL -> longURL
	mu       sync.RWMutex
	BaseURL  = "http://localhost:8888/"
)

// GenerateShortURL generates a short URL based on a counter
var counter int64 = 1

func InitStorage() {
	counter = 1
}

// SaveURL stores the mapping of short URL to long URL
func SaveURL(shortURL, longURL string) {
	mu.Lock()
	defer mu.Unlock()
	urlStore[shortURL] = longURL
}

// GetOriginalURL retrieves the long URL for a given short URL
func GetOriginalURL(shortURL string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	longURL, exists := urlStore[shortURL]
	return longURL, exists
}

// GenerateShortURL generates a new short URL based on the counter
func GenerateShortURL() string {
	mu.Lock()
	defer mu.Unlock()
	short := fmt.Sprintf("%d", counter)
	counter++
	return short
}
