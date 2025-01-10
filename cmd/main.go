package main

import (
	"fmt"
	"log"
	"net/http"

	handlers "github.com/aminmousaviunity/TinyWiny/internal/http_handlers"
	"github.com/aminmousaviunity/TinyWiny/internal/storage"
)

func main() {
	// Initialize in-memory storage
	storage.InitRedisPrd()

	// Define routes
	http.HandleFunc("/shorten", handlers.ShortenURLHandler) // POST: Create a short URL
	http.HandleFunc("/", handlers.RedirectHandler)          // GET: Redirect to the original

	// Start the server
	addr := ":8888"
	fmt.Printf("Starting server on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
