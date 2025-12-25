package providers

import (
	"context"
	"flight-aggregator/internal/models"
	"flight-aggregator/pkg/utils"
	"fmt"
	"strings"
	"time"
)

// AirAsiaProvider implements the Provider interface for AirAsia
type AirAsiaProvider struct {
	BaseProvider
}

// NewAirAsiaProviderFromConfig creates a new AirAsia provider from config
func NewAirAsiaProviderFromConfig(cfg ProviderConfig) *AirAsiaProvider {
	return &AirAsiaProvider{
		BaseProvider: NewBaseProviderFromConfig(cfg),
	}
}

// Search performs flight search for AirAsia
func (a *AirAsiaProvider) Search(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	// Simulate network delay
	if err := a.SimulateDelay(ctx); err != nil {
		return nil, fmt.Errorf("airasia: %w", ErrProviderTimeout)
	}

	// Simulate random failures (higher rate for AirAsia)
	if err := a.SimulateFailure(); err != nil {
		return nil, err
	}

	// Load mock data
	var response AirAsiaResponse
	if err := LoadMockData(a.mockDataPath, &response); err != nil {
		return nil, fmt.Errorf("airasia: %w: %v", ErrInvalidResponse, err)
	}

	// Convert to unified Flight model and filter based on search criteria
	flights := make([]models.Flight, 0, len(response.Flights))
	for _, af := range response.Flights {
		// Filter by origin and destination
		if af.FromAirport != req.Origin || af.ToAirport != req.Destination {
			continue
		}

		flight, err := a.convertToFlight(af)
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

// convertToFlight converts AirAsia-specific flight to unified Flight model
func (a *AirAsiaProvider) convertToFlight(af AirAsiaFlight) (models.Flight, error) {
	// Parse departure time with timezone
	departureTime, err := utils.ParseFlexibleTime(af.DepartTime)
	if err != nil {
		return models.Flight{}, fmt.Errorf("invalid departure time: %w", err)
	}

	// Load timezone for departure airport
	departureTZ, err := time.LoadLocation(utils.GetTimezone(af.FromAirport))
	if err != nil {
		departureTZ = time.UTC
	}
	departureTime = time.Date(
		departureTime.Year(), departureTime.Month(), departureTime.Day(),
		departureTime.Hour(), departureTime.Minute(), departureTime.Second(),
		departureTime.Nanosecond(), departureTZ,
	)

	// Parse arrival time with timezone
	arrivalTime, err := utils.ParseFlexibleTime(af.ArriveTime)
	if err != nil {
		return models.Flight{}, fmt.Errorf("invalid arrival time: %w", err)
	}

	// Load timezone for arrival airport
	arrivalTZ, err := time.LoadLocation(utils.GetTimezone(af.ToAirport))
	if err != nil {
		arrivalTZ = time.UTC
	}
	arrivalTime = time.Date(
		arrivalTime.Year(), arrivalTime.Month(), arrivalTime.Day(),
		arrivalTime.Hour(), arrivalTime.Minute(), arrivalTime.Second(),
		arrivalTime.Nanosecond(), arrivalTZ,
	)

	// Convert duration from hours to minutes
	durationMinutes := int(af.DurationHours * 60)

	// Determine number of stops
	stops := 0
	if !af.DirectFlight {
		stops = len(af.Stops)
	}

	// Extract airline code from flight number
	airlineCode := utils.ExtractAirlineCode(af.FlightCode)

	// Parse baggage info from BaggageNote
	carryOn := ""
	checked := ""
	if af.BaggageNote != "" {
		parts := strings.Split(af.BaggageNote, ",")
		if len(parts) >= 1 {
			carryOn = strings.TrimSpace(parts[0])
		}
		if len(parts) >= 2 {
			checked = strings.TrimSpace(parts[1])
		}
	}

	// Create flight object with unique ID
	flightID := fmt.Sprintf("%s_%s", af.FlightCode, a.Name())

	flight := models.Flight{
		ID:           flightID,
		Provider:     a.Name(),
		FlightNumber: af.FlightCode,
		Airline: models.Airline{
			Name: af.Airline,
			Code: airlineCode,
		},
		Departure: models.FlightLocation{
			Airport:   af.FromAirport,
			City:      utils.GetCityName(af.FromAirport),
			Datetime:  departureTime,
			Timestamp: departureTime.Unix(),
		},
		Arrival: models.FlightLocation{
			Airport:   af.ToAirport,
			City:      utils.GetCityName(af.ToAirport),
			Datetime:  arrivalTime,
			Timestamp: arrivalTime.Unix(),
		},
		Duration: models.Duration{
			TotalMinutes: durationMinutes,
			Formatted:    utils.FormatDuration(durationMinutes),
		},
		Stops: stops,
		Price: models.Money{
			Amount:          af.PriceIDR,
			Currency:        "IDR",
			FormattedAmount: utils.FormatPrice(af.PriceIDR, "IDR"),
			FormattedPrice:  utils.FormatPriceWithSymbol(af.PriceIDR, "IDR"),
		},
		CabinClass:     af.CabinClass,
		AvailableSeats: af.Seats,
		Amenities:      []string{},
		Baggage: models.BaggageInfo{
			CarryOn: carryOn,
			Checked: checked,
		},
	}

	return flight, nil
}
