package services

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQLService struct {
	db *sqlx.DB
}

// NewMySQLService creates new MySQLService
func NewMySQLService(user string, password string, address string, database string) (*MySQLService, error) {
	if user == "" {
		return nil, fmt.Errorf("user is required")
	}

	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	if database == "" {
		return nil, fmt.Errorf("database is required")
	}

	db := sqlx.MustConnect("mysql", fmt.Sprintf("%s:%s@%s/%s", user, password, address, database))

	// Configure connection pool
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(3)

	return &MySQLService{
		db: db,
	}, nil
}

// Session returns current database session
func (s *MySQLService) Session() *sqlx.DB {
	return s.db
}

// Close current database session
func (s *MySQLService) Close() {
	if s.db == nil {
		return
	}

	s.db.Close()
	s.db = nil
}

// Setup creates database tables if needed
func (s *MySQLService) Setup() error {
	// Create feed_updates table
	s.db.MustExec(`
		CREATE TABLE IF NOT EXISTS feed_updates (
			updated_at DATETIME NULL DEFAULT NULL,
			published_at DATETIME NULL DEFAULT NULL
		);
	`)

	s.db.MustExec(`
		CREATE UNIQUE INDEX IF NOT EXISTS feed_updates_published_at_idx ON feed_updates (published_at);
	`)

	// Create currencies table
	s.db.MustExec(`
		CREATE TABLE IF NOT EXISTS currencies (
			symbol CHAR(3) NULL DEFAULT NULL,
			value BIGINT NULL DEFAULT NULL,
			date DATETIME NULL DEFAULT NULL
		);
	`)

	s.db.MustExec(`
		CREATE INDEX IF NOT EXISTS currencies_symbol_date_idx ON currencies (symbol, date);
	`)

	// Create currencies_latest
	s.db.MustExec(`
		CREATE TABLE IF NOT EXISTS currencies_latest (
			symbol CHAR(3) PRIMARY KEY,
			value BIGINT NULL DEFAULT NULL
		);
	`)

	return nil
}
