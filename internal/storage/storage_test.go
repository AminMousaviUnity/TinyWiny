package storage

import (
	"testing"
	"time"
	"github.com/aminmousaviunity/TinyWiny/internal/storage"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestSaveAndGetURL(t *testing.T) {
	// Create a mock Redis client
	mockRedis, mock := redismock.NewClientMock()
	t.Logf("mockRedis type: %T", mockRedis)

	// Initialize the storage package with the mock Redis client
	storage.InitRedis(mockRedis)

	short := "http://xmlp.com"
	long := "http:example.com"

	// Expect the SaveURLWithExpiry operation
	mock.ExpectSet(short, long, time.Hour*24).SetVal("OK")

	// Test SaveURLWithExpiry
	err := storage.SaveURLWithExpiry(short, long, time.Hour*24)
	assert.NoError(t, err)

	// Expect the GetOriginalURL operation
	mock.ExpectGet(short).SetVal(long)

	// Test GetOriginalURL
	result, exists := storage.GetOriginalURL(short)

	assert.True(t, exists, "Expected URL to exist for short code")
	assert.Equal(t, long, result, "Expected long URL to match")

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	// Create a mock Redis client
	mockRedis, mock := redismock.NewClientMock()

	// Initialize the storage package with the mock Redis client
	storage.InitRedis(mockRedis)

	short := "http://unknown.com" // A short URL that doesn't exist

	// Expect the GetOriginalURL operation to return a Redis nil error
	mock.ExpectGet(short).RedisNil()

	// Test GetOriginalURL
	_, exists := storage.GetOriginalURL(short)
	assert.False(t, exists, "Expected URL to not exist for unknown short code")

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}
