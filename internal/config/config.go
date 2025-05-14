package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	RelayURL  string
	LoggerURL string
	OpenDelay time.Duration
	Port      string
}

func NewConfig() (*Config, error) {
	openDelay := 5
	if delayStr := os.Getenv("OPEN_DELAY"); delayStr != "" {
		_, err := fmt.Sscanf(delayStr, "%d", &openDelay)
		if err != nil {
			return nil, err
		}
	}

	relayURL := os.Getenv("RELAY_URL")
	if relayURL == "" {
		return nil, fmt.Errorf("RELAY_URL is required")
	}

	return &Config{
		RelayURL:  relayURL,
		LoggerURL: os.Getenv("LOGGER_URL"),
		OpenDelay: time.Duration(openDelay) * time.Second,
		Port:      os.Getenv("PORT"),
	}, nil
}
