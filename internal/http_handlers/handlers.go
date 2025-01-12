package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aminmousaviunity/TinyWiny/internal/models"
	"github.com/aminmousaviunity/TinyWiny/internal/services"
	"github.com/aminmousaviunity/TinyWiny/internal/storage"
)

// Handlers struct for dependency injection
type Handlers struct {
	BaseURL  string
	Storage  storage.StorageInterface
	Services services.ServiceInterface
}

// NewHandlers creates a new Handlers instance
func NewHandlers(baseURL string, storage storage.StorageInterface, services services.ServiceInterface) *Handlers {
	return &Handlers{
		BaseURL:  baseURL,
		Storage:  storage,
		Services: services,
	}
}

// ShortenURLHandler handles POST requests to create a short URL
func (h *Handlers) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the request body
	var req models.ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.LongURL == "" {
		http.Error(w, "Invalid JSON or missing long_url field", http.StatusBadRequest)
		return
	}

	// Generate short URL and save it
	shortURL := h.Services.GenerateShortURL(req.LongURL)
	err := h.Storage.SaveURLWithExpiry(r.Context(), shortURL, req.LongURL, 24*time.Hour) // Expires in 24 hours
	if err != nil {
		log.Printf("Error saving URL: %v", err) // Add this line
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	// Respond with the short URL
	resp := models.ShortenURLResponse{ShortURL: h.BaseURL + "/" + shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// RedirectHandler handles GET requests to redirect to the original URL
func (h *Handlers) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract short URL from the path
	shortURL := r.URL.Path[1:]

	// Lookup the original URL in Redis
	longURL, exists := h.Storage.GetOriginalURL(r.Context(), shortURL)
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, longURL, http.StatusFound)
}
