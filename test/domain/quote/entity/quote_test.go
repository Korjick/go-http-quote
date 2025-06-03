package entity_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Korjick/go-http-quote.git/domain/quote/entity"
)

func TestNewQuoteValidInput(t *testing.T) {
	id := entity.QuoteID(1)
	author := "Albert Einstein"
	text := "Imagination is more important than knowledge."

	quote, err := entity.NewQuote(id, author, text)

	if err != nil {
		t.Errorf("NewQuote() unexpected error = %v", err)
		return
	}

	if quote == nil {
		t.Error("NewQuote() returned nil quote")
		return
	}

	if quote.ID != id {
		t.Errorf("Quote.ID = %v, want %v", quote.ID, id)
	}

	if quote.Author != author {
		t.Errorf("Quote.Author = %v, want %v", quote.Author, author)
	}

	if quote.Text != text {
		t.Errorf("Quote.Text = %v, want %v", quote.Text, text)
	}

	if quote.CreatedAt.IsZero() {
		t.Error("Quote.CreatedAt should not be zero")
	}

	if time.Since(quote.CreatedAt) > time.Second {
		t.Error("Quote.CreatedAt should be recent")
	}
}

func TestNewQuoteEmptyAuthor(t *testing.T) {
	id := entity.QuoteID(1)
	author := ""
	text := "Some text"

	quote, err := entity.NewQuote(id, author, text)

	if err == nil {
		t.Error("NewQuote() expected error but got none")
		return
	}

	if !errors.Is(err, entity.ErrEmptyAuthor) {
		t.Errorf("NewQuote() error = %v, want %v", err, entity.ErrEmptyAuthor)
	}

	if quote != nil {
		t.Error("NewQuote() should return nil quote on error")
	}
}

func TestNewQuoteEmptyText(t *testing.T) {
	id := entity.QuoteID(1)
	author := "Author"
	text := ""

	quote, err := entity.NewQuote(id, author, text)

	if err == nil {
		t.Error("NewQuote() expected error but got none")
		return
	}

	if !errors.Is(err, entity.ErrEmptyText) {
		t.Errorf("NewQuote() error = %v, want %v", err, entity.ErrEmptyText)
	}

	if quote != nil {
		t.Error("NewQuote() should return nil quote on error")
	}
}

func TestNewQuoteWhitespaceOnlyAuthor(t *testing.T) {
	id := entity.QuoteID(1)
	author := "   "
	text := "Some text"

	quote, err := entity.NewQuote(id, author, text)

	if err == nil {
		t.Error("NewQuote() expected error but got none")
		return
	}

	if !errors.Is(err, entity.ErrEmptyAuthor) {
		t.Errorf("NewQuote() error = %v, want %v", err, entity.ErrEmptyAuthor)
	}

	if quote != nil {
		t.Error("NewQuote() should return nil quote on error")
	}
}

func TestNewQuoteWhitespaceOnlyText(t *testing.T) {
	id := entity.QuoteID(1)
	author := "Author"
	text := "   "

	quote, err := entity.NewQuote(id, author, text)

	if err == nil {
		t.Error("NewQuote() expected error but got none")
		return
	}

	if !errors.Is(err, entity.ErrEmptyText) {
		t.Errorf("NewQuote() error = %v, want %v", err, entity.ErrEmptyText)
	}

	if quote != nil {
		t.Error("NewQuote() should return nil quote on error")
	}
}
