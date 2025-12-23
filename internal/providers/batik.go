package providers

import (
	"context"
	"flight-aggregator/internal/models"
	"flight-aggregator/pkg/utils"
	"fmt"
	"strings"
	"time"
)

// BatikProvider implements the Provider interface for Batik Air
type BatikProvider struct {
	BaseProvider
}

// NewBatikProviderFromConfig creates a new Batik Air provider from config
func NewBatikProviderFromConfig(cfg ProviderConfig) *BatikProvider {
	return &BatikProvider{
		BaseProvider: NewBaseProviderFromConfig(cfg),
	}
}

// Search performs flight search for Batik Air
func (b *BatikProvider) Search(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	// Simulate network delay
	if err := b.SimulateDelay(ctx); err != nil {
		return nil, fmt.Errorf("batik: %w", ErrProviderTimeout)
	}

	// Simulate random failures
	if err := b.SimulateFailure(); err != nil {
		return nil, err
	}

	// Load mock data
	var response BatikResponse
	if err := LoadMockData(b.mockDataPath, &response); err != nil {
		return nil, fmt.Errorf("batik: %w: %v", ErrInvalidResponse, err)
	}

	if response.Code != 200 {
		return nil, fmt.Errorf("batik: unsuccessful response code: %d", response.Code)
	}

	// Convert to unified Flight model and filter based on search criteria
	flights := make([]models.Flight, 0, len(response.Results))
	for _, bf := range response.Results {
		// Filter by origin and destination
		if bf.Origin != req.Origin || bf.Destination != req.Destination {
			continue
		}

		flight, err := b.convertToFlight(bf)
		if err != nil {
			// Skip invalid flights but continue processing
			continue
		}

		// Filter by departure date (compare date only, ignoring time)
		flightDate := flight.Departure.Datetime.Format("2006-01-02")
		if flightDate != req.DepartureDate {
			continue
		}

		flights = append(flights, flight)
	}

	if len(flights) == 0 {
		return nil, ErrNoFlightsFound
	}

	return flights, nil
}

// convertToFlight converts Batik Air-specific flight to unified Flight model
func (b *BatikProvider) convertToFlight(bf BatikFlight) (models.Flight, error) {
	// Parse departure time with timezone
	departureTime, err := utils.ParseFlexibleTime(bf.DepartureDateTime)
	if err != nil {
		return models.Flight{}, fmt.Errorf("invalid departure time: %w", err)
	}

	// Load timezone for departure airport
	departureTZ, err := time.LoadLocation(utils.GetTimezone(bf.Origin))
	if err != nil {
		departureTZ = time.UTC
	}
	departureTime = time.Date(
		departureTime.Year(), departureTime.Month(), departureTime.Day(),
		departureTime.Hour(), departureTime.Minute(), departureTime.Second(),
		departureTime.Nanosecond(), departureTZ,
	)

	// Parse arrival time with timezone
	arrivalTime, err := utils.ParseFlexibleTime(bf.ArrivalDateTime)
	if err != nil {
		return models.Flight{}, fmt.Errorf("invalid arrival time: %w", err)
	}

	// Load timezone for arrival airport
	arrivalTZ, err := time.LoadLocation(utils.GetTimezone(bf.Destination))
	if err != nil {
		arrivalTZ = time.UTC
	}
	arrivalTime = time.Date(
		arrivalTime.Year(), arrivalTime.Month(), arrivalTime.Day(),
		arrivalTime.Hour(), arrivalTime.Minute(), arrivalTime.Second(),
		arrivalTime.Nanosecond(), arrivalTZ,
	)

	// Calculate duration in minutes from travel time string (e.g., "1h 45m")
	durationMinutes := utils.ParseTravelTime(bf.TravelTime)

	// Extract airline code from flight number
	airlineCode := utils.ExtractAirlineCode(bf.FlightNumber)

	// Parse baggage info from BaggageInfo
	carryOn := ""
	checked := ""
	if bf.BaggageInfo != "" {
		parts := strings.Split(bf.BaggageInfo, ",")
		if len(parts) >= 1 {
			carryOn = strings.TrimSpace(parts[0])
		}
		if len(parts) >= 2 {
			checked = strings.TrimSpace(parts[1])
		} else if len(parts) == 1 {
			// If only one part, put it in checked
			checked = strings.TrimSpace(parts[0])
			carryOn = "Standard baggage"
		}
	}

	// Create flight object with unique ID
	flightID := fmt.Sprintf("%s_%s", bf.FlightNumber, b.Name())

	// Build amenities list
	amenities := make([]string, 0)
	amenities = append(amenities, bf.OnboardServices...)

	flight := models.Flight{
		ID:           flightID,
		Provider:     b.Name(),
		FlightNumber: bf.FlightNumber,
		Airline: models.Airline{
			Name: bf.AirlineName,
			Code: airlineCode,
		},
		Departure: models.FlightLocation{
			Airport:   bf.Origin,
			City:      utils.GetCityName(bf.Origin),
			Datetime:  departureTime,
			Timestamp: departureTime.Unix(),
		},
		Arrival: models.FlightLocation{
			Airport:   bf.Destination,
			City:      utils.GetCityName(bf.Destination),
			Datetime:  arrivalTime,
			Timestamp: arrivalTime.Unix(),
		},
		Duration: models.Duration{
			TotalMinutes: durationMinutes,
			Formatted:    utils.FormatDuration(durationMinutes),
		},
		Stops:          bf.NumberOfStops,
		Price:          models.Money{Amount: bf.Fare.TotalPrice, Currency: bf.Fare.CurrencyCode},
		CabinClass:     bf.Fare.Class,
		Aircraft:       bf.AircraftModel,
		AvailableSeats: bf.SeatsAvailable,
		Amenities:      amenities,
		Baggage: models.BaggageInfo{
			CarryOn: carryOn,
			Checked: checked,
		},
	}

	return flight, nil
}
