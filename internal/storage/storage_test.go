package storage

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
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
	err := storage.SaveURLWithExpiry(ctx, short, long, time.Hour*24)
	assert.NoError(t, err, "Expected no error while saving URL")

	// Expect the GetOriginalURL operation
	mock.ExpectGet(short).SetVal(long)

	// Test GetOriginalURL
	result, err := storage.GetOriginalURL(ctx, short)
	assert.NoError(t, err, "Expected no error for existing URL")
	assert.Equal(t, long, result, "Expected long URL to match")

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	// Create a mock Redis client
	mockRedis, mock := redismock.NewClientMock()
	storage := NewRedisStorage(mockRedis)

	short := "http://unknown.com" // A short URL that doesn't exist

	// Expect the GetOriginalURL operation to return a Redis nil error
	mock.ExpectGet(short).RedisNil()

	// Test GetOriginalURL
	ctx := context.Background()
	result, err := storage.GetOriginalURL(ctx, short)
	assert.Equal(t, redis.Nil, err, "Expected redis.Nil for unknown short URL")
	assert.Equal(t, "", result, "Expected result to be empty for unknown short URL")

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}
