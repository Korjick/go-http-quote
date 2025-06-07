package repository

import "github.com/Korjick/go-http-quote/domain/quote/entity"

type QuoteRepository interface {
	Create(author, text string) (*entity.Quote, error)
	GetAll() ([]*entity.Quote, error)
	GetByAuthor(author string) ([]*entity.Quote, error)
	GetRandom() (*entity.Quote, error)
	Delete(id entity.QuoteID) error
}
