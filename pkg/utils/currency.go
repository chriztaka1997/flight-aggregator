package utils

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// CurrencyInfo holds information about a currency
type CurrencyInfo struct {
	Code   string
	Symbol string
	Name   string
}

// Common Indonesian currencies
var (
	IDR = CurrencyInfo{Code: "IDR", Symbol: "Rp", Name: "Indonesian Rupiah"}
	USD = CurrencyInfo{Code: "USD", Symbol: "$", Name: "US Dollar"}
	SGD = CurrencyInfo{Code: "SGD", Symbol: "S$", Name: "Singapore Dollar"}
	MYR = CurrencyInfo{Code: "MYR", Symbol: "RM", Name: "Malaysian Ringgit"}
	EUR = CurrencyInfo{Code: "EUR", Symbol: "€", Name: "Euro"}
	GBP = CurrencyInfo{Code: "GBP", Symbol: "£", Name: "British Pound"}
)

// CurrencySymbols maps currency codes to their symbols
var CurrencySymbols = map[string]string{
	"IDR": "Rp",
	"USD": "$",
	"SGD": "S$",
	"MYR": "RM",
	"EUR": "€",
	"GBP": "£",
	"JPY": "¥",
	"AUD": "A$",
	"CNY": "¥",
	"THB": "฿",
}

// GetCurrencySymbol returns the symbol for a currency code
func GetCurrencySymbol(code string) string {
	if symbol, ok := CurrencySymbols[code]; ok {
		return symbol
	}
	return code // Return code if symbol not found
}

// FormatCurrency formats an amount with the currency code
// Example: FormatCurrency(1000000, "IDR") -> "IDR 1,000,000"
func FormatCurrency(amount float64, currencyCode string) string {
	p := message.NewPrinter(language.English)
	return fmt.Sprintf("%s %s", currencyCode, p.Sprintf("%.0f", amount))
}

// FormatCurrencyWithSymbol formats an amount with the currency symbol
// Example: FormatCurrencyWithSymbol(1000000, "IDR") -> "Rp 1,000,000"
func FormatCurrencyWithSymbol(amount float64, currencyCode string) string {
	symbol := GetCurrencySymbol(currencyCode)
	p := message.NewPrinter(language.English)
	return fmt.Sprintf("%s %s", symbol, p.Sprintf("%.0f", amount))
}

// FormatCurrencyCompact formats an amount in a compact form
// Example: FormatCurrencyCompact(1500000, "IDR") -> "Rp 1.5M"
func FormatCurrencyCompact(amount float64, currencyCode string) string {
	symbol := GetCurrencySymbol(currencyCode)

	var value float64
	var unit string

	switch {
	case amount >= 1_000_000_000: // Billions
		value = amount / 1_000_000_000
		unit = "B"
	case amount >= 1_000_000: // Millions
		value = amount / 1_000_000
		unit = "M"
	case amount >= 1_000: // Thousands
		value = amount / 1_000
		unit = "K"
	default:
		return FormatCurrencyWithSymbol(amount, currencyCode)
	}

	// Format with appropriate decimal places, removing trailing zeros
	var formatted string
	if value >= 100 {
		formatted = fmt.Sprintf("%.0f", value)
	} else if value >= 10 {
		formatted = formatDecimal(value, 1)
	} else {
		formatted = formatDecimal(value, 2)
	}

	return fmt.Sprintf("%s %s%s", symbol, formatted, unit)
}

// formatDecimal formats a decimal number with max decimal places, removing trailing zeros
func formatDecimal(value float64, maxDecimals int) string {
	format := fmt.Sprintf("%%.%df", maxDecimals)
	str := fmt.Sprintf(format, value)

	// Remove trailing zeros
	for len(str) > 0 && str[len(str)-1] == '0' {
		str = str[:len(str)-1]
	}

	// Remove trailing decimal point if no decimals remain
	if len(str) > 0 && str[len(str)-1] == '.' {
		str = str[:len(str)-1]
	}

	return str
}

// FormatPrice formats a price for display in Indonesian locale
// Example: FormatPrice(1234567.89, "IDR") -> "Rp 1.234.567,89"
func FormatPrice(amount float64, currencyCode string) string {
	symbol := GetCurrencySymbol(currencyCode)

	// For IDR and similar currencies, no decimal places
	if currencyCode == "IDR" || currencyCode == "JPY" {
		return formatWithThousandsSeparator(amount, symbol, 0, ".", ",")
	}

	// For other currencies, 2 decimal places
	return formatWithThousandsSeparator(amount, symbol, 2, ",", ".")
}

// formatWithThousandsSeparator formats a number with custom separators
func formatWithThousandsSeparator(amount float64, symbol string, decimals int, thousandsSep, decimalSep string) string {
	// Extract integer and decimal parts
	intPart := int64(amount)
	decPart := amount - float64(intPart)

	// Format integer part with thousands separator
	intStr := fmt.Sprintf("%d", intPart)
	formatted := ""

	for i, digit := range intStr {
		if i > 0 && (len(intStr)-i)%3 == 0 {
			formatted += thousandsSep
		}
		formatted += string(digit)
	}

	// Add decimal part if needed
	if decimals > 0 {
		formatted += decimalSep
		decFormat := fmt.Sprintf("%%0%dd", decimals)
		formatted += fmt.Sprintf(decFormat, int(decPart*pow10(decimals)))
	}

	return fmt.Sprintf("%s %s", symbol, formatted)
}

// pow10 returns 10^n
func pow10(n int) float64 {
	result := 1.0
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

// ParsePrice parses a formatted price string back to a float
// Example: ParsePrice("Rp 1.234.567") -> 1234567.0
func ParsePrice(priceStr string) (float64, error) {
	// Remove common currency symbols and separators
	cleaned := priceStr

	// Remove currency symbols
	for _, symbol := range CurrencySymbols {
		cleaned = removeString(cleaned, symbol+" ")
		cleaned = removeString(cleaned, symbol)
	}

	// Remove common separators
	cleaned = removeString(cleaned, ".")
	cleaned = removeString(cleaned, ",")
	cleaned = removeString(cleaned, " ")

	// Parse as float
	var amount float64
	_, err := fmt.Sscanf(cleaned, "%f", &amount)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	return amount, nil
}

// removeString removes all occurrences of substr from s
func removeString(s, substr string) string {
	result := ""
	i := 0
	for i < len(s) {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			i += len(substr)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}

// ComparePrice compares two prices and returns:
// -1 if price1 < price2
//
//	0 if price1 == price2
//	1 if price1 > price2
func ComparePrice(price1, price2 float64) int {
	const epsilon = 0.01 // Tolerance for floating point comparison

	diff := price1 - price2
	if diff < -epsilon {
		return -1
	} else if diff > epsilon {
		return 1
	}
	return 0
}

// PriceRange represents a price range
type PriceRange struct {
	Min      float64
	Max      float64
	Currency string
}

// FormatPriceRange formats a price range for display
// Example: FormatPriceRange({500000, 1000000, "IDR"}) -> "Rp 500K - Rp 1M"
func FormatPriceRange(pr PriceRange) string {
	minStr := FormatCurrencyCompact(pr.Min, pr.Currency)
	maxStr := FormatCurrencyCompact(pr.Max, pr.Currency)
	return fmt.Sprintf("%s - %s", minStr, maxStr)
}

// Contains checks if an amount is within the price range
func (pr PriceRange) Contains(amount float64) bool {
	return amount >= pr.Min && amount <= pr.Max
}
