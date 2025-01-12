package storage

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestSaveAndGetURL(t *testing.T) {
	// Create a mock Redis client
	mockRedis, mock := redismock.NewClientMock()
	storage := NewRedisStorage(mockRedis)

	short := "1MnZAnMm"
	long := "http://google.com"

	// Expect the SaveURLWithExpiry operation
	mock.ExpectSet(short, long, time.Hour*24).SetVal("OK")

	// Test SaveURLWithExpiry
	ctx := context.Background()
	err := storage.SaveURLWithExpiry(ctx, short, long, time.Hour*24) // Call on the storage instance
	assert.NoError(t, err)

	// Expect the GetOriginalURL operation
	mock.ExpectGet(short).SetVal(long)

	// Test GetOriginalURL
	result, exists := storage.GetOriginalURL(ctx, short) // Call on the storage instance
	assert.True(t, exists, "Expected URL to exist for short code")
	assert.Equal(t, long, result, "Expected long URL to match")

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	// Create a mock Redis client
	mockRedis, mock := redismock.NewClientMock()
	storage := NewRedisStorage(mockRedis) // Correctly initialize the RedisStorage instance

	short := "http://unknown.com" // A short URL that doesn't exist

	// Expect the GetOriginalURL operation to return a Redis nil error
	mock.ExpectGet(short).RedisNil()

	// Test GetOriginalURL
	ctx := context.Background()
	_, exists := storage.GetOriginalURL(ctx, short) // Call on the storage instance
	assert.False(t, exists, "Expected URL to not exist for unknown short code")

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}
