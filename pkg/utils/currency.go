package utils

import (
	"fmt"
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

// FormatPrice formats a price for display with currency-specific formatting (without symbol)
// Example: FormatPrice(1234567.89, "IDR") -> "1.234.567,89"
func FormatPrice(amount float64, currencyCode string) string {
	var thousandsSep, decimalSep string
	var decimals int

	// Currency-specific formatting rules
	switch currencyCode {
	case "USD", "GBP", "SGD", "MYR":
		// English-style: comma for thousands, dot for decimal
		thousandsSep = ","
		decimalSep = "."
		decimals = 2
	case "EUR":
		// European style: dot for thousands, comma for decimal
		thousandsSep = "."
		decimalSep = ","
		decimals = 2
	case "IDR":
		// Indonesian Rupiah: dot for thousands, comma for decimal, 2 decimals
		thousandsSep = "."
		decimalSep = ","
		decimals = 2
	case "JPY", "KRW":
		// Japanese Yen/Korean Won: comma for thousands, no decimals
		thousandsSep = ","
		decimalSep = ""
		decimals = 0
	default:
		// Default: English-style with 2 decimals
		thousandsSep = ","
		decimalSep = "."
		decimals = 2
	}

	return formatWithThousandsSeparator(amount, decimals, thousandsSep, decimalSep)
}

// FormatPriceWithSymbol formats a price with currency symbol
// Example: FormatPriceWithSymbol(1234567.89, "IDR") -> "Rp 1.234.567,89"
func FormatPriceWithSymbol(amount float64, currencyCode string) string {
	symbol := GetCurrencySymbol(currencyCode)
	formatted := FormatPrice(amount, currencyCode)
	return fmt.Sprintf("%s %s", symbol, formatted)
}

// formatWithThousandsSeparator formats a number with custom separators
func formatWithThousandsSeparator(amount float64, decimals int, thousandsSep, decimalSep string) string {
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

	// Add decimal part (necessary for currencies with decimals)
	if decimals > 0 {
		formatted += decimalSep
		// Round the decimal part properly
		decFormat := fmt.Sprintf("%%0%dd", decimals)
		roundedDec := int(decPart*pow10(decimals) + 0.5)
		formatted += fmt.Sprintf(decFormat, roundedDec)
	}

	return fmt.Sprintf("%s", formatted)
}

// pow10 returns 10^n
func pow10(n int) float64 {
	result := 1.0
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}
