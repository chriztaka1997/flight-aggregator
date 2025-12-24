package models

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// SearchRequest represents a flight search request
type SearchRequest struct {
	Origin          string         `json:"origin" validate:"required,len=3"`
	Destination     string         `json:"destination" validate:"required,len=3"`
	DepartureDate   string         `json:"departureDate" validate:"required"`
	ReturnDate      *string        `json:"returnDate,omitempty"`
	Passengers      int            `json:"passengers" validate:"min=1"`
	CabinClass      string         `json:"cabinClass" validate:"required"`
	Filters         *FilterOptions `json:"filters,omitempty"`
	SortBy          string         `json:"sortBy,omitempty"`
	SortOrder       string         `json:"sortOrder,omitempty"`
	ReturnFilters   *FilterOptions `json:"returnFilters,omitempty"`
	ReturnSortBy    string         `json:"returnSortBy,omitempty"`
	ReturnSortOrder string         `json:"returnSortOrder,omitempty"`
}

// FilterOptions represents filtering criteria for flights
type FilterOptions struct {
	MinPrice      *float64   `json:"minPrice,omitempty"`
	MaxPrice      *float64   `json:"maxPrice,omitempty"`
	MaxStops      *int       `json:"maxStops,omitempty"`
	Airlines      []string   `json:"airlines,omitempty"`
	DepartureTime *TimeRange `json:"departureTime,omitempty"`
	ArrivalTime   *TimeRange `json:"arrivalTime,omitempty"`
	MaxDuration   *int       `json:"maxDuration,omitempty"` // minutes
}

// TimeRange represents a time range filter (hours in 24-hour format)
type TimeRange struct {
	Start int `json:"start"` // Start hour (0-23)
	End   int `json:"end"`   // End hour (0-23)
}

// SearchResponse represents the search results
type SearchResponse struct {
	SearchCriteria        SearchCriteria `json:"search_criteria"`
	Metadata              SearchMetadata `json:"metadata"`
	Flights               []Flight       `json:"flights"`
	BestValueFlight       *Flight        `json:"best_value_flight,omitempty"`
	ReturnFlights         []Flight       `json:"return_flights,omitempty"`
	BestValueReturnFlight *Flight        `json:"best_value_return_flight,omitempty"`
}

// SearchCriteria represents the search parameters used for the query
type SearchCriteria struct {
	Origin        string  `json:"origin"`
	Destination   string  `json:"destination"`
	DepartureDate string  `json:"departure_date"`
	ReturnDate    *string `json:"return_date,omitempty"`
	Passengers    int     `json:"passengers"`
	CabinClass    string  `json:"cabin_class"`
}

// SearchMetadata contains metadata about the search operation
type SearchMetadata struct {
	TotalResults       int               `json:"total_results"`
	ProvidersQueried   int               `json:"providers_queried"`
	ProvidersSucceeded int               `json:"providers_succeeded"`
	ProvidersFailed    int               `json:"providers_failed"`
	SearchTimeMs       int               `json:"search_time_ms"`
	CacheHit           bool              `json:"cache_hit"`
	ProviderResults    map[string]int    `json:"provider_results,omitempty"`
	ProviderErrors     map[string]string `json:"provider_errors,omitempty"`
}
