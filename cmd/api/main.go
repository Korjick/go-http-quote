package main

import (
	"log"
	"net/http"

	"github.com/Korjick/go-http-quote/infrastructure/repository/in_memory"

	"github.com/Korjick/go-http-quote/application/service"
	"github.com/Korjick/go-http-quote/presentation/http/quote"
)

func main() {
	repo := in_memory.NewInMemoryQuoteRepository()
	quoteService := service.NewQuoteService(repo)

	quotePrefix := "/quotes"
	quoteHandler := quote.NewQuoteController(quoteService, quotePrefix)

	http.Handle(quotePrefix+"/", quoteHandler)
	http.Handle(quotePrefix, quoteHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
