package providers

import (
	"context"
	"flight-aggregator/internal/models"
	"flight-aggregator/pkg/utils"
	"fmt"
	"time"
)

// GarudaProvider implements the Provider interface for Garuda Indonesia
type GarudaProvider struct {
	BaseProvider
}

// NewGarudaProviderFromConfig creates a new Garuda Indonesia provider from config
func NewGarudaProviderFromConfig(cfg ProviderConfig) *GarudaProvider {
	return &GarudaProvider{
		BaseProvider: NewBaseProviderFromConfig(cfg),
	}
}

// Search performs flight search for Garuda Indonesia
func (g *GarudaProvider) Search(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	// Simulate network delay
	if err := g.SimulateDelay(ctx); err != nil {
		return nil, fmt.Errorf("garuda: %w", ErrProviderTimeout)
	}

	// Simulate random failures
	if err := g.SimulateFailure(); err != nil {
		return nil, err
	}

	// Load mock data
	var response GarudaResponse
	if err := LoadMockData(g.mockDataPath, &response); err != nil {
		return nil, fmt.Errorf("garuda: %w: %v", ErrInvalidResponse, err)
	}

	if response.Status != "success" {
		return nil, fmt.Errorf("garuda: unsuccessful response status: %s", response.Status)
	}

	// Convert to unified Flight model and filter based on search criteria
	flights := make([]models.Flight, 0, len(response.Flights))
	for _, gf := range response.Flights {
		// Filter by origin and destination
		if gf.Departure.Airport != req.Origin || gf.Arrival.Airport != req.Destination {
			continue
		}

		flight, err := g.convertToFlight(gf)
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

// convertToFlight converts Garuda-specific flight to unified Flight model
func (g *GarudaProvider) convertToFlight(gf GarudaFlight) (models.Flight, error) {
	// Parse departure time with timezone
	departureTime, err := utils.ParseFlexibleTime(gf.Departure.Time)
	if err != nil {
		return models.Flight{}, fmt.Errorf("invalid departure time: %w", err)
	}

	// Load timezone for departure airport
	departureTZ, err := time.LoadLocation(utils.GetTimezone(gf.Departure.Airport))
	if err != nil {
		departureTZ = time.UTC
	}
	departureTime = time.Date(
		departureTime.Year(), departureTime.Month(), departureTime.Day(),
		departureTime.Hour(), departureTime.Minute(), departureTime.Second(),
		departureTime.Nanosecond(), departureTZ,
	)

	// Parse arrival time with timezone
	arrivalTime, err := utils.ParseFlexibleTime(gf.Arrival.Time)
	if err != nil {
		return models.Flight{}, fmt.Errorf("invalid arrival time: %w", err)
	}

	// Load timezone for arrival airport
	arrivalTZ, err := time.LoadLocation(utils.GetTimezone(gf.Arrival.Airport))
	if err != nil {
		arrivalTZ = time.UTC
	}
	arrivalTime = time.Date(
		arrivalTime.Year(), arrivalTime.Month(), arrivalTime.Day(),
		arrivalTime.Hour(), arrivalTime.Minute(), arrivalTime.Second(),
		arrivalTime.Nanosecond(), arrivalTZ,
	)

	// Extract airline code from flight ID
	airlineCode := utils.ExtractAirlineCode(gf.FlightID)

	// Format baggage info
	carryOnText := fmt.Sprintf("%d bag", gf.Baggage.CarryOn)
	if gf.Baggage.CarryOn > 1 || gf.Baggage.CarryOn == 0 {
		carryOnText = fmt.Sprintf("%d bags", gf.Baggage.CarryOn)
	}

	checkedText := fmt.Sprintf("%d bag", gf.Baggage.Checked)
	if gf.Baggage.Checked > 1 || gf.Baggage.Checked == 0 {
		checkedText = fmt.Sprintf("%d bags", gf.Baggage.Checked)
	}

	// Create flight object with unique ID
	flightID := fmt.Sprintf("%s_%s", gf.FlightID, g.Name())

	flight := models.Flight{
		ID:           flightID,
		Provider:     g.Name(),
		FlightNumber: gf.FlightID,
		Airline: models.Airline{
			Name: gf.Airline,
			Code: airlineCode,
		},
		Departure: models.FlightLocation{
			Airport:   gf.Departure.Airport,
			City:      utils.GetCityName(gf.Departure.Airport),
			Datetime:  departureTime,
			Timestamp: departureTime.Unix(),
		},
		Arrival: models.FlightLocation{
			Airport:   gf.Arrival.Airport,
			City:      utils.GetCityName(gf.Arrival.Airport),
			Datetime:  arrivalTime,
			Timestamp: arrivalTime.Unix(),
		},
		Duration: models.Duration{
			TotalMinutes: gf.DurationMinutes,
			Formatted:    utils.FormatDuration(gf.DurationMinutes),
		},
		Stops: gf.Stops,
		Price: models.Money{
			Amount:          gf.Price.Amount,
			Currency:        gf.Price.Currency,
			FormattedAmount: utils.FormatPrice(gf.Price.Amount, gf.Price.Currency),
			FormattedPrice:  utils.FormatPriceWithSymbol(gf.Price.Amount, gf.Price.Currency),
		},
		CabinClass:     gf.FareClass,
		Aircraft:       gf.Aircraft,
		AvailableSeats: gf.AvailableSeats,
		Amenities:      gf.Amenities,
		Baggage: models.BaggageInfo{
			CarryOn: carryOnText,
			Checked: checkedText,
		},
	}

	return flight, nil
}
