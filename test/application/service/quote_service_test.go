package service_test

import (
	"errors"
	"testing"

	"github.com/Korjick/go-http-quote/application/service"
	"github.com/Korjick/go-http-quote/domain/quote/entity"
	"github.com/Korjick/go-http-quote/infrastructure/repository/in_memory"
)

func TestQuoteService_CreateQuoteValidInput(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()
	svc := service.NewQuoteService(repo)

	author := "Albert Einstein"
	text := "Imagination is more important than knowledge."

	quote, err := svc.CreateQuote(author, text)

	if err != nil {
		t.Errorf("CreateQuote() unexpected error = %v", err)
		return
	}

	if quote == nil {
		t.Error("CreateQuote() returned nil quote")
		return
	}

	if quote.Author != author {
		t.Errorf("Quote.Author = %v, want %v", quote.Author, author)
	}

	if quote.Text != text {
		t.Errorf("Quote.Text = %v, want %v", quote.Text, text)
	}

	if quote.ID <= 0 {
		t.Error("Quote.ID should be positive")
	}
}

func TestQuoteService_CreateQuoteEmptyAuthor(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()
	svc := service.NewQuoteService(repo)

	author := ""
	text := "Some text"

	quote, err := svc.CreateQuote(author, text)

	if err == nil {
		t.Error("CreateQuote() expected error but got none")
		return
	}

	if !errors.Is(err, entity.ErrEmptyAuthor) {
		t.Errorf("CreateQuote() error = %v, want %v", err, entity.ErrEmptyAuthor)
	}

	if quote != nil {
		t.Error("CreateQuote() should return nil quote on error")
	}
}

func TestQuoteService_CreateQuoteEmptyText(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()
	svc := service.NewQuoteService(repo)

	author := "Author"
	text := ""

	quote, err := svc.CreateQuote(author, text)

	if err == nil {
		t.Error("CreateQuote() expected error but got none")
		return
	}

	if !errors.Is(err, entity.ErrEmptyText) {
		t.Errorf("CreateQuote() error = %v, want %v", err, entity.ErrEmptyText)
	}

	if quote != nil {
		t.Error("CreateQuote() should return nil quote on error")
	}
}

func TestQuoteService_GetAllQuotes(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()
	svc := service.NewQuoteService(repo)

	quotes, err := svc.GetAllQuotes()
	if err != nil {
		t.Errorf("GetAllQuotes() error = %v, want nil", err)
	}
	if len(quotes) != 0 {
		t.Errorf("GetAllQuotes() returned %d quotes, want 0", len(quotes))
	}

	_, err = svc.CreateQuote("Author 1", "Quote 1")
	if err != nil {
		t.Fatalf("CreateQuote() error = %v", err)
	}

	_, err = svc.CreateQuote("Author 2", "Quote 2")
	if err != nil {
		t.Fatalf("CreateQuote() error = %v", err)
	}

	quotes, err = svc.GetAllQuotes()
	if err != nil {
		t.Errorf("GetAllQuotes() error = %v, want nil", err)
	}
	if len(quotes) != 2 {
		t.Errorf("GetAllQuotes() returned %d quotes, want 2", len(quotes))
	}
}

func TestQuoteService_GetQuotesByAuthor(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()
	svc := service.NewQuoteService(repo)

	_, err := svc.CreateQuote("Einstein", "Quote 1")
	if err != nil {
		t.Fatalf("CreateQuote() error = %v", err)
	}

	_, err = svc.CreateQuote("Einstein", "Quote 2")
	if err != nil {
		t.Fatalf("CreateQuote() error = %v", err)
	}

	_, err = svc.CreateQuote("Jobs", "Quote 3")
	if err != nil {
		t.Fatalf("CreateQuote() error = %v", err)
	}

	einsteinQuotes, err := svc.GetQuotesByAuthor("Einstein")
	if err != nil {
		t.Errorf("GetQuotesByAuthor() error = %v, want nil", err)
	}
	if len(einsteinQuotes) != 2 {
		t.Errorf("GetQuotesByAuthor('Einstein') returned %d quotes, want 2", len(einsteinQuotes))
	}

	jobsQuotes, err := svc.GetQuotesByAuthor("Jobs")
	if err != nil {
		t.Errorf("GetQuotesByAuthor() error = %v, want nil", err)
	}
	if len(jobsQuotes) != 1 {
		t.Errorf("GetQuotesByAuthor('Jobs') returned %d quotes, want 1", len(jobsQuotes))
	}

	nonExistentQuotes, err := svc.GetQuotesByAuthor("Non-existent")
	if err != nil {
		t.Errorf("GetQuotesByAuthor() error = %v, want nil", err)
	}
	if len(nonExistentQuotes) != 0 {
		t.Errorf("GetQuotesByAuthor('Non-existent') returned %d quotes, want 0", len(nonExistentQuotes))
	}
}

func TestQuoteService_GetRandomQuote(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()
	svc := service.NewQuoteService(repo)

	_, err := svc.GetRandomQuote()
	if !errors.Is(err, entity.ErrQuoteNotFound) {
		t.Errorf("GetRandomQuote() error = %v, want %v", err, entity.ErrQuoteNotFound)
	}

	createdQuote, err := svc.CreateQuote("Test Author", "Test Quote")
	if err != nil {
		t.Fatalf("CreateQuote() error = %v", err)
	}

	randomQuote, err := svc.GetRandomQuote()
	if err != nil {
		t.Errorf("GetRandomQuote() error = %v, want nil", err)
	}

	if randomQuote.ID != createdQuote.ID {
		t.Errorf("GetRandomQuote() returned quote with ID %v, want %v", randomQuote.ID, createdQuote.ID)
	}

	for i := 2; i <= 10; i++ {
		_, err = svc.CreateQuote("Author", "Quote")
		if err != nil {
			t.Fatalf("CreateQuote() error = %v", err)
		}
	}

	for i := 0; i < 5; i++ {
		randomQuote, err = svc.GetRandomQuote()
		if err != nil {
			t.Errorf("GetRandomQuote() error = %v, want nil", err)
		}
		if randomQuote == nil {
			t.Error("GetRandomQuote() returned nil quote")
		}
	}
}

func TestQuoteService_DeleteQuote(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()
	svc := service.NewQuoteService(repo)

	err := svc.DeleteQuote(999)
	if !errors.Is(err, entity.ErrQuoteNotFound) {
		t.Errorf("DeleteQuote() error = %v, want %v", err, entity.ErrQuoteNotFound)
	}

	createdQuote, err := svc.CreateQuote("Test Author", "Test Quote")
	if err != nil {
		t.Fatalf("CreateQuote() error = %v", err)
	}

	err = svc.DeleteQuote(createdQuote.ID)
	if err != nil {
		t.Errorf("DeleteQuote() error = %v, want nil", err)
	}

	quotes, err := svc.GetAllQuotes()
	if err != nil {
		t.Errorf("GetAllQuotes() error = %v", err)
	}
	if len(quotes) != 0 {
		t.Errorf("After deletion, GetAllQuotes() returned %d quotes, want 0", len(quotes))
	}
}
