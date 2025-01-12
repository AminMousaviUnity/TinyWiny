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

// Centralized setup for Handlers
func setupHandlers(t *testing.T) (*Handlers, redismock.ClientMock, string) {
	// Mock Redis client and mock controller
	mockRedis, mock := redismock.NewClientMock()

	// Initialize Redis storage with the mock Redis client
	mockStorage := storage.NewRedisStorage(mockRedis)

	// Use the real service implementation
	mockService := &services.Service{}
	baseURL := "http://localhost:8888"

	// Initialize handlers with mock dependencies
	h := NewHandlers(baseURL, mockStorage, mockService)
	return h, mock, baseURL
}

func TestShortenURLHandler(t *testing.T) {
	h, mock, baseURL := setupHandlers(t)

	// Table-driven tests for multiple cases
	tests := []struct {
		name           string
		input          string
		expectedCode   int
		mockBehavior   func()
		expectedOutput string
	}{
		{
			name:         "Valid URL",
			input:        `{"long_url":"http://example.com"}`,
			expectedCode: http.StatusCreated,
			mockBehavior: func() {
				mock.ExpectSet("o3mm9u6v", "http://example.com", 24*time.Hour).SetVal("OK")
			},
			expectedOutput: baseURL + "/o3mm9u6v",
		},
		{
			name:         "Missing long_url field",
			input:        `{}`,
			expectedCode: http.StatusBadRequest,
			mockBehavior: func() {}, // No mock needed for invalid input
			expectedOutput: "",
		},
		{
			name:         "Redis Save Error",
			input:        `{"long_url":"http://example.com"}`,
			expectedCode: http.StatusInternalServerError,
			mockBehavior: func() {
				mock.ExpectSet("o3mm9u6v", "http://example.com", 24*time.Hour).SetErr(assert.AnError)
			},
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock behavior
			tt.mockBehavior()

			// Prepare the request
			req, err := http.NewRequest("POST", "/shorten", bytes.NewBufferString(tt.input))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Record the response
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.ShortenURLHandler)
			handler.ServeHTTP(rr, req)

			// Validate the response code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// If output is expected, validate the response body
			if tt.expectedOutput != "" {
				var resp struct {
					ShortURL string `json:"short_url"`
				}
				err = json.NewDecoder(rr.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, resp.ShortURL)
			}

			// Verify mock expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	h, mock, _ := setupHandlers(t)

	// Table-driven tests for multiple cases
	tests := []struct {
		name         string
		shortURL     string
		expectedCode int
		mockBehavior func()
		expectedLoc  string
	}{
		{
			name:         "Valid Short URL",
			shortURL:     "o3mm9u6v",
			expectedCode: http.StatusFound,
			mockBehavior: func() {
				mock.ExpectGet("o3mm9u6v").SetVal("http://example.com")
			},
			expectedLoc: "http://example.com",
		},
		{
			name:         "Short URL Not Found",
			shortURL:     "unknown",
			expectedCode: http.StatusNotFound,
			mockBehavior: func() {
				mock.ExpectGet("unknown").RedisNil()
			},
			expectedLoc: "",
		},
		{
			name:         "Redis Error",
			shortURL:     "error",
			expectedCode: http.StatusInternalServerError,
			mockBehavior: func() {
				mock.ExpectGet("error").SetErr(assert.AnError)
			},
			expectedLoc: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock behavior
			tt.mockBehavior()

			// Prepare the request
			req, err := http.NewRequest("GET", "/"+tt.shortURL, nil)
			assert.NoError(t, err)

			// Record the response
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.RedirectHandler)
			handler.ServeHTTP(rr, req)

			// Validate the response code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// If a redirect is expected, validate the Location header
			if tt.expectedLoc != "" {
				location := rr.Header().Get("Location")
				assert.Equal(t, tt.expectedLoc, location)
			}

			// Verify mock expectations
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
