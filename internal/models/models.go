package models

type ShortenURLRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenURLResponse struct {
	ShortURL string `json:"short_url"`
}
