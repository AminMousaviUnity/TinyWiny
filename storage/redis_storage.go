package storage

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() {
	// Get Redis host and port from environment variables
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" || redisPort == "" {
		log.Fatalf("Environment variabels REDIS_HOST or REDIS_PORT are not set")
	}

	// Connect to Redis
	rdb = redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})

	// Test the connection
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}

// SaveURL stores the mapping of short URL to long URL in Redis
func SaveURL(shortURL, longURL string) error {
	ctx := context.Background()
	return rdb.Set(ctx, shortURL, longURL, 0).Err()
}

// SaveURLWithExpiry stores the mapping with a time-to-live (TTL)
func SaveURLWithExpiry(shortURL, longURL string, ttl time.Duration) error {
	ctx := context.Background()
	return rdb.Set(ctx, shortURL, longURL, ttl).Err()
}

// GetOriginalURL retrieves the long URL for a given short URL from Redis
func GetOriginalURL(shortURL string) (string, bool) {
	ctx := context.Background()
	longURL, err := rdb.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		log.Printf("Redis error: %v", err)
		return "", false
	}
	return longURL, true
}
