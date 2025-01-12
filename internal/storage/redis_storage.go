package storage

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// StorageInterface defines methods for interacting with storage
type StorageInterface interface {
	SaveURLWithExpiry(ctx context.Context, shortURL, longURL string, ttl time.Duration) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, bool)
}

// RedisStorage is the concrete implementation of StorageInterface
type RedisStorage struct {
	Client *redis.Client
}

// NewRedisStorage creates a new RedisStorage instance
func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{Client: client}
}

// SaveURLWithExpiry stores the mapping with a TTL
func (s *RedisStorage) SaveURLWithExpiry(ctx context.Context, shortURL, longURL string, ttl time.Duration) error {
	return s.Client.Set(ctx, shortURL, longURL, ttl).Err()
}

// GetOriginalURL retrieves the long URL for a given short URL
func (s *RedisStorage) GetOriginalURL(ctx context.Context, shortURL string) (string, bool) {
	longURL, err := s.Client.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		log.Printf("Redis error: %v", err)
		return "", false
	}
	return longURL, true
}
