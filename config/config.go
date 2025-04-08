package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	envDBHost     = "DB_HOST"
	envDBPort     = "DB_PORT"
	envDBName     = "DB_NAME"
	envDBUser     = "DB_USER"
	envDBPassword = "DB_PASSWORD"
	envDBSSLMode  = "DB_SSLMODE"
)

// Config holds the database configuration.
type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	config := &Config{
		Host:     os.Getenv(envDBHost),
		Port:     os.Getenv(envDBPort),
		Name:     os.Getenv(envDBName),
		User:     os.Getenv(envDBUser),
		Password: os.Getenv(envDBPassword),
		SSLMode:  os.Getenv(envDBSSLMode),
	}

	// Set defaults
	if config.Host == "" {
		config.Host = "localhost"
	}

	if config.Port == "" {
		config.Port = "5432"
	}

	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}

	// Validate required fields
	if config.Name == "" || config.User == "" || config.Password == "" {
		return nil, fmt.Errorf("missing required database configuration (DB_NAME, DB_USER, DB_PASSWORD)")
	}

	return config, nil
}
