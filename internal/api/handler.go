package api

import (
	"encoding/json"
	"flight-aggregator/internal/models"
	"flight-aggregator/internal/service"
	"log"
	"net/http"
	"strings"
)

// Handler handles HTTP requests
type Handler struct {
	searchService *service.SearchService
}

// NewHandler creates a new API handler
func NewHandler(searchService *service.SearchService) *Handler {
	return &Handler{
		searchService: searchService,
	}
}

// Search handles flight search requests
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	var req models.SearchRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		respondWithErrorDetailed(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Perform search
	response, err := h.searchService.Search(r.Context(), req)
	if err != nil {
		// Determine appropriate error code based on error type
		statusCode := http.StatusInternalServerError
		errorType := "Internal server error"

		errMsg := err.Error()

		// Check for validation errors (common patterns)
		if strings.Contains(errMsg, "invalid") ||
			strings.Contains(errMsg, "required") ||
			strings.Contains(errMsg, "must be") {
			statusCode = http.StatusBadRequest
			errorType = "Validation error"
		}

		// Check for timeout errors
		if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "context deadline exceeded") {
			statusCode = http.StatusGatewayTimeout
			errorType = "Request timeout"
		}

		log.Printf("Search failed: %v", err)
		respondWithErrorDetailed(w, statusCode, errorType, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

// Health checks if the service is healthy
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// ListProviders lists all available providers and their status
func (h *Handler) ListProviders(w http.ResponseWriter, r *http.Request) {
	providers := h.searchService.GetProviders()
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"providers": providers,
	})
}

// Helper functions
func respondWithErrorDetailed(w http.ResponseWriter, code int, errorType string, message string) {
	respondWithJSON(w, code, models.ErrorResponse{
		Error:   errorType,
		Message: message,
		Code:    code,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
