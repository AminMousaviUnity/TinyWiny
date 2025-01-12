package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	handlers "github.com/aminmousaviunity/TinyWiny/internal/http_handlers"
	"github.com/aminmousaviunity/TinyWiny/internal/services"
	"github.com/aminmousaviunity/TinyWiny/internal/storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize Redis client
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" || redisPort == "" {
		log.Fatal("Environment variables REDIS_HOST and REDIS_PORT must be set")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})

	// Initialize storage, services, and handlers
	storage := storage.NewRedisStorage(redisClient)
	services := &services.Service{}
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Fatal("BaseURL must be set")
	}
	h := handlers.NewHandlers(baseURL, storage, services)

	// Define routes
	http.HandleFunc("/shorten", h.ShortenURLHandler) // POST: Create a short URL
	http.HandleFunc("/", h.RedirectHandler)          // GET: Redirect to the original

	// Start the server
	addr := ":8888"
	fmt.Printf("Starting server on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
