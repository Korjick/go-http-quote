package dto

import (
	"github.com/Korjick/go-http-quote.git/domain/quote/entity"
)

func EntityToDTO(quote *entity.Quote) QuoteResponse {
	return QuoteResponse{
		ID:        int64(quote.ID),
		Author:    quote.Author,
		Quote:     quote.Text,
		CreatedAt: quote.CreatedAt,
	}
}

func EntitiesToDTO(quotes []*entity.Quote) []QuoteResponse {
	dtos := make([]QuoteResponse, len(quotes))
	for i, quote := range quotes {
		dtos[i] = EntityToDTO(quote)
	}
	return dtos
}
