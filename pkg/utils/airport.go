package utils

import "fmt"

// GetCityName returns the city name for Indonesian airports
func GetCityName(airportCode string) string {
	cityMap := map[string]string{
		"CGK": "Jakarta",
		"DPS": "Denpasar",
		"SUB": "Surabaya",
		"KNO": "Medan",
		"PLM": "Palembang",
		"PDG": "Padang",
		"BTJ": "Banda Aceh",
		"PKU": "Pekanbaru",
		"BDO": "Bandung",
		"SRG": "Semarang",
		"JOG": "Yogyakarta",
		"UPG": "Makassar",
		"BPN": "Balikpapan",
		"LOP": "Lombok",
		"BDJ": "Banjarmasin",
		"SOC": "Solo",
		"PKY": "Palangkaraya",
		"MDC": "Manado",
		"DJJ": "Jayapura",
		"AMQ": "Ambon",
		"TIM": "Timika",
		"SRR": "Sorong",
	}

	if city, ok := cityMap[airportCode]; ok {
		return city
	}
	return airportCode // Return airport code if city not found
}

// FormatDuration converts minutes to a formatted string (e.g., "4h 20m")
func FormatDuration(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60

	if hours > 0 && mins > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	} else if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dm", mins)
}

// ExtractAirlineCode extracts the airline code from flight number (e.g., "QZ7250" -> "QZ")
func ExtractAirlineCode(flightNumber string) string {
	if len(flightNumber) < 2 {
		return ""
	}
	// Airline codes are typically 2-3 letters at the start
	for i, r := range flightNumber {
		if r >= '0' && r <= '9' {
			return flightNumber[:i]
		}
	}
	return flightNumber
}
