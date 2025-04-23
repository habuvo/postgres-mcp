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
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"sslmode"`
}

// LoadConfigs loads multiple database configurations from a JSON string in environment variable.
func LoadConfigs(envVar string) (map[string]Config, error) {
	jsonConfig := os.Getenv(envVar)
	if jsonConfig == "" {
		return nil, fmt.Errorf("environment variable %s is empty", envVar)
	}

	var configs map[string]Config
	err := json.Unmarshal([]byte(jsonConfig), &configs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON config: %v", err)
	}

	// Validate each config
	for name, cfg := range configs {
		if cfg.Name == "" || cfg.User == "" || cfg.Password == "" {
			return nil, fmt.Errorf("missing required fields in config for database %s", name)
		}
		// Set defaults
		if cfg.Host == "" {
			cfg.Host = "localhost"
		}
		if cfg.Port == "" {
			cfg.Port = "5432"
		}
		if cfg.SSLMode == "" {
			cfg.SSLMode = "disable"
		}
		configs[name] = cfg
	}

	return configs, nil
}

// LoadConfig loads a single database configuration from environment variables.
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
