package dto

import "time"

type QuoteResponse struct {
	ID        int64     `json:"id"`
	Author    string    `json:"author"`
	Quote     string    `json:"quote"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
