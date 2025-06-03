package entity

import "errors"

var (
	ErrEmptyAuthor   = errors.New("author cannot be empty")
	ErrEmptyText     = errors.New("quote text cannot be empty")
	ErrQuoteNotFound = errors.New("quote not found")
)
