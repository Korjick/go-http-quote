package quote_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Korjick/go-http-quote.git/application/service"
	"github.com/Korjick/go-http-quote.git/infrastructure/repository/in_memory"
	"github.com/Korjick/go-http-quote.git/presentation/http/quote"
	"github.com/Korjick/go-http-quote.git/presentation/http/quote/dto"
)

func setupTestController() *quote.Controller {
	repo := in_memory.NewInMemoryQuoteRepository()
	svc := service.NewQuoteService(repo)
	return quote.NewQuoteController(svc, "/quotes")
}

func TestController_CreateQuoteValidInput(t *testing.T) {
	controller := setupTestController()

	requestBody := dto.CreateQuoteRequest{
		Author: "Albert Einstein",
		Quote:  "Imagination is more important than knowledge.",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("CreateQuote() status = %v, want %v", w.Code, http.StatusCreated)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("CreateQuote() Content-Type = %v, want application/json", contentType)
	}

	var response dto.QuoteResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Author != requestBody.Author {
		t.Errorf("Response Author = %v, want %v", response.Author, requestBody.Author)
	}

	if response.Quote != requestBody.Quote {
		t.Errorf("Response Quote = %v, want %v", response.Quote, requestBody.Quote)
	}

	if response.ID <= 0 {
		t.Error("Response ID should be positive")
	}

	if response.CreatedAt.IsZero() {
		t.Error("Response CreatedAt should not be zero")
	}
}

