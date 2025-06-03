package entity

import (
	"strings"
	"time"
)

type QuoteID int64

type Quote struct {
	ID        QuoteID   `json:"id"`
	Author    string    `json:"author"`
	Text      string    `json:"quote"`
	CreatedAt time.Time `json:"created_at"`
}

func NewQuote(id QuoteID, author, text string) (*Quote, error) {
	quote := &Quote{
		ID:        id,
		Author:    author,
		Text:      text,
		CreatedAt: time.Now(),
	}

	if err := quote.validate(); err != nil {
		return nil, err
	}

	return quote, nil
}

func (q *Quote) validate() error {
	if strings.TrimSpace(q.Author) == "" {
		return ErrEmptyAuthor
	}
	if strings.TrimSpace(q.Text) == "" {
		return ErrEmptyText
	}
	return nil
}
