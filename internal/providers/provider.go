package providers

import (
	"context"
	"flight-aggregator/internal/models"
)

// Provider defines the interface that all airline providers must implement
type Provider interface {
	// Name returns the provider name
	Name() string

	// Search performs a flight search
	Search(ctx context.Context, req models.SearchRequest) ([]models.Flight, error)

	// HealthCheck returns true if the provider is healthy
	HealthCheck() bool
}
