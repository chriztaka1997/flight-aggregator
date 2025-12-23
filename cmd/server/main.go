package main

import (
	"flight-aggregator/internal/api"
	"flight-aggregator/internal/service"
	"flight-aggregator/pkg/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Load configuration from .env.yaml
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Println("Configuration loaded successfully from .env.yaml")

	// Initialize search service with providers from config
	searchService := service.NewSearchServiceWithConfig(cfg)

	// Initialize API handler
	handler := api.NewHandler(searchService)

	// Setup routes
	router := api.SetupRoutes(handler)

	// Initialize rate limiter from config
	rateLimiter := api.NewRateLimiter(
		float64(cfg.RateLimit.Requests)/60.0, // Convert requests per minute to requests per second
		cfg.RateLimit.Requests,               // Burst
	)

	// Add middleware (order matters!)
	router.Use(api.RecoveryMiddleware)          // Recover from panics
	router.Use(api.LoggingMiddleware)           // Handle CORS
	router.Use(rateLimiter.RateLimitMiddleware) // Apply rate limiting

	// Configure server using config from .env.yaml
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.GetReadTimeout(),
		WriteTimeout: cfg.Server.GetWriteTimeout(),
		IdleTimeout:  cfg.Server.GetIdleTimeout(),
	}

	// Start server
	log.Printf("\nServer configuration:")
	log.Printf("  - Port: %d", cfg.Server.Port)
	log.Printf("  - Cache TTL: %s", cfg.Cache.TTL)
	log.Printf("  - Provider Timeout: %s", cfg.Provider.Timeout)
	log.Printf("  - Logging Level: %s", cfg.Logging.Level)
	log.Printf("  - Rate Limit: %d requests per %s", cfg.RateLimit.Requests, cfg.RateLimit.Window)
	log.Printf("\nStarting server on http://localhost:%d", cfg.Server.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
