package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"TinyWiny/handlers"
	"TinyWiny/storage"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupMockRedis() (redismock.ClientMock, *redis.Client) {
	mockRedis, mock := redismock.NewClientMock()
	storage.InitRedis(mockRedis) // Initialize Redis in your storage layer
	return mock, mockRedis
}

func createRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	var reqBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&reqBody).Encode(body); err != nil {
			t.Fatalf("Failed to encode request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, &reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestShortenURLHandler(t *testing.T) {
	mock, _ := setupMockRedis()

	short := "http://xmpl.com"
	long := "http://example.com"
	mock.ExpectSet(short, long, 24*time.Hour).SetVal("OK")

	req := createRequest(t, "POST", "/shorten", handlers.ShortenURLRequest{LongURL: long})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ShortenURLHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Unexpected status code")

	var resp handlers.ShortenURLResponse
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&resp), "Failed to decode response")

	expectedShortURL := "http://localhost:8888/" + short
	assert.Equal(t, expectedShortURL, resp.ShortURL, "Incorrect short URL in response")

	assert.NoError(t, mock.ExpectationsWereMet(), "Redis expectations were not met")
}

func TestRedirectHandler(t *testing.T) {
	mock, _ := setupMockRedis()

	short := "http://xmpl.com"
	long := "http://example.com"
	mock.ExpectGet(short).SetVal(long)

	req := createRequest(t, "GET", "/"+short, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.RedirectHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusFound, rr.Code, "Unexpected status code")
	assert.Equal(t, long, rr.Header().Get("Location"), "Incorrect redirect location")

	assert.NoError(t, mock.ExpectationsWereMet(), "Redis expectations were not met")
}

func TestRedirectHandler_NotFound(t *testing.T) {
	mock, _ := setupMockRedis()

	short := "nonexistentCode"
	mock.ExpectGet(short).RedisNil()

	req := createRequest(t, "GET", "/"+short, nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.RedirectHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "Expected 404 for nonexistent code")
	assert.Contains(t, rr.Body.String(), "URL not found", "Expected error message in response")

	assert.NoError(t, mock.ExpectationsWereMet(), "Redis expectations were not met")
}
