package filter

import (
	"flight-aggregator/internal/models"
)

// FilterEngine handles filtering of flight results
type FilterEngine struct{}

// NewFilterEngine creates a new filter engine
func NewFilterEngine() *FilterEngine {
	return &FilterEngine{}
}

// Apply applies all filters to the flight list
func (f *FilterEngine) Apply(flights []models.Flight, filters models.FilterOptions) []models.Flight {
	result := flights

	// Apply price filter
	if filters.MinPrice != nil || filters.MaxPrice != nil {
		result = f.filterByPrice(result, filters.MinPrice, filters.MaxPrice)
	}

	// Apply stops filter
	if filters.MaxStops != nil {
		result = f.filterByStops(result, *filters.MaxStops)
	}

	// Apply departure time filter
	if filters.DepartureTime != nil {
		result = f.filterByDepartureTime(result, filters.DepartureTime)
	}

	// Apply arrival time filter
	if filters.ArrivalTime != nil {
		result = f.filterByArrivalTime(result, filters.ArrivalTime)
	}

	// Apply duration filter
	if filters.MaxDuration != nil {
		result = f.filterByDuration(result, *filters.MaxDuration)
	}

	return result
}

// filterByPrice filters flights within price range
func (f *FilterEngine) filterByPrice(flights []models.Flight, minPrice, maxPrice *float64) []models.Flight {
	filtered := make([]models.Flight, 0)

	for _, flight := range flights {
		price := flight.Price.Amount

		// Check minimum price
		if minPrice != nil && price < *minPrice {
			continue
		}

		// Check maximum price
		if maxPrice != nil && price > *maxPrice {
			continue
		}

		filtered = append(filtered, flight)
	}

	return filtered
}

// filterByStops filters flights by maximum number of stops
func (f *FilterEngine) filterByStops(flights []models.Flight, maxStops int) []models.Flight {
	filtered := make([]models.Flight, 0)

	for _, flight := range flights {
		if flight.Stops <= maxStops {
			filtered = append(filtered, flight)
		}
	}

	return filtered
}

// filterByDepartureTime filters flights by departure time range
func (f *FilterEngine) filterByDepartureTime(flights []models.Flight, timeRange *models.TimeRange) []models.Flight {
	filtered := make([]models.Flight, 0)

	for _, flight := range flights {
		hour := flight.Departure.Datetime.Hour()

		// Check if within range (inclusive)
		if hour >= timeRange.Start && hour <= timeRange.End {
			filtered = append(filtered, flight)
		}
	}

	return filtered
}

// filterByArrivalTime filters flights by arrival time range
func (f *FilterEngine) filterByArrivalTime(flights []models.Flight, timeRange *models.TimeRange) []models.Flight {
	filtered := make([]models.Flight, 0)

	for _, flight := range flights {
		hour := flight.Arrival.Datetime.Hour()

		// Check if within range (inclusive)
		if hour >= timeRange.Start && hour <= timeRange.End {
			filtered = append(filtered, flight)
		}
	}

	return filtered
}

// filterByDuration filters flights by maximum duration in minutes
func (f *FilterEngine) filterByDuration(flights []models.Flight, maxDuration int) []models.Flight {
	filtered := make([]models.Flight, 0)

	for _, flight := range flights {
		if flight.Duration.TotalMinutes <= maxDuration {
			filtered = append(filtered, flight)
		}
	}

	return filtered
}