func TestController_CreateQuoteEmptyAuthor(t *testing.T) {
	controller := setupTestController()

	requestBody := dto.CreateQuoteRequest{
		Author: "",
		Quote:  "Some quote",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("CreateQuote() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("CreateQuote() Content-Type = %v, want application/json", contentType)
	}

	var errorResp dto.ErrorResponse
	err = json.NewDecoder(w.Body).Decode(&errorResp)
	if err != nil {
		t.Errorf("Failed to decode error response: %v", err)
	}
	if errorResp.Error == "" {
		t.Error("Expected error message but got empty string")
	}
}

func TestController_CreateQuoteEmptyText(t *testing.T) {
	controller := setupTestController()

	requestBody := dto.CreateQuoteRequest{
		Author: "Some Author",
		Quote:  "",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("CreateQuote() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("CreateQuote() Content-Type = %v, want application/json", contentType)
	}

	var errorResp dto.ErrorResponse
	err = json.NewDecoder(w.Body).Decode(&errorResp)
	if err != nil {
		t.Errorf("Failed to decode error response: %v", err)
	}
	if errorResp.Error == "" {
		t.Error("Expected error message but got empty string")
	}
}

func TestController_GetAllQuotes(t *testing.T) {
	controller := setupTestController()

	quotes := []dto.CreateQuoteRequest{
		{Author: "Einstein", Quote: "Quote 1"},
		{Author: "Einstein", Quote: "Quote 2"},
		{Author: "Jobs", Quote: "Quote 3"},
	}

	for _, q := range quotes {
		jsonBody, _ := json.Marshal(q)
		req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		controller.ServeHTTP(w, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetQuotes() status = %v, want %v", w.Code, http.StatusOK)
	}

	var responses []dto.QuoteResponse
	err := json.NewDecoder(w.Body).Decode(&responses)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(responses) != 3 {
		t.Errorf("GetQuotes() returned %d quotes, want 3", len(responses))
	}
}

func TestController_GetQuotesByAuthorEinstein(t *testing.T) {
	controller := setupTestController()

	quotes := []dto.CreateQuoteRequest{
		{Author: "Einstein", Quote: "Quote 1"},
		{Author: "Einstein", Quote: "Quote 2"},
		{Author: "Jobs", Quote: "Quote 3"},
	}

	for _, q := range quotes {
		jsonBody, _ := json.Marshal(q)
		req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		controller.ServeHTTP(w, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/quotes?author=Einstein", nil)
	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetQuotes() status = %v, want %v", w.Code, http.StatusOK)
	}

	var responses []dto.QuoteResponse
	err := json.NewDecoder(w.Body).Decode(&responses)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(responses) != 2 {
		t.Errorf("GetQuotes() returned %d quotes, want 2", len(responses))
	}
}

func TestController_GetQuotesByAuthorJobs(t *testing.T) {
	controller := setupTestController()

	quotes := []dto.CreateQuoteRequest{
		{Author: "Einstein", Quote: "Quote 1"},
		{Author: "Einstein", Quote: "Quote 2"},
		{Author: "Jobs", Quote: "Quote 3"},
	}

	for _, q := range quotes {
		jsonBody, _ := json.Marshal(q)
		req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		controller.ServeHTTP(w, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/quotes?author=Jobs", nil)
	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetQuotes() status = %v, want %v", w.Code, http.StatusOK)
	}

	var responses []dto.QuoteResponse
	err := json.NewDecoder(w.Body).Decode(&responses)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(responses) != 1 {
		t.Errorf("GetQuotes() returned %d quotes, want 1", len(responses))
	}
}

func TestController_GetQuotesByNonExistentAuthor(t *testing.T) {
	controller := setupTestController()

	quotes := []dto.CreateQuoteRequest{
		{Author: "Einstein", Quote: "Quote 1"},
		{Author: "Einstein", Quote: "Quote 2"},
		{Author: "Jobs", Quote: "Quote 3"},
	}

	for _, q := range quotes {
		jsonBody, _ := json.Marshal(q)
		req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		controller.ServeHTTP(w, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/quotes?author=NonExistent", nil)
	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetQuotes() status = %v, want %v", w.Code, http.StatusOK)
	}

	var responses []dto.QuoteResponse
	err := json.NewDecoder(w.Body).Decode(&responses)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(responses) != 0 {
		t.Errorf("GetQuotes() returned %d quotes, want 0", len(responses))
	}
}

func TestController_GetRandomQuote(t *testing.T) {
	controller := setupTestController()

	req := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("GetRandomQuote() on empty repo status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	createReq := dto.CreateQuoteRequest{
		Author: "Test Author",
		Quote:  "Test Quote",
	}
	jsonBody, _ := json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(jsonBody))
	w = httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	w = httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetRandomQuote() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response dto.QuoteResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Author != createReq.Author {
		t.Errorf("Random quote Author = %v, want %v", response.Author, createReq.Author)
	}

	if response.Quote != createReq.Quote {
		t.Errorf("Random quote Quote = %v, want %v", response.Quote, createReq.Quote)
	}
}

func TestController_DeleteQuote(t *testing.T) {
	controller := setupTestController()

	req := httptest.NewRequest(http.MethodDelete, "/quotes/999", nil)
	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("DeleteQuote() non-existent status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	createReq := dto.CreateQuoteRequest{
		Author: "Test Author",
		Quote:  "Test Quote",
	}
	jsonBody, _ := json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(jsonBody))
	w = httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	var createResponse dto.QuoteResponse
	_ = json.NewDecoder(w.Body).Decode(&createResponse)

	deleteURL := "/quotes/" + strconv.FormatInt(createResponse.ID, 10)
	req = httptest.NewRequest(http.MethodDelete, deleteURL, nil)
	w = httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("DeleteQuote() status = %v, want %v", w.Code, http.StatusNoContent)
	}

	req = httptest.NewRequest(http.MethodGet, "/quotes", nil)
	w = httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	var quotes []dto.QuoteResponse
	_ = json.NewDecoder(w.Body).Decode(&quotes)

	if len(quotes) != 0 {
		t.Errorf("After deletion, found %d quotes, want 0", len(quotes))
	}
}

func TestController_PutMethodNotAllowed(t *testing.T) {
	controller := setupTestController()

	req := httptest.NewRequest(http.MethodPut, "/quotes", nil)
	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("PUT method status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestController_PatchMethodNotAllowed(t *testing.T) {
	controller := setupTestController()

	req := httptest.NewRequest(http.MethodPatch, "/quotes", nil)
	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("PATCH method status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}

func TestController_InvalidJSON(t *testing.T) {
	controller := setupTestController()

	req := httptest.NewRequest(http.MethodPost, "/quotes", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Invalid JSON status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	var errorResp dto.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&errorResp)
	if err != nil {
		t.Errorf("Failed to decode error response: %v", err)
	}

	if errorResp.Error == "" {
		t.Error("Expected error message but got empty string")
	}
}

func TestController_InvalidQuoteID(t *testing.T) {
	controller := setupTestController()

	req := httptest.NewRequest(http.MethodDelete, "/quotes/invalid", nil)
	w := httptest.NewRecorder()
	controller.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Invalid quote ID status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	var errorResp dto.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&errorResp)
	if err != nil {
		t.Errorf("Failed to decode error response: %v", err)
	}

	if !strings.Contains(errorResp.Error, "Invalid quote ID") {
		t.Errorf("Expected 'Invalid quote ID' error, got: %s", errorResp.Error)
	}
}
