package service

import (
	"context"
	"flight-aggregator/internal/aggregator"
	"flight-aggregator/internal/cache"
	"flight-aggregator/internal/filter"
	"flight-aggregator/internal/models"
	"flight-aggregator/internal/providers"
	"flight-aggregator/internal/ranking"
	"flight-aggregator/internal/validator"
	"flight-aggregator/pkg/config"
	"log"
	"time"
)

// SearchService handles flight search orchestration
type SearchService struct {
	providers  []providers.Provider
	aggregator *aggregator.Aggregator
	cache      *cache.Cache
	filter     *filter.FilterEngine
	sorter     *filter.Sorter
	scorer     *ranking.Scorer
	validator  *validator.Validator
}

// NewSearchServiceWithConfig creates a new search service with config-based providers
func NewSearchServiceWithConfig(cfg *config.Config) *SearchService {
	var providerList []providers.Provider

	// Initialize each provider if enabled
	if garudaCfg, exists := cfg.Provider.GetProviderConfig("garuda"); exists && garudaCfg.Enabled {
		log.Printf("Initializing provider: %s (delay: %v, failure rate: %.1f%%)",
			garudaCfg.Name, garudaCfg.GetResponseTime(), garudaCfg.FailureRate*100)
		providerList = append(providerList, providers.NewGarudaProviderFromConfig(providers.ProviderConfig{
			Name:         garudaCfg.Name,
			ResponseTime: garudaCfg.GetResponseTime(),
			FailureRate:  garudaCfg.FailureRate,
			DataPath:     garudaCfg.DataPath,
		}))
	}

	if lionairCfg, exists := cfg.Provider.GetProviderConfig("lionair"); exists && lionairCfg.Enabled {
		log.Printf("Initializing provider: %s (delay: %v, failure rate: %.1f%%)",
			lionairCfg.Name, lionairCfg.GetResponseTime(), lionairCfg.FailureRate*100)
		providerList = append(providerList, providers.NewLionAirProviderFromConfig(providers.ProviderConfig{
			Name:         lionairCfg.Name,
			ResponseTime: lionairCfg.GetResponseTime(),
			FailureRate:  lionairCfg.FailureRate,
			DataPath:     lionairCfg.DataPath,
		}))
	}

	if batikCfg, exists := cfg.Provider.GetProviderConfig("batik"); exists && batikCfg.Enabled {
		log.Printf("Initializing provider: %s (delay: %v, failure rate: %.1f%%)",
			batikCfg.Name, batikCfg.GetResponseTime(), batikCfg.FailureRate*100)
		providerList = append(providerList, providers.NewBatikProviderFromConfig(providers.ProviderConfig{
			Name:         batikCfg.Name,
			ResponseTime: batikCfg.GetResponseTime(),
			FailureRate:  batikCfg.FailureRate,
			DataPath:     batikCfg.DataPath,
		}))
	}

	if airAsiaCfg, exists := cfg.Provider.GetProviderConfig("airasia"); exists && airAsiaCfg.Enabled {
		log.Printf("Initializing provider: %s (delay: %v, failure rate: %.1f%%)",
			airAsiaCfg.Name, airAsiaCfg.GetResponseTime(), airAsiaCfg.FailureRate*100)
		providerList = append(providerList, providers.NewAirAsiaProviderFromConfig(providers.ProviderConfig{
			Name:         airAsiaCfg.Name,
			ResponseTime: airAsiaCfg.GetResponseTime(),
			FailureRate:  airAsiaCfg.FailureRate,
			DataPath:     airAsiaCfg.DataPath,
		}))
	}

	log.Printf("Initialized %d providers from configuration", len(providerList))

	// Initialize components
	aggregatorTimeout, _ := time.ParseDuration(cfg.Provider.Timeout)
	cacheTTL, _ := time.ParseDuration(cfg.Cache.TTL)

	return &SearchService{
		providers:  providerList,
		aggregator: aggregator.NewAggregator(providerList, aggregatorTimeout),
		cache:      cache.New(cacheTTL),
		filter:     filter.NewFilterEngine(),
		sorter:     filter.NewSorter(),
		scorer:     ranking.NewScorerFromConfig(cfg),
		validator:  validator.NewValidator(),
	}
}

