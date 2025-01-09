package handlers

import (
	"TinyWiny/storage"
	"encoding/json"
	"net/http"
	"time"
)

type ShortenURLRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenURLResponse struct {
	ShortURL string `json:"short_url"`
}

// ShortenURLHandler handles POST requests to create a short URL
func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.LongURL == "" {
		http.Error(w, "Invalid JSON or missing long_url field", http.StatusBadRequest)
		return
	}

	// Generate short URL and save the mapping in Redis
	shortURL := storage.GenerateShortURL(req.LongURL)
	err := storage.SaveURLWithExpiry(shortURL, req.LongURL, 24*time.Hour) // Expires in 24 hours
	if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	// Respond with the short URL
	resp := ShortenURLResponse{ShortURL: "http://localhost:8888/" + shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// RedirectHandler handles GET requests to redirect to the original URL
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract short URL from the path
	shortURL := r.URL.Path[1:]

	// Lookup the original URL in Redis
	longURL, exists := storage.GetOriginalURL(shortURL)
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, longURL, http.StatusFound)
}
