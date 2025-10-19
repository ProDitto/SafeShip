package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	APIPort     string
}

func Load() (*Config, error) {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	// Validate port is a number
	if _, err := strconv.Atoi(port); err != nil {
		port = "8080"
	}

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		APIPort:     port,
	}, nil
}

