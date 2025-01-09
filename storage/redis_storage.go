package storage

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb RedisClient

// RedisClient defines the interface for interacting with Redis
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Ping(ctx context.Context) *redis.StatusCmd
}

// InitRedis sets the global Redis client (used for both production and testing)
func InitRedis(client RedisClient) {
	rdb = client
}

// InitRedisPrd initializes the Redis client for prd use
func InitRedisPrd() {
	// Get Redis host and port from env variabels
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" || redisPort == "" {
		log.Fatalf("Environment variables REDIS_HOST or REDIS_PORT are not set")
	}

	// Create a read Redis client
	client := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})

	// Test the connection
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}

	// Assign the real client
	InitRedis(client)
}

// SaveURL stores the mapping of short URL to long URL in Redis
func SaveURL(shortURL, longURL string) error {
	ctx := context.Background()
	return rdb.Set(ctx, shortURL, longURL, 0).Err()
}

// SaveURLWithExpiry stores the mapping with a TTL
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
