package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ConnectDBs establishes connections to multiple PostgreSQL databases.
func ConnectDBs(configs map[string]Config) (map[string]*sql.DB, error) {
	dbs := make(map[string]*sql.DB)
	var firstErr error

	for name, config := range configs {
		connStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode,
		)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("failed to connect to database %s: %v", name, err)
			}
			continue
		}

		if err := db.Ping(); err != nil {
			db.Close()
			if firstErr == nil {
				firstErr = fmt.Errorf("failed to ping database %s: %v", name, err)
			}
			continue
		}

		dbs[name] = db
	}

	if len(dbs) == 0 && firstErr != nil {
		return nil, firstErr
	}

	return dbs, firstErr
}

// ConnectDB establishes a connection to a single PostgreSQL database.
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
