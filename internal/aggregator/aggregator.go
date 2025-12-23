package aggregator

import (
	"context"
	"flight-aggregator/internal/models"
	"flight-aggregator/internal/providers"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ProviderResult represents the result from a single provider
type ProviderResult struct {
	Provider string
	Flights  []models.Flight
	Error    error
	Duration time.Duration
}

// AggregatedResults contains all results from multiple providers
type AggregatedResults struct {
	Flights         []models.Flight
	ProviderResults map[string]int    // provider name -> number of flights
	ProviderErrors  map[string]string // provider name -> error message
	TotalDuration   time.Duration
}

// Aggregator handles parallel queries to multiple flight providers
type Aggregator struct {
	providers []providers.Provider
	timeout   time.Duration
}

// NewAggregator creates a new aggregator with the given providers and timeout
func NewAggregator(providerList []providers.Provider, timeout time.Duration) *Aggregator {
	return &Aggregator{
		providers: providerList,
		timeout:   timeout,
	}
}

// SearchAll queries all providers in parallel and aggregates results
func (a *Aggregator) SearchAll(ctx context.Context, req models.SearchRequest) (*AggregatedResults, error) {
	startTime := time.Now()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	// Filter providers based on airline filter if specified
	providersToQuery := a.providers
	if req.Filters != nil && len(req.Filters.Airlines) > 0 {
		// Create a map for faster lookup (case-insensitive)
		airlineFilter := make(map[string]bool)
		for _, airline := range req.Filters.Airlines {
			airlineFilter[strings.ToLower(airline)] = true
		}

		// Filter providers that match the airline filter
		filteredProviders := make([]providers.Provider, 0)
		for _, provider := range a.providers {
			if airlineFilter[strings.ToLower(provider.Name())] {
				filteredProviders = append(filteredProviders, provider)
			}
		}

		// Only update if we found matching providers
		if len(filteredProviders) > 0 {
			providersToQuery = filteredProviders
		}
	}

	// Create channels for communication
	results := make(chan ProviderResult, len(providersToQuery))
	var wg sync.WaitGroup

	// Fan-out: Launch goroutines for each provider
	for _, provider := range providersToQuery {
		wg.Add(1)
		go func(p providers.Provider) {
			defer wg.Done()
			a.queryProvider(ctx, p, req, results)
		}(provider)
	}

	// Close results channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Fan-in: Collect results
	aggregated := a.collectResults(results)
	aggregated.TotalDuration = time.Since(startTime)

	// Check if we got at least some results
	if len(aggregated.Flights) == 0 {
		return aggregated, fmt.Errorf("no flights found from any provider")
	}

	return aggregated, nil
}

// queryProvider queries a single provider and sends result to channel
func (a *Aggregator) queryProvider(ctx context.Context, provider providers.Provider, req models.SearchRequest, results chan<- ProviderResult) {
	providerStart := time.Now()

	// Execute search with context
	flights, err := provider.Search(ctx, req)

	// Send result to channel
	results <- ProviderResult{
		Provider: provider.Name(),
		Flights:  flights,
		Error:    err,
		Duration: time.Since(providerStart),
	}
}

// collectResults gathers all provider results from the channel
func (a *Aggregator) collectResults(results <-chan ProviderResult) *AggregatedResults {
	aggregated := &AggregatedResults{
		Flights:         make([]models.Flight, 0),
		ProviderResults: make(map[string]int),
		ProviderErrors:  make(map[string]string),
	}

	// Collect from channel until closed
	for result := range results {
		if result.Error != nil {
			// Track provider errors
			aggregated.ProviderErrors[result.Provider] = result.Error.Error()
		} else {
			// Add successful results
			aggregated.Flights = append(aggregated.Flights, result.Flights...)
			aggregated.ProviderResults[result.Provider] = len(result.Flights)
		}
	}

	return aggregated
}

// GetProviders returns the list of providers
func (a *Aggregator) GetProviders() []providers.Provider {
	return a.providers
}

// GetTimeout returns the configured timeout
func (a *Aggregator) GetTimeout() time.Duration {
	return a.timeout
}
