package providers

import "errors"

// Common provider errors
var (
	ErrProviderTimeout     = errors.New("provider request timeout")
	ErrProviderUnavailable = errors.New("provider unavailable")
	ErrInvalidResponse     = errors.New("invalid response from provider")
	ErrNoFlightsFound      = errors.New("no flights found")
)
