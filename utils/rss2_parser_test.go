package utils

import (
	"io/ioutil"
	"testing"
)

func TestParse(t *testing.T) {
	// Parse RSS2 test file
	rss2 := RSS2{}
	xmlContent, _ := ioutil.ReadFile("../ecb_rss.xml")
	err := rss2.Parse(xmlContent)
	if err != nil {
		t.Errorf("Could not parse XML file")
	}

	// Make sure title is parsed correctly
	if rss2.Title != "Valūtu kursi" {
		t.Errorf("Failed to parse title, expected %s got: %s", "Valūtu kursi", rss2.Title)
	}

	// Make sure description is parsed correctly
	if rss2.Description != "Eiropas Centrālās bankas publicētie eiro atsauces kursi" {
		t.Errorf("Failed to parse description, expected %s got: %s", "Eiropas Centrālās bankas publicētie eiro atsauces kursi", rss2.Description)
	}

	// Make sure lastBuildDate is parsed correctly
	if rss2.BuiltAt != "Mon, 11 May 2020 18:42:02 +0300" {
		t.Errorf("Failed to parse lastBuildDate, expected %s got: %s", "Mon, 11 May 2020 18:42:02 +0300", rss2.BuiltAt)
	}

	// Make sure ttl is parsed correctly
	if rss2.TTL != 5 {
		t.Errorf("Failed to parse ttl, expected %d got: %d", 5, rss2.TTL)
	}

	// Validate item count
	if len(rss2.Items) != 4 {
		t.Errorf("Invalid item count, expected %d got: %d", 4, len(rss2.Items))
	}

	// Validate first item
	item := rss2.Items[0]

	// Make sure item title is parsed correctly
	if item.Title != "Eiropas Centrālās bankas publicētie eiro atsauces kursi. 6. May_LONG" {
		t.Errorf("Failed to parse item title, expected %s got: %s", "Eiropas Centrālās bankas publicētie eiro atsauces kursi. 6. May_LONG", item.Title)
	}

	// Make sure item link is parsed correctly
	if item.Link != "https://www.bank.lv/" {
		t.Errorf("Failed to parse item link, expected %s got: %s", "https://www.bank.lv/", item.Link)
	}

	// Make sure item guid is parsed correctly
	if item.GUID != "https://www.bank.lv/#06.05" {
		t.Errorf("Failed to parse item guid, expected %s got: %s", "https://www.bank.lv/#06.05", item.GUID)
	}

	// Make sure item description is parsed correctly
	if len(item.Description) != 502 {
		t.Errorf("Invalid length of description, expected %d got: %d", 502, len(item.Description))
	}

	start := item.Description[0:3]
	end := item.Description[len(item.Description)-4:]

	// Validate start of description
	if start != "AUD" {
		t.Errorf("Failed to parse item description, expected %s got: %s", "AUD", start)
	}

	// Validate end of description
	if end != "000 " {
		t.Errorf("Failed to parse item description, expected %s got: %s", "000 ", end)
	}
}
