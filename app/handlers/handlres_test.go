package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"TinyWiny/app/storage"
)

func TestShortenURLHandler(t *testing.T) {
	// Initialize storage
	storage.InitStorage()

	// Create a sample requset
	payload := `{"long_url":"http://example.com"}`
	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShortenURLHandler)
	handler.ServeHTTP(rr, req)

	// Check response code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check response body
	var resp ShortenURLResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("Could not decode response: %v", err)
	}

	if resp.ShortURL == "" {
		t.Errorf("Expected short URL, got empty string")
	}
}

func TestRedirectHandler(t *testing.T) {
	// Initialize storage
	storage.InitStorage()

	// Prepopulate storage
	short := "123"
	long := "http://example.com"
	storage.SaveURL(short, long)

	// Create a request to the short URL
	req, err := http.NewRequest("GET", "/"+short, nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RedirectHandler)
	handler.ServeHTTP(rr, req)

	// Check response code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Check the redirect location
	location := rr.Header().Get("Location")
	if location != long {
		t.Errorf("Expected redirect to %s, got %s", long, location)
	}
}