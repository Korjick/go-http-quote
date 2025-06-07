package in_memory

import (
	"github.com/Korjick/go-http-quote/domain/quote/repository"
	"math/rand"
	"strings"
	"sync"

	"github.com/Korjick/go-http-quote/domain/quote/entity"
)

type inMemoryQuoteRepository struct {
	quotes []*entity.Quote
	mutex  sync.RWMutex
}

func NewInMemoryQuoteRepository() repository.QuoteRepository {
	return &inMemoryQuoteRepository{
		quotes: make([]*entity.Quote, 0),
	}
}

func (r *inMemoryQuoteRepository) Create(author, text string) (*entity.Quote, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id := entity.QuoteID(len(r.quotes) + 1)
	quote, err := entity.NewQuote(id, author, text)
	if err != nil {
		return nil, err
	}

	r.quotes = append(r.quotes, quote)
	return quote, nil
}

func (r *inMemoryQuoteRepository) GetAll() ([]*entity.Quote, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*entity.Quote, len(r.quotes))
	copy(result, r.quotes)
	return result, nil
}

func (r *inMemoryQuoteRepository) GetByAuthor(author string) ([]*entity.Quote, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*entity.Quote
	for i := range r.quotes {
		if strings.EqualFold(r.quotes[i].Author, author) {
			result = append(result, r.quotes[i])
		}
	}
	return result, nil
}

func (r *inMemoryQuoteRepository) GetRandom() (*entity.Quote, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(r.quotes) == 0 {
		return nil, entity.ErrQuoteNotFound
	}

	index := rand.Intn(len(r.quotes))
	return r.quotes[index], nil
}

func (r *inMemoryQuoteRepository) Delete(id entity.QuoteID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, quote := range r.quotes {
		if quote.ID == id {
			r.quotes = append(r.quotes[:i], r.quotes[i+1:]...)
			return nil
		}
	}
	return entity.ErrQuoteNotFound
}
