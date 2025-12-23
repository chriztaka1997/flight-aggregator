package models

import "time"

// Flight represents a unified flight model across all providers
type Flight struct {
	ID             string         `json:"id"`
	Provider       string         `json:"provider"`
	FlightNumber   string         `json:"flight_number"`
	Airline        Airline        `json:"airline"`
	Departure      FlightLocation `json:"departure"`
	Arrival        FlightLocation `json:"arrival"`
	Duration       Duration       `json:"duration"`
	Stops          int            `json:"stops"`
	Price          Money          `json:"price"`
	CabinClass     string         `json:"cabin_class"`
	AvailableSeats int            `json:"available_seats"`
	Aircraft       string         `json:"aircraft"`
	Amenities      []string       `json:"amenities"`
	Baggage        BaggageInfo    `json:"baggage"`
}

// Airline represents airline information
type Airline struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// FlightLocation represents departure or arrival information
type FlightLocation struct {
	Airport   string    `json:"airport"`
	City      string    `json:"city"`
	Datetime  time.Time `json:"datetime"`
	Timestamp int64     `json:"timestamp"`
}

// Duration represents flight duration information
type Duration struct {
	TotalMinutes int    `json:"total_minutes"`
	Formatted    string `json:"formatted"`
}

// Money represents monetary value with currency
type Money struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// BaggageInfo represents baggage allowance details
type BaggageInfo struct {
	CarryOn string `json:"carry_on"`
	Checked string `json:"checked"`
}
