package utils

import (
	"encoding/xml"
)

type RSS2 struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`

	// Required
	Title       string `xml:"channel>title"`
	Link        string `xml:"channel>link"`
	Description string `xml:"channel>description"`

	// Optional
	BuiltAt   string `xml:"channel>lastBuildDate"`
	Generator string `xml:"channel>generator"`
	Image     Image  `xml:"channel>image"`
	Language  string `xml:"channel>language"`
	TTL       int    `xml:"channel>ttl"`
	Items     []Item `xml:"channel>item"`
}

type Image struct {
	Url    string `xml:"url"`
	Title  string `xml:"title"`
	Link   string `xml:"link"`
	Width  int    `xml:"width"`
	Height int    `xml:"height"`
}

type Item struct {
	// Required
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`

	// Optional
	GUID    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
}

// Parse XML/RSS2 file contents
func (r *RSS2) Parse(content []byte) error {
	err := xml.Unmarshal(content, r)
	if err != nil {
		return err
	}

	return nil
}
