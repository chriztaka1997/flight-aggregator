package models

// ProviderResult represents the result from a single provider
type ProviderResult struct {
	Provider string
	Flights  []Flight
	Error    error
}

// AggregatedResults represents the combined results from all providers
type AggregatedResults struct {
	Results         []Flight
	ProviderResults map[string]int    // provider name -> count of flights
	ProviderErrors  map[string]string // provider name -> error message
}
