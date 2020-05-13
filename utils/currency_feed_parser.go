package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CurrencyFeed struct {
	Items []FeedItem
}

type FeedItem struct {
	Date       time.Time
	Currencies []Currency
}

type Currency struct {
	Symbol string
	Value  float64
}

// ParseItem parses currency feed item and appends results to the CurrencyFeed items list
func (f *CurrencyFeed) ParseItem(date string, content string) error {
	// Validate date length
	if len(date) != 31 {
		return fmt.Errorf("Invalid date length")
	}

	// Validate content length
	if len(content) < 7 {
		return fmt.Errorf("Invalid content length")
	}

	// Parse date with RFC1123Z format
	pubAt, err := time.Parse(time.RFC1123Z, date)
	if err != nil {
		return err
	}

	item := FeedItem{
		Date: pubAt,
	}

	// Trim trailing whitespace if exists
	trimmedContent := content
	if content[len(content)-1] == 32 {
		trimmedContent = content[:len(content)-1]
	}

	// Split content
	parts := strings.Split(trimmedContent, " ")
	length := len(parts)

	// Make sure parts data length is an odd number
	if (length & 1) != 0 {
		return fmt.Errorf("Invalid data length")
	}

	// Iterate over seperated currency bits
	for i := 0; i < (length - 1); i += 2 {
		value, err := strconv.ParseFloat(parts[i+1], 64)
		if err != nil {
			return err
		}

		item.Currencies = append(item.Currencies, Currency{parts[i], value})
	}

	// Append parsed item to the item list
	f.Items = append(f.Items, item)

	return nil
}
