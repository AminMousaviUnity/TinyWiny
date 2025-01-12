package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aminmousaviunity/TinyWiny/internal/services"
	"github.com/aminmousaviunity/TinyWiny/internal/storage"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLHandler(t *testing.T) {
	// Create a mock Redis client
	mockRedis, mock := redismock.NewClientMock()
	mockStorage := storage.NewRedisStorage(mockRedis) // Use RedisStorage instance
	mockService := &services.Service{}               // Use the actual Service implementation
	baseURL := "http://localhost:8888"

	// Initialize Handlers with mock dependencies
	h := NewHandlers(baseURL, mockStorage, mockService)

	// Mock Redis behavior for SaveURLWithExpiry
	short := "o3mm9u6v"
	long := "http://example.com"
	mock.ExpectSet(short, long, 24*time.Hour).SetVal("OK")

	// Prepare the request payload
	reqBody, err := json.Marshal(struct {
		LongURL string `json:"long_url"`
	}{LongURL: long})
	if err != nil {
		t.Fatalf("Could not marshal request body: %v", err)
	}

	// Create an HTTP POST request
	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ShortenURLHandler)
	handler.ServeHTTP(rr, req)

	// Check response code
	assert.Equal(t, http.StatusCreated, rr.Code, "Handler returned wrong status code")

	// Parse the response body
	var resp struct {
		ShortURL string `json:"short_url"`
	}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	// Check the short URL in the response
	expectedShortURL := baseURL + "/" + short
	assert.Equal(t, expectedShortURL, resp.ShortURL, "Response contains incorrect short URL")

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedirectHandler(t *testing.T) {
	// Create a mock Redis client
	mockRedis, mock := redismock.NewClientMock()
	mockStorage := storage.NewRedisStorage(mockRedis) // Use RedisStorage instance
	mockService := &services.Service{}               // Use the actual Service implementation
	baseURL := "http://localhost:8888"

	// Initialize Handlers with mock dependencies
	h := NewHandlers(baseURL, mockStorage, mockService)

	// Mock Redis behavior for GetOriginalURL
	short := "o3mm9u6v"
	long := "http://example.com"
	mock.ExpectGet(short).SetVal(long)

	// Create a request to the short URL
	req, err := http.NewRequest("GET", "/"+short, nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.RedirectHandler)
	handler.ServeHTTP(rr, req)

	// Check response code
	assert.Equal(t, http.StatusFound, rr.Code, "Handler returned wrong status code")

	// Check the redirect location
	location := rr.Header().Get("Location")
	assert.Equal(t, long, location, "Expected redirect to correct location")

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}
