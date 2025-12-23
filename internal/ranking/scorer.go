package ranking

import (
	"flight-aggregator/internal/models"
	"flight-aggregator/pkg/config"
	"math"
	"sort"
)

// Weights represents the importance of different factors in scoring
type Weights struct {
	Price         float64 // Weight for price (lower is better)
	Duration      float64 // Weight for duration (shorter is better)
	Stops         float64 // Weight for number of stops (fewer is better)
	DepartureTime float64 // Weight for departure time preference
}

// Scorer calculates best value scores for flights
type Scorer struct {
	weights Weights
}

// NewScorerWithWeights creates a scorer with custom weights
func NewScorerWithWeights(weights Weights) *Scorer {
	return &Scorer{
		weights: weights,
	}
}

// NewScorerFromConfig creates a scorer using weights from configuration
func NewScorerFromConfig(cfg *config.Config) *Scorer {
	return &Scorer{
		weights: Weights{
			Price:         cfg.Scoring.Weights.Price,
			Duration:      cfg.Scoring.Weights.Duration,
			Stops:         cfg.Scoring.Weights.Stops,
			DepartureTime: cfg.Scoring.Weights.DepartureTime,
		},
	}
}

// FlightScore represents a flight with its calculated score
type FlightScore struct {
	Flight    models.Flight
	Score     float64
	Breakdown ScoreBreakdown
}

// ScoreBreakdown shows how the score was calculated
type ScoreBreakdown struct {
	PriceScore         float64
	DurationScore      float64
	StopsScore         float64
	DepartureTimeScore float64
}

// ScoreFlights calculates scores for all flights and returns them sorted by best score
func (s *Scorer) ScoreFlights(flights []models.Flight) []FlightScore {
	if len(flights) == 0 {
		return []FlightScore{}
	}

	scored := make([]FlightScore, len(flights))

	// Calculate min/max values for normalization
	minPrice, maxPrice := s.findPriceRange(flights)
	minDuration, maxDuration := s.findDurationRange(flights)

	for i, flight := range flights {
		breakdown := ScoreBreakdown{
			PriceScore:         s.scorePriceNormalized(flight.Price.Amount, minPrice, maxPrice),
			DurationScore:      s.scoreDurationNormalized(flight.Duration.TotalMinutes, minDuration, maxDuration),
			StopsScore:         s.scoreStops(flight.Stops),
			DepartureTimeScore: s.scoreDepartureTime(flight.Departure.Datetime.Hour()),
		}

		// Calculate weighted total score (0-100)
		totalScore := (breakdown.PriceScore*s.weights.Price +
			breakdown.DurationScore*s.weights.Duration +
			breakdown.StopsScore*s.weights.Stops +
			breakdown.DepartureTimeScore*s.weights.DepartureTime) * 100

		scored[i] = FlightScore{
			Flight:    flight,
			Score:     totalScore,
			Breakdown: breakdown,
		}
	}

	// Sort by score (highest first)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	return scored
}

// findPriceRange finds min and max prices
func (s *Scorer) findPriceRange(flights []models.Flight) (float64, float64) {
	if len(flights) == 0 {
		return 0, 0
	}

	minPrice := flights[0].Price.Amount
	maxPrice := flights[0].Price.Amount

	for _, flight := range flights {
		if flight.Price.Amount < minPrice {
			minPrice = flight.Price.Amount
		}
		if flight.Price.Amount > maxPrice {
			maxPrice = flight.Price.Amount
		}
	}

	return minPrice, maxPrice
}

// findDurationRange finds min and max durations
func (s *Scorer) findDurationRange(flights []models.Flight) (int, int) {
	if len(flights) == 0 {
		return 0, 0
	}

	minDuration := flights[0].Duration.TotalMinutes
	maxDuration := flights[0].Duration.TotalMinutes

	for _, flight := range flights {
		if flight.Duration.TotalMinutes < minDuration {
			minDuration = flight.Duration.TotalMinutes
		}
		if flight.Duration.TotalMinutes > maxDuration {
			maxDuration = flight.Duration.TotalMinutes
		}
	}

	return minDuration, maxDuration
}

// scorePriceNormalized scores price on 0-1 scale (lower price = higher score)
func (s *Scorer) scorePriceNormalized(price, minPrice, maxPrice float64) float64 {
	if maxPrice == minPrice {
		return 1.0
	}

	// Invert: lower price gets higher score
	normalized := 1.0 - ((price - minPrice) / (maxPrice - minPrice))
	return math.Max(0, math.Min(1, normalized))
}

// scoreDurationNormalized scores duration on 0-1 scale (shorter = higher score)
func (s *Scorer) scoreDurationNormalized(duration, minDuration, maxDuration int) float64 {
	if maxDuration == minDuration {
		return 1.0
	}

	// Invert: shorter duration gets higher score
	normalized := 1.0 - (float64(duration-minDuration) / float64(maxDuration-minDuration))
	return math.Max(0, math.Min(1, normalized))
}

// scoreStops scores based on number of stops (0-1 scale)
func (s *Scorer) scoreStops(stops int) float64 {
	// Direct flights get 1.0, each stop reduces score
	switch stops {
	case 0:
		return 1.0
	case 1:
		return 0.7
	case 2:
		return 0.4
	default:
		return 0.2
	}
}

// scoreDepartureTime scores departure time based on preference (0-1 scale)
// Prefers reasonable departure times (not too early, not too late)
func (s *Scorer) scoreDepartureTime(hour int) float64 {
	// Preferred times: 7am-10pm (19 hours)
	// Peak preferred: 8am-8pm (12 hours)

	switch {
	case hour >= 8 && hour <= 20: // 8am-8pm
		return 1.0
	case hour >= 6 && hour < 8: // 6am-8am
		return 0.8
	case hour > 20 && hour <= 22: // 8pm-10pm
		return 0.8
	case hour >= 5 && hour < 6: // 5am-6am
		return 0.6
	case hour > 22 && hour <= 23: // 10pm-11pm
		return 0.6
	default: // Midnight-5am (very early/late)
		return 0.3
	}
}

// GetTopFlights returns the top N best value flights
func (s *Scorer) GetTopFlights(flights []models.Flight, n int) []FlightScore {
	scored := s.ScoreFlights(flights)

	if len(scored) <= n {
		return scored
	}

	return scored[:n]
}
