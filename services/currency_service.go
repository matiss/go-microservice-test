package services

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type CurrencyService struct {
	mysql *MySQLService
}

type Currency struct {
	Symbol string         `db:"symbol" json:"symbol"`
	Value  int64          `db:"value" json:"value"`
	Date   mysql.NullTime `db:"date" json:"date,omitempty"`
}

// NewCurrencyService creates new CurrencyService
func NewCurrencyService(mysql *MySQLService) *CurrencyService {
	return &CurrencyService{
		mysql: mysql,
	}
}

// Latest currencies
func (s *CurrencyService) Latest() ([]Currency, error) {

	// Get current MySQL session
	sess := s.mysql.Session()

	// Fetch currencies by symbol
	currencies := make([]Currency, 0)
	err := sess.Select(&currencies, "SELECT * FROM currencies_latest LIMIT 100")
	if err != nil {
		return nil, err
	}

	return currencies, nil
}

// BySymbol returns historical values by symbol
func (s *CurrencyService) BySymbol(symbol string, limit int) ([]Currency, error) {
	// Make sure symbol is not empty string
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	// Constraint limit
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	// Get current MySQL session
	sess := s.mysql.Session()

	// Fetch currencies by symbol
	currencies := make([]Currency, 0)
	err := sess.Select(&currencies, "SELECT * FROM currencies WHERE symbol = ? ORDER BY date DESC LIMIT ?", symbol, limit)
	if err != nil {
		return nil, err
	}

	return currencies, nil
}
