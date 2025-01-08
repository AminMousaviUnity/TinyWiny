package handlers

import (
	"TinyWiny/app/storage"
	"encoding/json"
	"net/http"
)

// ShortenURLRequest represents the JSON payload for the shorten request
type ShortenURLRequest struct {
	LongURL string `json:"long_url"`
}

// ShortenURLResponse represents the JSON response for the shorten request
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

	// Generate short URL and save the mapping
	shortURL := storage.GenerateShortURL()
	storage.SaveURL(shortURL, req.LongURL)

	// Respond with the short URL
	resp := ShortenURLResponse{ShortURL: storage.BaseURL + shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// RedirectHandler handles GET requests to redirect to the original URL
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract short URL from the path
	shortURL := r.URL.Path[1:]

	// Lookup the original URL
	longURL, exists := storage.GetOriginalURL(shortURL)
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, longURL, http.StatusFound)
}
