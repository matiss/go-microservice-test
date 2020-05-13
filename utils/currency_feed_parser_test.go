package utils

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestParseItem(t *testing.T) {
	// Parse RSS2 test file
	rss2 := RSS2{}
	xmlContent, _ := ioutil.ReadFile("../ecb_rss.xml")
	err := rss2.Parse(xmlContent)
	if err != nil {
		t.Errorf("Could not parse XML file")
	}

	// Parse currency feed
	feed := CurrencyFeed{}

	// Iterate over feed items and parse currencies
	for _, item := range rss2.Items {
		err := feed.ParseItem(item.PubDate, string(item.Description))
		if err != nil {
			t.Error(err)
		}
	}

	// Validate feed item count
	if len(feed.Items) != 4 {
		t.Errorf("Invalid count of feed items, expected %d got: %d", 4, len(feed.Items))
		// No further tests can be performed if item count is not valid
		return
	}

	// Test first item of the feed
	first := feed.Items[0]

	// Get timezone (EEST)
	location, err := time.LoadLocation("Europe/Riga")
	if err != nil {
		t.Error(err)
		return
	}

	// Test date: "2020-05-06 03:00:00 +0300 EEST"
	testDate := time.Date(2020, 5, 6, 3, 0, 0, 0, location)

	// Test date
	if !testDate.Equal(first.Date) {
		t.Errorf("Invalid date, expected %s got: %s", testDate, first.Date)
	}

	// Make sure currency count is right
	if len(first.Currencies) != 32 {
		t.Errorf("Invalid count of currencies, expected %d got: %d", 32, len(first.Currencies))
		// No further tests can be performed if currency count is not valid
		return
	}

	// Validate first and last currencies
	firstCurrency := first.Currencies[0]
	lastCurrency := first.Currencies[len(first.Currencies)-1]

	if firstCurrency.Symbol != "AUD" {
		t.Errorf("Invalid symbol, expected %s got: %s", "AUD", firstCurrency.Symbol)
	}

	if firstCurrency.Value != 1.7046 {
		t.Errorf("Invalid value, expected %f got: %f", 1.7046, firstCurrency.Value)
	}

	if lastCurrency.Symbol != "ZAR" {
		t.Errorf("Invalid symbol, expected %s got: %s", "ZAR", firstCurrency.Symbol)
	}

	if lastCurrency.Value != 20.0603 {
		t.Errorf("Invalid value, expected %f got: %f", 20.0603, firstCurrency.Value)
	}

}
