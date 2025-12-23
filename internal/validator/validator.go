package validator

import (
	"flight-aggregator/internal/models"
	"fmt"
	"strings"
	"time"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validator handles data validation
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateSearchRequest validates a search request
func (v *Validator) ValidateSearchRequest(req models.SearchRequest) error {
	// Validate origin
	if err := v.validateAirportCode(req.Origin, "Origin"); err != nil {
		return err
	}

	// Validate destination
	if err := v.validateAirportCode(req.Destination, "Destination"); err != nil {
		return err
	}

	// Origin and destination must be different
	if strings.EqualFold(req.Origin, req.Destination) {
		return ValidationError{
			Field:   "Destination",
			Message: "origin and destination must be different",
		}
	}

	// Validate departure date
	departureDate, err := v.validateDate(req.DepartureDate, "DepartureDate")
	if err != nil {
		return err
	}

	// Validate return date if provided
	if req.ReturnDate != nil && *req.ReturnDate != "" {
		returnDate, err := v.validateDate(*req.ReturnDate, "ReturnDate")
		if err != nil {
			return err
		}

		// Return date must be on or after departure date
		if returnDate.Before(departureDate) {
			return ValidationError{
				Field:   "ReturnDate",
				Message: "return date must be on or after departure date",
			}
		}
	}

	// Validate passengers
	if req.Passengers < 1 {
		return ValidationError{Field: "Passengers", Message: "must have at least 1 passenger"}
	}

	if req.Passengers > 9 {
		return ValidationError{Field: "Passengers", Message: "maximum 9 passengers per search"}
	}

	// Validate cabin class
	validCabinClasses := map[string]bool{
		"economy":  true,
		"premium":  true,
		"business": true,
		"first":    true,
	}

	if !validCabinClasses[strings.ToLower(req.CabinClass)] {
		return ValidationError{
			Field:   "CabinClass",
			Message: "cabin class must be economy, premium, business, or first",
		}
	}

	// Validate filters if provided
	if req.Filters != nil {
		if err := v.ValidateFilters(*req.Filters); err != nil {
			return err
		}
	}

	// Validate return filters if provided
	if req.ReturnFilters != nil {
		if err := v.ValidateFilters(*req.ReturnFilters); err != nil {
			return err
		}
	}

	return nil
}

// validateAirportCode validates airport code format (IATA 3-letter code)
func (v *Validator) validateAirportCode(code, field string) error {
	if code == "" {
		return ValidationError{Field: field, Message: "airport code is required"}
	}

	if len(code) != 3 {
		return ValidationError{
			Field:   field,
			Message: "airport code must be 3 characters (IATA code)",
		}
	}

	// Check if all characters are letters
	for _, char := range code {
		if !((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')) {
			return ValidationError{
				Field:   field,
				Message: "airport code must contain only letters",
			}
		}
	}

	return nil
}

// ValidateFilters validates filter options
func (v *Validator) ValidateFilters(filters models.FilterOptions) error {
	// Validate price range
	if filters.MinPrice != nil && *filters.MinPrice < 0 {
		return ValidationError{Field: "MinPrice", Message: "minimum price cannot be negative"}
	}

	if filters.MaxPrice != nil && *filters.MaxPrice < 0 {
		return ValidationError{Field: "MaxPrice", Message: "maximum price cannot be negative"}
	}

	if filters.MinPrice != nil && filters.MaxPrice != nil {
		if *filters.MinPrice > *filters.MaxPrice {
			return ValidationError{
				Field:   "MaxPrice",
				Message: "maximum price must be greater than minimum price",
			}
		}
	}

	// Validate max stops
	if filters.MaxStops != nil && *filters.MaxStops < 0 {
		return ValidationError{Field: "MaxStops", Message: "maximum stops cannot be negative"}
	}

	// Validate time ranges
	if filters.DepartureTime != nil {
		if err := v.validateTimeRange(*filters.DepartureTime, "DepartureTime"); err != nil {
			return err
		}
	}

	if filters.ArrivalTime != nil {
		if err := v.validateTimeRange(*filters.ArrivalTime, "ArrivalTime"); err != nil {
			return err
		}
	}

	// Validate max duration
	if filters.MaxDuration != nil && *filters.MaxDuration <= 0 {
		return ValidationError{Field: "MaxDuration", Message: "maximum duration must be positive"}
	}

	return nil
}

// validateTimeRange validates a time range
func (v *Validator) validateTimeRange(timeRange models.TimeRange, field string) error {
	if timeRange.Start < 0 || timeRange.Start > 23 {
		return ValidationError{
			Field:   field + ".Start",
			Message: "start hour must be between 0 and 23",
		}
	}

	if timeRange.End < 0 || timeRange.End > 23 {
		return ValidationError{
			Field:   field + ".End",
			Message: "end hour must be between 0 and 23",
		}
	}

	if timeRange.Start > timeRange.End {
		return ValidationError{
			Field:   field,
			Message: "start hour must be less than or equal to end hour",
		}
	}

	return nil
}

// validateDate validates date format and returns parsed time
func (v *Validator) validateDate(dateStr, field string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, ValidationError{
			Field:   field,
			Message: "date is required",
		}
	}

	// Try common date formats
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"02-01-2006",
		"02/01/2006",
	}

	var parsedDate time.Time
	var lastErr error

	for _, format := range formats {
		date, err := time.Parse(format, dateStr)
		if err == nil {
			parsedDate = date
			lastErr = nil
			break
		}
		lastErr = err
	}

	if lastErr != nil {
		return time.Time{}, ValidationError{
			Field:   field,
			Message: "invalid date format (expected YYYY-MM-DD, YYYY/MM/DD, DD-MM-YYYY, or DD/MM/YYYY)",
		}
	}

	return parsedDate, nil
}
