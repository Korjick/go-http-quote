package in_memory_test

import (
	"errors"
	"testing"

	"github.com/Korjick/go-http-quote/domain/quote/entity"
	"github.com/Korjick/go-http-quote/infrastructure/repository/in_memory"
)

func TestInMemoryQuoteRepository_Create(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()

	author := "Test Author"
	text := "Test Quote"

	quote, err := repo.Create(author, text)
	if err != nil {
		t.Errorf("Create() error = %v, want nil", err)
	}

	if quote == nil {
		t.Error("Create() returned nil quote")
		return
	}

	if quote.ID != 1 {
		t.Errorf("Created quote ID = %v, want 1", quote.ID)
	}
	if quote.Author != author {
		t.Errorf("Created quote Author = %v, want %v", quote.Author, author)
	}

	if quote.Text != text {
		t.Errorf("Created quote Text = %v, want %v", quote.Text, text)
	}

	allQuotes, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll() error = %v, want nil", err)
	}

	if len(allQuotes) != 1 {
		t.Errorf("GetAll() returned %d quotes, want 1", len(allQuotes))
	}
}

func TestInMemoryQuoteRepository_CreateInvalidInput(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()

	_, err := repo.Create("", "Some text")
	if err == nil {
		t.Error("Create() with empty author should return error")
	}

	_, err = repo.Create("Author", "")
	if err == nil {
		t.Error("Create() with empty text should return error")
	}
}

func TestInMemoryQuoteRepository_GetAll(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()

	quotes, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll() error = %v, want nil", err)
	}

	if len(quotes) != 0 {
		t.Errorf("GetAll() returned %d quotes, want 0", len(quotes))
	}

	_, _ = repo.Create("Author 1", "Quote 1")
	_, _ = repo.Create("Author 2", "Quote 2")

	quotes, err = repo.GetAll()
	if err != nil {
		t.Errorf("GetAll() error = %v, want nil", err)
	}

	if len(quotes) != 2 {
		t.Errorf("GetAll() returned %d quotes, want 2", len(quotes))
	}

	if quotes[0].Text != "Quote 1" {
		t.Errorf("First quote text = %v, want 'Quote 1'", quotes[0].Text)
	}

	if quotes[1].Text != "Quote 2" {
		t.Errorf("Second quote text = %v, want 'Quote 2'", quotes[1].Text)
	}
}

func TestInMemoryQuoteRepository_GetByAuthor(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()

	_, _ = repo.Create("Albert Einstein", "Quote 1")
	_, _ = repo.Create("Steve Jobs", "Quote 2")
	_, _ = repo.Create("Albert Einstein", "Quote 3")

	quotes, err := repo.GetByAuthor("Albert Einstein")
	if err != nil {
		t.Errorf("GetByAuthor() error = %v, want nil", err)
	}

	if len(quotes) != 2 {
		t.Errorf("GetByAuthor() returned %d quotes, want 2", len(quotes))
	}

	quotes, err = repo.GetByAuthor("albert einstein")
	if err != nil {
		t.Errorf("GetByAuthor() error = %v, want nil", err)
	}

	if len(quotes) != 2 {
		t.Errorf("GetByAuthor() case-insensitive returned %d quotes, want 2", len(quotes))
	}

	quotes, err = repo.GetByAuthor("Non-existent Author")
	if err != nil {
		t.Errorf("GetByAuthor() error = %v, want nil", err)
	}

	if len(quotes) != 0 {
		t.Errorf("GetByAuthor() for non-existent author returned %d quotes, want 0", len(quotes))
	}
}

func TestInMemoryQuoteRepository_GetRandom(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()

	_, err := repo.GetRandom()
	if !errors.Is(err, entity.ErrQuoteNotFound) {
		t.Errorf("GetRandom() on empty repo error = %v, want %v", err, entity.ErrQuoteNotFound)
	}

	_, _ = repo.Create("Author 1", "Quote 1")
	_, _ = repo.Create("Author 2", "Quote 2")
	_, _ = repo.Create("Author 3", "Quote 3")

	quote, err := repo.GetRandom()
	if err != nil {
		t.Errorf("GetRandom() error = %v, want nil", err)
	}

	if quote == nil {
		t.Error("GetRandom() returned nil quote")
		return
	}

	validIDs := []entity.QuoteID{1, 2, 3}
	found := false
	for _, id := range validIDs {
		if quote.ID == id {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("GetRandom() returned quote with ID %v, want one of %v", quote.ID, validIDs)
	}
}

func TestInMemoryQuoteRepository_Delete(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()

	err := repo.Delete(999)
	if !errors.Is(err, entity.ErrQuoteNotFound) {
		t.Errorf("Delete() on empty repo error = %v, want %v", err, entity.ErrQuoteNotFound)
	}

	quote1, _ := repo.Create("Author 1", "Quote 1")
	quote2, _ := repo.Create("Author 2", "Quote 2")
	quote3, _ := repo.Create("Author 3", "Quote 3")

	err = repo.Delete(quote2.ID)
	if err != nil {
		t.Errorf("Delete() error = %v, want nil", err)
	}

	allQuotes, _ := repo.GetAll()
	if len(allQuotes) != 2 {
		t.Errorf("After deletion, found %d quotes, want 2", len(allQuotes))
	}

	remainingIDs := []entity.QuoteID{allQuotes[0].ID, allQuotes[1].ID}
	expectedIDs := []entity.QuoteID{quote1.ID, quote3.ID}

	for _, expectedID := range expectedIDs {
		found := false
		for _, remainingID := range remainingIDs {
			if remainingID == expectedID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected quote ID %v not found after deletion", expectedID)
		}
	}

	err = repo.Delete(quote2.ID)
	if !errors.Is(err, entity.ErrQuoteNotFound) {
		t.Errorf("Delete() already deleted quote error = %v, want %v", err, entity.ErrQuoteNotFound)
	}
}

func TestInMemoryQuoteRepository_ThreadSafety(t *testing.T) {
	repo := in_memory.NewInMemoryQuoteRepository()

	done := make(chan bool, 10)

	for i := 0; i < 5; i++ {
		go func(id int) {
			_, _ = repo.Create("Author", "Quote from goroutine")
			done <- true
		}(i)
	}

	for i := 0; i < 5; i++ {
		go func() {
			_, _ = repo.GetAll()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	quotes, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll() after concurrent operations error = %v", err)
	}

	if len(quotes) != 5 {
		t.Errorf("After concurrent operations, found %d quotes, want 5", len(quotes))
	}
}
