package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ConnectDB establishes a connection to the PostgreSQL database.
func ConnectDB(config *Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()

		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	return db, nil
}
