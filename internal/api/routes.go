package api

import (
	"github.com/gorilla/mux"
)

// SetupRoutes configures all API routes
func SetupRoutes(h *Handler) *mux.Router {
	router := mux.NewRouter()

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Search endpoint
	api.HandleFunc("/search", h.Search).Methods("POST")

	// Health check endpoint
	api.HandleFunc("/health", h.Health).Methods("GET")

	// Provider status endpoint
	api.HandleFunc("/providers", h.ListProviders).Methods("GET")

	return router
}
