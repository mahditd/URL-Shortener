package dto

import "time"

type ShortenResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

type MetadataResponse struct {
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}
