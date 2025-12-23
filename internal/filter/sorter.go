package filter

import (
	"flight-aggregator/internal/models"
	"sort"
)

// Sorter handles sorting of flight results
type Sorter struct{}

// NewSorter creates a new sorter
func NewSorter() *Sorter {
	return &Sorter{}
}

// Sort sorts flights by the specified field and order
// sortBy options: "price", "duration", "departure", "arrival", "stops"
// sortOrder options: "asc" (ascending), "desc" (descending)
func (s *Sorter) Sort(flights []models.Flight, sortBy, sortOrder string) []models.Flight {
	// Create a copy to avoid modifying original slice
	result := make([]models.Flight, len(flights))
	copy(result, flights)

	// Determine sort order multiplier
	ascending := sortOrder != "desc"

	switch sortBy {
	case "price":
		s.sortByPrice(result, ascending)
	case "duration":
		s.sortByDuration(result, ascending)
	case "departure":
		s.sortByDepartureTime(result, ascending)
	case "arrival":
		s.sortByArrivalTime(result, ascending)
	case "stops":
		s.sortByStops(result, ascending)
	default:
		// Default: sort by price ascending
		s.sortByPrice(result, true)
	}

	return result
}

// sortByPrice sorts flights by price
func (s *Sorter) sortByPrice(flights []models.Flight, ascending bool) {
	sort.Slice(flights, func(i, j int) bool {
		if ascending {
			return flights[i].Price.Amount < flights[j].Price.Amount
		}
		return flights[i].Price.Amount > flights[j].Price.Amount
	})
}

// sortByDuration sorts flights by total duration
func (s *Sorter) sortByDuration(flights []models.Flight, ascending bool) {
	sort.Slice(flights, func(i, j int) bool {
		if ascending {
			return flights[i].Duration.TotalMinutes < flights[j].Duration.TotalMinutes
		}
		return flights[i].Duration.TotalMinutes > flights[j].Duration.TotalMinutes
	})
}

// sortByDepartureTime sorts flights by departure time
func (s *Sorter) sortByDepartureTime(flights []models.Flight, ascending bool) {
	sort.Slice(flights, func(i, j int) bool {
		if ascending {
			return flights[i].Departure.Datetime.Before(flights[j].Departure.Datetime)
		}
		return flights[i].Departure.Datetime.After(flights[j].Departure.Datetime)
	})
}

// sortByArrivalTime sorts flights by arrival time
func (s *Sorter) sortByArrivalTime(flights []models.Flight, ascending bool) {
	sort.Slice(flights, func(i, j int) bool {
		if ascending {
			return flights[i].Arrival.Datetime.Before(flights[j].Arrival.Datetime)
		}
		return flights[i].Arrival.Datetime.After(flights[j].Arrival.Datetime)
	})
}

// sortByStops sorts flights by number of stops
func (s *Sorter) sortByStops(flights []models.Flight, ascending bool) {
	sort.Slice(flights, func(i, j int) bool {
		if ascending {
			return flights[i].Stops < flights[j].Stops
		}
		return flights[i].Stops > flights[j].Stops
	})
}
