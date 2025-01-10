package services

import (
	"strings"
	"sync"
)

var mu sync.RWMutex

// GenerateShortURL generates a new short URL based on the counter
func GenerateShortURL(LongURL string) string {
	mu.Lock()
	defer mu.Unlock()
	shortURL := generateShortURL(LongURL)
	return shortURL
}

func generateShortURL(LongURL string) string {
	// Extract the part of the address without the TLD
	addrWithoutTLD, TLD := breakDownAddr(LongURL)

	vowels := "aeiouAEIOU"
	var result strings.Builder

	// Remove vowels from the address without the TLD
	for _, char := range addrWithoutTLD {
		if !strings.ContainsRune(vowels, char) {
			result.WriteRune(char)
		}
	}

	result.WriteString(TLD)
	return result.String()
}

// breakDownAddr breaks down the URL before and after the Top-Level Domain
func breakDownAddr(fullAddr string) (string, string) {
	// Position of last dot
	lastDotIndex := strings.LastIndex(fullAddr, ".")

	if lastDotIndex == -1 {
		return fullAddr, ""
	}

	return fullAddr[:lastDotIndex], fullAddr[lastDotIndex:]
}
