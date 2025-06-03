package service

import (
	"github.com/Korjick/go-http-quote.git/domain/quote/entity"
	"github.com/Korjick/go-http-quote.git/domain/quote/repository"
)

type QuoteService struct {
	repo repository.QuoteRepository
}

func NewQuoteService(repo repository.QuoteRepository) *QuoteService {
	return &QuoteService{
		repo: repo,
	}
}

func (s *QuoteService) CreateQuote(author, text string) (*entity.Quote, error) {
	return s.repo.Create(author, text)
}

func (s *QuoteService) GetAllQuotes() ([]*entity.Quote, error) {
	return s.repo.GetAll()
}

func (s *QuoteService) GetQuotesByAuthor(author string) ([]*entity.Quote, error) {
	return s.repo.GetByAuthor(author)
}

func (s *QuoteService) GetRandomQuote() (*entity.Quote, error) {
	return s.repo.GetRandom()
}

func (s *QuoteService) DeleteQuote(id entity.QuoteID) error {
	return s.repo.Delete(id)
}
