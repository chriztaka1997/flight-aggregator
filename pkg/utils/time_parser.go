package utils

import (
	"fmt"
	"time"
)

// ParseTime attempts to parse time in multiple formats with optional timezone
// If timezone is empty, assumes timezone is embedded in timeStr
// If timezone is provided, applies it to the parsed time
func ParseTime(timeStr string, timezone ...string) (time.Time, error) {
	// Determine if timezone was provided
	tz := ""
	if len(timezone) > 0 {
		tz = timezone[0]
	}

	// List of formats to try, ordered from most to least specific
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z0700", // Format with timezone without colon
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05-0700", // Format with timezone without colon
		"2006-01-02T15:04:05+07:00",
		"2006-01-02T15:04:05+0700", // Format with timezone without colon
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			// If timezone was provided and time has no timezone info, apply it
			if tz != "" && t.Location() == time.UTC && (format == "2006-01-02T15:04:05" || format == "2006-01-02 15:04:05") {
				loc, err := time.LoadLocation(tz)
				if err != nil {
					// Fallback to UTC if timezone is invalid
					return t, nil
				}
				return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc), nil
			}
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

// ParseFlexibleTime is a convenience wrapper for ParseTime without timezone
func ParseFlexibleTime(timeStr string) (time.Time, error) {
	return ParseTime(timeStr)
}

// ParseTimeWithTimezone is a convenience wrapper for ParseTime with timezone
func ParseTimeWithTimezone(timeStr, timezone string) (time.Time, error) {
	return ParseTime(timeStr, timezone)
}

// ParseTravelTime parses duration string like "1h 45m" or "2h 30m" to minutes
// Handles formats: "1h 45m", "2h", "45m"
func ParseTravelTime(travelTime string) int {
	var hours, minutes int

	// Try parsing formats like "1h 45m", "2h", "45m"
	fmt.Sscanf(travelTime, "%dh %dm", &hours, &minutes)
	if hours == 0 && minutes == 0 {
		fmt.Sscanf(travelTime, "%dh", &hours)
	}
	if hours == 0 && minutes == 0 {
		fmt.Sscanf(travelTime, "%dm", &minutes)
	}

	return hours*60 + minutes
}
