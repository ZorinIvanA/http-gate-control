package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    ServerPort         string
    RelayIP            string
    BlockDuration      time.Duration
    LogPath            string
    MetricsFlushPeriod time.Duration
    SwaggerEnabled     bool
}

func Load() (*Config, error) {
    blockDuration, _ := strconv.Atoi(getenv("BLOCK_DURATION_SECONDS", "10"))
    metricsFlushPeriod, _ := strconv.Atoi(getenv("METRICS_FLUSH_PERIOD_SECONDS", "3600"))
    swaggerEnabled := getenv("SWAGGER_ENABLED", "true") == "true"

    return &Config{
        ServerPort:         getenv("SERVER_PORT", "8080"),
        RelayIP:            getenv("RELAY_IP", "192.168.1.100"),
        BlockDuration:      time.Duration(blockDuration),
        LogPath:            getenv("LOG_PATH", "/var/log/gate-service.log"),
        MetricsFlushPeriod: time.Duration(metricsFlushPeriod),
        SwaggerEnabled:     swaggerEnabled,
    }, nil
}

func getenv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}