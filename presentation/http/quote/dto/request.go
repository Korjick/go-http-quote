package dto

type CreateQuoteRequest struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}
