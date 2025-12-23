package providers

import (
	"context"
	"flight-aggregator/internal/models"
	"flight-aggregator/pkg/utils"
	"fmt"
	"time"
)

// LionAirProvider implements the Provider interface for Lion Air
type LionAirProvider struct {
	BaseProvider
}

// NewLionAirProviderFromConfig creates a new Lion Air provider from config
func NewLionAirProviderFromConfig(cfg ProviderConfig) *LionAirProvider {
	return &LionAirProvider{
		BaseProvider: NewBaseProviderFromConfig(cfg),
	}
}

// Search performs flight search for Lion Air
func (l *LionAirProvider) Search(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	// Simulate network delay
	if err := l.SimulateDelay(ctx); err != nil {
		return nil, fmt.Errorf("lionair: %w", ErrProviderTimeout)
	}

	// Simulate random failures
	if err := l.SimulateFailure(); err != nil {
		return nil, err
	}

	// Load mock data
	var response LionAirResponse
	if err := LoadMockData(l.mockDataPath, &response); err != nil {
		return nil, fmt.Errorf("lionair: %w: %v", ErrInvalidResponse, err)
	}

	if !response.Success {
		return nil, fmt.Errorf("lionair: unsuccessful response")
	}

	// Convert to unified Flight model and filter based on search criteria
	flights := make([]models.Flight, 0, len(response.Data.AvailableFlights))
	for _, lf := range response.Data.AvailableFlights {
		// Filter by origin and destination
		if lf.Route.From.Code != req.Origin || lf.Route.To.Code != req.Destination {
			continue
		}

		flight, err := l.convertToFlight(lf)
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

// convertToFlight converts Lion Air-specific flight to unified Flight model
func (l *LionAirProvider) convertToFlight(lf LionAirFlight) (models.Flight, error) {
	// Parse departure time with timezone
	departureTime, err := utils.ParseTimeWithTimezone(lf.Schedule.Departure, lf.Schedule.DepartureTimezone)
	if err != nil {
		return models.Flight{}, fmt.Errorf("invalid departure time: %w", err)
	}

	// Load timezone for departure airport
	departureTZ, err := time.LoadLocation(lf.Schedule.DepartureTimezone)
	if err != nil {
		departureTZ = time.UTC
	}
	departureTime = time.Date(
		departureTime.Year(), departureTime.Month(), departureTime.Day(),
		departureTime.Hour(), departureTime.Minute(), departureTime.Second(),
		departureTime.Nanosecond(), departureTZ,
	)

	// Parse arrival time with timezone
	arrivalTime, err := utils.ParseTimeWithTimezone(lf.Schedule.Arrival, lf.Schedule.ArrivalTimezone)
	if err != nil {
		return models.Flight{}, fmt.Errorf("invalid arrival time: %w", err)
	}

	// Load timezone for arrival airport
	arrivalTZ, err := time.LoadLocation(lf.Schedule.ArrivalTimezone)
	if err != nil {
		arrivalTZ = time.UTC
	}
	arrivalTime = time.Date(
		arrivalTime.Year(), arrivalTime.Month(), arrivalTime.Day(),
		arrivalTime.Hour(), arrivalTime.Minute(), arrivalTime.Second(),
		arrivalTime.Nanosecond(), arrivalTZ,
	)

	// Determine number of stops
	stops := 0
	if !lf.IsDirect {
		stops = lf.StopCount
	}

	// Extract airline code from flight ID
	airlineCode := utils.ExtractAirlineCode(lf.ID)

	// Create flight object with unique ID
	flightID := fmt.Sprintf("%s_%s", lf.ID, l.Name())

	// Convert duration from minutes to Duration struct
	durationMinutes := lf.FlightTime

	// Build amenities list
	amenities := make([]string, 0)
	if lf.Services.WifiAvailable {
		amenities = append(amenities, "WiFi")
	}
	if lf.Services.MealsIncluded {
		amenities = append(amenities, "Meals")
	}

	flight := models.Flight{
		ID:           flightID,
		Provider:     l.Name(),
		FlightNumber: lf.ID,
		Airline: models.Airline{
			Name: lf.Carrier.Name,
			Code: airlineCode,
		},
		Departure: models.FlightLocation{
			Airport:   lf.Route.From.Code,
			City:      utils.GetCityName(lf.Route.From.Code),
			Datetime:  departureTime,
			Timestamp: departureTime.Unix(),
		},
		Arrival: models.FlightLocation{
			Airport:   lf.Route.To.Code,
			City:      utils.GetCityName(lf.Route.To.Code),
			Datetime:  arrivalTime,
			Timestamp: arrivalTime.Unix(),
		},
		Duration: models.Duration{
			TotalMinutes: durationMinutes,
			Formatted:    utils.FormatDuration(durationMinutes),
		},
		Stops:          stops,
		Price:          models.Money{Amount: lf.Pricing.Total, Currency: lf.Pricing.Currency},
		CabinClass:     lf.Pricing.FareType,
		Aircraft:       lf.PlaneType,
		AvailableSeats: lf.SeatsLeft,
		Amenities:      amenities,
		Baggage: models.BaggageInfo{
			CarryOn: lf.Services.Baggage.Cabin,
			Checked: lf.Services.Baggage.Hold,
		},
	}

	return flight, nil
}
