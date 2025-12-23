package providers

import (
	"encoding/json"
	"fmt"
	"os"
)

// GarudaResponse represents Garuda Indonesia API response structure
type GarudaResponse struct {
	Status  string         `json:"status"`
	Flights []GarudaFlight `json:"flights"`
}

type GarudaFlight struct {
	FlightID        string          `json:"flight_id"`
	Airline         string          `json:"airline"`
	AirlineCode     string          `json:"airline_code"`
	Departure       GarudaLocation  `json:"departure"`
	Arrival         GarudaLocation  `json:"arrival"`
	DurationMinutes int             `json:"duration_minutes"`
	Stops           int             `json:"stops"`
	Aircraft        string          `json:"aircraft"`
	Price           PriceInfo       `json:"price"`
	AvailableSeats  int             `json:"available_seats"`
	FareClass       string          `json:"fare_class"`
	Baggage         GarudaBaggage   `json:"baggage"`
	Amenities       []string        `json:"amenities,omitempty"`
	Segments        []GarudaSegment `json:"segments,omitempty"`
}

type GarudaLocation struct {
	Airport  string `json:"airport"`
	City     string `json:"city"`
	Time     string `json:"time"`
	Terminal string `json:"terminal,omitempty"`
}

type GarudaBaggage struct {
	CarryOn int `json:"carry_on"`
	Checked int `json:"checked"`
}

type GarudaSegment struct {
	FlightNumber    string         `json:"flight_number"`
	Departure       GarudaLocation `json:"departure"`
	Arrival         GarudaLocation `json:"arrival"`
	DurationMinutes int            `json:"duration_minutes"`
	LayoverMinutes  int            `json:"layover_minutes,omitempty"`
}

// LionAirResponse represents Lion Air API response structure
type LionAirResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AvailableFlights []LionAirFlight `json:"available_flights"`
	} `json:"data"`
}

type LionAirFlight struct {
	ID      string `json:"id"`
	Carrier struct {
		Name string `json:"name"`
		IATA string `json:"iata"`
	} `json:"carrier"`
	Route struct {
		From LionAirAirport `json:"from"`
		To   LionAirAirport `json:"to"`
	} `json:"route"`
	Schedule struct {
		Departure         string `json:"departure"`
		DepartureTimezone string `json:"departure_timezone"`
		Arrival           string `json:"arrival"`
		ArrivalTimezone   string `json:"arrival_timezone"`
	} `json:"schedule"`
	FlightTime int  `json:"flight_time"`
	IsDirect   bool `json:"is_direct"`
	StopCount  int  `json:"stop_count,omitempty"`
	Layovers   []struct {
		Airport         string `json:"airport"`
		DurationMinutes int    `json:"duration_minutes"`
	} `json:"layovers,omitempty"`
	Pricing struct {
		Total    float64 `json:"total"`
		Currency string  `json:"currency"`
		FareType string  `json:"fare_type"`
	} `json:"pricing"`
	SeatsLeft int    `json:"seats_left"`
	PlaneType string `json:"plane_type"`
	Services  struct {
		WifiAvailable bool `json:"wifi_available"`
		MealsIncluded bool `json:"meals_included"`
		Baggage       struct {
			Cabin string `json:"cabin"`
			Hold  string `json:"hold"`
		} `json:"baggage_allowance"`
	} `json:"services"`
}

type LionAirAirport struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

// BatikResponse represents Batik Air API response structure
type BatikResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Results []BatikFlight `json:"results"`
}

type BatikFlight struct {
	FlightNumber      string `json:"flightNumber"`
	AirlineName       string `json:"airlineName"`
	AirlineIATA       string `json:"airlineIATA"`
	Origin            string `json:"origin"`
	Destination       string `json:"destination"`
	DepartureDateTime string `json:"departureDateTime"`
	ArrivalDateTime   string `json:"arrivalDateTime"`
	TravelTime        string `json:"travelTime"`
	NumberOfStops     int    `json:"numberOfStops"`
	Fare              struct {
		BasePrice    float64 `json:"basePrice"`
		Taxes        float64 `json:"taxes"`
		TotalPrice   float64 `json:"totalPrice"`
		CurrencyCode string  `json:"currencyCode"`
		Class        string  `json:"class"`
	} `json:"fare"`
	SeatsAvailable  int      `json:"seatsAvailable"`
	AircraftModel   string   `json:"aircraftModel"`
	BaggageInfo     string   `json:"baggageInfo,omitempty"`
	OnboardServices []string `json:"onboardServices,omitempty"`
	Connections     []struct {
		Airport  string `json:"stopAirport"`
		WaitTime string `json:"stopDuration"`
	} `json:"connections,omitempty"`
}

// AirAsiaResponse represents AirAsia API response structure
type AirAsiaResponse struct {
	Status  string          `json:"status"`
	Flights []AirAsiaFlight `json:"flights"`
}

type AirAsiaFlight struct {
	FlightCode    string  `json:"flight_code"`
	Airline       string  `json:"airline"`
	FromAirport   string  `json:"from_airport"`
	ToAirport     string  `json:"to_airport"`
	DepartTime    string  `json:"depart_time"`
	ArriveTime    string  `json:"arrive_time"`
	DurationHours float64 `json:"duration_hours"`
	DirectFlight  bool    `json:"direct_flight"`
	Stops         []struct {
		Airport         string `json:"airport"`
		WaitTimeMinutes int    `json:"wait_time_minutes"`
	} `json:"stops,omitempty"`
	PriceIDR    float64 `json:"price_idr"`
	Seats       int     `json:"seats"`
	CabinClass  string  `json:"cabin_class"`
	BaggageNote string  `json:"baggage_note,omitempty"`
}

type PriceInfo struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// LoadMockData loads mock data from a JSON file
func LoadMockData(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from %s: %w", filePath, err)
	}

	return nil
}
