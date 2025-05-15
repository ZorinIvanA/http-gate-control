package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port      string
	RelayURL  string
	LoggerURL string
	OpenDelay time.Duration
	RateLimit int
}

func MustLoad() *Config {
	delayStr := getEnv("OPEN_DELAY", "5")
	delay, err := strconv.Atoi(delayStr)
	if err != nil {
		panic(fmt.Errorf("invalid OPEN_DELAY: %w", err))
	}

	rateLimitStr := getEnv("RATE_LIMIT", "100")
	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil {
		panic(fmt.Errorf("invalid RATE_LIMIT: %w", err))
	}

	return &Config{
		Port:      getEnv("PORT", "8080"),
		RelayURL:  os.Getenv("RELAY_URL"),
		LoggerURL: os.Getenv("LOGGER_URL"),
		OpenDelay: time.Duration(delay) * time.Second,
		RateLimit: rateLimit,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
