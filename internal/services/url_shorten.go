package services

import (
	"crypto/sha256"
	"encoding/base64"
	"net/url"
)

type ServiceInterface interface {
	GenerateShortURL(longURL string) string
}

// Service is a concrete implementation of ServiceInterface
type Service struct{}

func (s *Service) GenerateShortURL(longURL string) string {
	return generateShortURL(longURL)
}

// GenerageShortURL generates a short URL based on the input long URL
func generateShortURL(longURL string) string {
	// Parse and validate the URL
	parsedURL, err := url.Parse(longURL)
	if err != nil || parsedURL.Host == "" {
		return "" // Return empty string for invalid URLs
	}

	// Extract the hostname and path
	host := parsedURL.Host
	path := parsedURL.Path

	// Generate a bash-based short URL
	shortURL := hashURL(host + path)
	return shortURL
}

// hashURL generates a short string based on a hash of the input
func hashURL(input string) string {
	hash := sha256.Sum256([]byte(input))
	encoded := base64.URLEncoding.EncodeToString(hash[:6]) // Use first 6 bytes of the string
	return encoded
}
