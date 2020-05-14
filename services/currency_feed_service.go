package services

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/matiss/go-microservice-test/utils"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type CurrencyFeedService struct {
	mysql        *MySQLService
	url          string
	Feed         utils.CurrencyFeed
	pendingItems []utils.FeedItem
	oldestDate   time.Time
}

// NewCurrencyFeedService creates new CurrencyFeedService
func NewCurrencyFeedService(mysql *MySQLService, url string) *CurrencyFeedService {
	return &CurrencyFeedService{
		mysql: mysql,
		url:   url,
	}
}

// Fetch latest currency RSS2 feed and store it locally
func (s *CurrencyFeedService) Fetch() error {
	// Fetch file and read its contents
	res, err := http.Get(s.url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Parse RSS2
	rss2 := utils.RSS2{}

	err = rss2.Parse(b)
	if err != nil {
		return err
	}

	// Parse feed
	feed := utils.CurrencyFeed{}

	// Iterate over feed items and parse currencies
	for _, item := range rss2.Items {
		err := feed.ParseItem(item.PubDate, item.Description)
		if err != nil {
			return err
		}
	}

	// Store feed locally
	s.Feed = feed

	return nil
}

// Check if feed needs to persited in database
func (s *CurrencyFeedService) Check() (bool, error) {

	// Map feed item dates
	dates := make([]time.Time, 0)
	count := len(s.Feed.Items)

	// Limit count to 50 for safety reasons
	if count > 50 {
		count = 50
	}

	for i := 0; i < count; i++ {
		dates = append(dates, s.Feed.Items[i].Date)
	}

	// Get current MySQL session
	sess := s.mysql.Session()

	// Get stored dates in database
	query, args, err := sqlx.In("SELECT published_at FROM feed_updates WHERE published_at IN (?) LIMIT 10;", dates)
	if err != nil {
		return false, err
	}

	rows, err := sess.Queryx(query, args...)

	storedDates := make([]time.Time, 0)

	for rows.Next() {
		var pubAt mysql.NullTime
		err = rows.Scan(&pubAt)
		if err != nil {
			return false, err
		}

		if pubAt.Valid {
			storedDates = append(storedDates, pubAt.Time)
		}
	}

	// Get feed dates that are not in database
	pendingItems := make([]utils.FeedItem, 0)

	for _, item := range s.Feed.Items {
		found := false

		// Iterate over dates from database and compare to item date
		for _, date := range storedDates {
			if item.Date.Equal(date) {
				found = true
			}
		}

		// Append items that are not in database
		if !found {
			pendingItems = append(pendingItems, item)
		}
	}

	s.pendingItems = pendingItems

	return (len(pendingItems) > 0), nil
}

// Persist feed data into database
func (s *CurrencyFeedService) Persist() error {
	if len(s.pendingItems) == 0 {
		return nil
	}

	// Get current MySQL session
	sess := s.mysql.Session()

	timestamp := time.Now().UTC()
	latestDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	latestIndex := 0

	for i, item := range s.pendingItems {
		// Get index for latest date
		if item.Date.After(latestDate) {
			latestDate = item.Date
			latestIndex = i
		}

		// Insert timestamp in feed_updates table
		_, err := sess.Exec("INSERT INTO feed_updates (updated_at, published_at) VALUES (?, ?);", timestamp, item.Date)
		if err != nil {
			return err
		}

		// Insert currencies in currencies table
		for _, currency := range item.Currencies {
			// Convert float64 to int64. Stored in micro currency 0.000001
			valueInt := int64(currency.Value * currencyStoreFactor)

			_, err := sess.Exec("INSERT INTO currencies (symbol, value, date) VALUES (?, ?, ?);", currency.Symbol, valueInt, item.Date)
			if err != nil {
				return err
			}
		}
	}

	// Update currencies_latest table with latest values
	for _, currency := range s.pendingItems[latestIndex].Currencies {
		// Convert float64 to int64. Stored in micro currency 0.000001
		valueInt := int64(currency.Value * currencyStoreFactor)

		_, err := sess.Exec("INSERT INTO currencies_latest (symbol, value) VALUES (?, ?) ON DUPLICATE KEY UPDATE value = ?;", currency.Symbol, valueInt, valueInt)
		if err != nil {
			return err
		}
	}

	// Clear pending items
	s.pendingItems = make([]utils.FeedItem, 0)

	return nil
}
