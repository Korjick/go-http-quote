package dto_test

import (
	"github.com/Korjick/go-http-quote.git/domain/quote/entity"
	"github.com/Korjick/go-http-quote.git/presentation/http/quote/dto"
	"testing"
)

func TestEntityToDTO(t *testing.T) {
	quote, err := entity.NewQuote(42, "Albert Einstein", "Imagination is more important than knowledge.")
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	response := dto.EntityToDTO(quote)

	if response.ID != int64(quote.ID) {
		t.Errorf("EntityToDTO() ID = %v, want %v", response.ID, int64(quote.ID))
	}

	if response.Author != quote.Author {
		t.Errorf("EntityToDTO() Author = %v, want %v", response.Author, quote.Author)
	}

	if response.Quote != quote.Text {
		t.Errorf("EntityToDTO() Quote = %v, want %v", response.Quote, quote.Text)
	}

	if !response.CreatedAt.Equal(quote.CreatedAt) {
		t.Errorf("EntityToDTO() CreatedAt = %v, want %v", response.CreatedAt, quote.CreatedAt)
	}
}

func TestEntitiesToDTO(t *testing.T) {
	quote1, err := entity.NewQuote(1, "Einstein", "Quote 1")
	if err != nil {
		t.Fatalf("Failed to create quote1: %v", err)
	}

	quote2, err := entity.NewQuote(2, "Jobs", "Quote 2")
	if err != nil {
		t.Fatalf("Failed to create quote2: %v", err)
	}

	quotes := []*entity.Quote{quote1, quote2}

	responses := dto.EntitiesToDTO(quotes)

	if len(responses) != len(quotes) {
		t.Errorf("EntitiesToDTO() returned %d DTOs, want %d", len(responses), len(quotes))
	}

	for i, response := range responses {
		quote := quotes[i]

		if response.ID != int64(quote.ID) {
			t.Errorf("EntitiesToDTO()[%d] ID = %v, want %v", i, response.ID, int64(quote.ID))
		}

		if response.Author != quote.Author {
			t.Errorf("EntitiesToDTO()[%d] Author = %v, want %v", i, response.Author, quote.Author)
		}

		if response.Quote != quote.Text {
			t.Errorf("EntitiesToDTO()[%d] Quote = %v, want %v", i, response.Quote, quote.Text)
		}

		if !response.CreatedAt.Equal(quote.CreatedAt) {
			t.Errorf("EntitiesToDTO()[%d] CreatedAt = %v, want %v", i, response.CreatedAt, quote.CreatedAt)
		}
	}
}

func TestEntitiesToDTO_EmptySlice(t *testing.T) {
	var quotes []*entity.Quote
	responses := dto.EntitiesToDTO(quotes)

	if len(responses) != 0 {
		t.Errorf("EntitiesToDTO() with empty slice returned %d DTOs, want 0", len(responses))
	}
}

func TestEntitiesToDTO_NilSlice(t *testing.T) {
	var quotes []*entity.Quote
	responses := dto.EntitiesToDTO(quotes)

	if len(responses) != 0 {
		t.Errorf("EntitiesToDTO() with nil slice returned %d DTOs, want 0", len(responses))
	}
}
