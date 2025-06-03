package main

import (
	"github.com/Korjick/go-http-quote.git/infrastructure/repository/in_memory"
	"log"
	"net/http"

	"github.com/Korjick/go-http-quote.git/application/service"
	"github.com/Korjick/go-http-quote.git/presentation/http/quote"
)

func main() {
	repo := in_memory.NewInMemoryQuoteRepository()
	quoteService := service.NewQuoteService(repo)

	quotePrefix := "/quotes"
	quoteHandler := quote.NewQuoteController(quoteService, quotePrefix)

	http.Handle(quotePrefix, quoteHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