// Search performs a flight search with full orchestration
func (s *SearchService) Search(ctx context.Context, req models.SearchRequest) (*models.SearchResponse, error) {
	startTime := time.Now()

	// Step 1: Validate request
	if err := s.validator.ValidateSearchRequest(req); err != nil {
		return nil, err
	}

	// Step 2: Check cache
	cacheKey := s.cache.GenerateKey(req)
	if cached, ok := s.cache.Get(cacheKey); ok {
		log.Printf("Cache hit for key: %s", cacheKey)
		response := cached.(*models.SearchResponse)
		// Mark as cache hit
		response.Metadata.CacheHit = true
		return response, nil
	}

	log.Printf("Cache miss for key: %s", cacheKey)

	// Step 3: Aggregate from providers
	aggregated, err := s.aggregator.SearchAll(ctx, req)
	if err != nil {
		// Return partial results if we have any
		if aggregated != nil && len(aggregated.Flights) > 0 {
			log.Printf("Partial results: got %d flights with errors", len(aggregated.Flights))
		} else {
			return nil, err
		}
	}

	flights := aggregated.Flights

	// Step 4: Apply filters if provided
	if req.Filters != nil {
		log.Printf("Applying filters to %d flights", len(flights))
		flights = s.filter.Apply(flights, *req.Filters)
		log.Printf("After filtering: %d flights remaining", len(flights))
	}

	// Step 5: Calculate scores and identify best value flight
	var bestValueFlight *models.Flight
	if len(flights) > 0 {
		log.Printf("Scoring %d flights", len(flights))
		scoredFlights := s.scorer.ScoreFlights(flights)
		// Extract the best value flight (highest score)
		if len(scoredFlights) > 0 {
			bestValueFlight = &scoredFlights[0].Flight
			log.Printf("Best value flight: %s with score %.2f", bestValueFlight.FlightNumber, scoredFlights[0].Score)
		}
		// Keep flights in original order (don't reorder)
	}

	// Step 6: Apply custom sorting if requested
	if req.SortBy != "" {
		log.Printf("Sorting flights by %s (%s)", req.SortBy, req.SortOrder)
		flights = s.sorter.Sort(flights, req.SortBy, req.SortOrder)
	}

	// Calculate providers succeeded and failed
	providersSucceeded := 0
	for _, count := range aggregated.ProviderResults {
		if count > 0 {
			providersSucceeded++
		}
	}
	providersFailed := len(aggregated.ProviderErrors)

	flightMetaData := models.SearchMetadata{
		TotalResults:       len(flights),
		ProvidersQueried:   len(aggregated.ProviderResults),
		ProvidersSucceeded: providersSucceeded,
		ProvidersFailed:    providersFailed,
		SearchTimeMs:       int(time.Since(startTime).Milliseconds()),
		CacheHit:           false,
		ProviderResults:    aggregated.ProviderResults,
		ProviderErrors:     aggregated.ProviderErrors,
	}

	// Step 6.5: Search for return flights if return date is provided
	var returnFlights []models.Flight
	var bestValueReturnFlight *models.Flight
	var returnMetadata *models.SearchMetadata
	if req.ReturnDate != nil && *req.ReturnDate != "" {
		log.Printf("Searching for return flights on %s", *req.ReturnDate)
		returnStartTime := time.Now()

		// Create a return flight search request (swap origin/destination)
		returnReq := models.SearchRequest{
			Origin:        req.Destination,
			Destination:   req.Origin,
			DepartureDate: *req.ReturnDate,
			Passengers:    req.Passengers,
			CabinClass:    req.CabinClass,
			Filters:       req.ReturnFilters,
			SortBy:        req.ReturnSortBy,
			SortOrder:     req.ReturnSortOrder,
		}

		// Check cache for return flights
		returnCacheKey := s.cache.GenerateKey(returnReq)
		if cached, ok := s.cache.Get(returnCacheKey); ok {
			log.Printf("Cache hit for return flights key: %s", returnCacheKey)
			cachedResponse := cached.(*models.SearchResponse)
			returnFlights = cachedResponse.Flights
			bestValueReturnFlight = cachedResponse.BestValueFlight
			returnMetadata = cachedResponse.ReturnMetadata
			if returnMetadata != nil {
				returnMetadata.CacheHit = true
				returnMetadata.SearchTimeMs = int(time.Since(returnStartTime).Milliseconds())
			}
			if bestValueReturnFlight != nil {
				log.Printf("Best value return flight from cache: %s", bestValueReturnFlight.FlightNumber)
			}
		} else {
			log.Printf("Cache miss for return flights key: %s", returnCacheKey)

			// Aggregate from providers for return flights
			returnAggregated, err := s.aggregator.SearchAll(ctx, returnReq)
			if err != nil {
				log.Printf("Error searching return flights: %v", err)
				if returnAggregated != nil && len(returnAggregated.Flights) > 0 {
					log.Printf("Partial return results: got %d flights with errors", len(returnAggregated.Flights))
				}
			}

			if returnAggregated != nil {
				returnFlights = returnAggregated.Flights

				// Apply filters if provided
				if req.ReturnFilters != nil {
					log.Printf("Applying filters to %d return flights", len(returnFlights))
					returnFlights = s.filter.Apply(returnFlights, *req.ReturnFilters)
					log.Printf("After filtering: %d return flights remaining", len(returnFlights))
				}

				// Calculate scores and identify best value return flight
				if len(returnFlights) > 0 {
					log.Printf("Scoring %d return flights", len(returnFlights))
					scoredReturnFlights := s.scorer.ScoreFlights(returnFlights)
					// Extract the best value return flight (highest score)
					if len(scoredReturnFlights) > 0 {
						bestValueReturnFlight = &scoredReturnFlights[0].Flight
						log.Printf("Best value return flight: %s with score %.2f", bestValueReturnFlight.FlightNumber, scoredReturnFlights[0].Score)
					}
					// Keep return flights in original order (don't reorder)
				}

				// Apply custom sorting if requested
				if req.ReturnSortBy != "" {
					log.Printf("Sorting return flights by %s (%s)", req.ReturnSortBy, req.ReturnSortOrder)
					returnFlights = s.sorter.Sort(returnFlights, req.ReturnSortBy, req.ReturnSortOrder)
				}

				// Build return metadata
				returnProvidersSucceeded := 0
				for _, count := range returnAggregated.ProviderResults {
					if count > 0 {
						returnProvidersSucceeded++
					}
				}
				returnProvidersFailed := len(returnAggregated.ProviderErrors)

				returnMetadata = &models.SearchMetadata{
					TotalResults:       len(returnFlights),
					ProvidersQueried:   len(returnAggregated.ProviderResults),
					ProvidersSucceeded: returnProvidersSucceeded,
					ProvidersFailed:    returnProvidersFailed,
					SearchTimeMs:       int(time.Since(returnStartTime).Milliseconds()),
					CacheHit:           false,
					ProviderResults:    returnAggregated.ProviderResults,
					ProviderErrors:     returnAggregated.ProviderErrors,
				}
			}
		}
	}

	// Build response
	response := &models.SearchResponse{
		SearchCriteria: models.SearchCriteria{
			Origin:        req.Origin,
			Destination:   req.Destination,
			DepartureDate: req.DepartureDate,
			ReturnDate:    req.ReturnDate,
			Passengers:    req.Passengers,
			CabinClass:    req.CabinClass,
		},
		Metadata:              flightMetaData,
		Flights:               flights,
		BestValueFlight:       bestValueFlight,
		ReturnFlights:         returnFlights,
		BestValueReturnFlight: bestValueReturnFlight,
		ReturnMetadata:        returnMetadata,
	}

	// Cache response
	s.cache.Set(cacheKey, response)
	log.Printf("Cached response for key: %s", cacheKey)

	return response, nil
}

// GetProviders returns list of available providers
func (s *SearchService) GetProviders() []string {
	providerNames := make([]string, len(s.providers))
	for i, p := range s.providers {
		providerNames[i] = p.Name()
	}
	return providerNames
}
