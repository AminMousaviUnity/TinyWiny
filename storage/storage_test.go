package storage

import (
	"testing"
)

func TestSaveAndGetURL(t *testing.T) {
	InitStorage() // Reset storage before testing
	short := "123"
	long := "http://example.com"

	SaveURL(short, long)
	result, exists := GetOriginalURL(short)

	if !exists {
		t.Errorf("Expected URL to exist for short code %s", short)
	}

	if result != long {
		t.Errorf("Expected %s, got %s", long, result)
	}
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	InitStorage()
	_, exists := GetOriginalURL("nonexistent")
	if exists {
		t.Errorf("Expected URL to not exist for nonexistant short code")
	}
}
