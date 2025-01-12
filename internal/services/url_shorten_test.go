package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortURL_ValidURLs(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic URL",
			input:    "http://example.com/path",
			expected: hashURL("example.com/path"),
		},
		{
			name:     "URL with Query Params",
			input:    "https://example.com/path?query=123",
			expected: hashURL("example.com/path"), // Query params don't affect the short URL
		},
		{
			name:     "URL with Trailing Slash",
			input:    "http://example.com/",
			expected: hashURL("example.com/"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GenerateShortURL(tt.input)
			assert.Equal(t, tt.expected, result, "Short URL did not match the expected value")
		})
	}
}

func TestGenerateShortURL_InvalidURLs(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name  string
		input string
	}{
		{name: "Empty URL", input: ""},
		{name: "Malformated URL", input: "http:/invalid-url"},
		{name: "No Host", input: "http:///"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GenerateShortURL(tt.input)
			assert.Equal(t, "", result, "Expected an ampty string for invalid URL")
		})
	}
}

func TestHashURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // Expected length of the hash
	}{
		{
			name:     "Basic Input",
			input:    "example.com/path",
			expected: 8, // Base64 encoding of first 6 bytes results in an 8-character string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hashURL(tt.input)
			assert.Equal(t, tt.expected, len(result), "Hash length did not match expected length")
		})
	}
}
