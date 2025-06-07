package quote

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	utils "github.com/Korjick/go-http-quote/presentation/http"

	"github.com/Korjick/go-http-quote/application/service"
	"github.com/Korjick/go-http-quote/domain/quote/entity"
	"github.com/Korjick/go-http-quote/presentation/http/quote/dto"
)

type Controller struct {
	service *service.QuoteService
	prefix  string
}

func NewQuoteController(service *service.QuoteService, prefix string) *Controller {
	controller := &Controller{
		service: service,
		prefix:  prefix,
	}
	return controller
}

func (h *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.URL.Path = strings.TrimPrefix(r.URL.Path, h.prefix)
	switch r.Method {
	case http.MethodPost:
		h.createQuote(w, r)
	case http.MethodGet:
		switch {
		case strings.HasPrefix(r.URL.Path, "/random"):
			h.getRandomQuote(w, r)
		default:
			h.getQuotes(w, r)
		}
	case http.MethodDelete:
		h.deleteQuote(w, r)
	default:
		h.handleDomainError(w, errors.ErrUnsupported)
	}
}

func (h *Controller) handleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.ErrEmptyAuthor), errors.Is(err, entity.ErrEmptyText), errors.Is(err, entity.ErrQuoteNotFound):
		utils.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	default:
		utils.WriteJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Error: http.StatusText(http.StatusInternalServerError)})
	}
}

func (h *Controller) createQuote(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: http.StatusText(http.StatusBadRequest)})
		return
	}

	quote, err := h.service.CreateQuote(req.Author, req.Quote)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	response := dto.EntityToDTO(quote)
	utils.WriteJSON(w, http.StatusCreated, response)
}

func (h *Controller) getQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	var quotes []*entity.Quote
	var err error

	if author != "" {
		quotes, err = h.service.GetQuotesByAuthor(author)
	} else {
		quotes, err = h.service.GetAllQuotes()
	}

	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	response := dto.EntitiesToDTO(quotes)
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Controller) getRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := h.service.GetRandomQuote()
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	response := dto.EntityToDTO(quote)
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Controller) deleteQuote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid quote ID"})
		return
	}

	err = h.service.DeleteQuote(entity.QuoteID(id))
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
