package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Logger интерфейс для логгирования событий
type Logger interface {
	LogEvent(ctx context.Context, event string, metadata map[string]interface{})
}

// HTTPLogger реализация Logger через HTTP
type HTTPLogger struct {
	url string
}

func NewHTTPLogger(url string) Logger {
	return &HTTPLogger{
		url: url,
	}
}

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Event     string                 `json:"event"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

func (l *HTTPLogger) LogEvent(ctx context.Context, event string, metadata map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Event:     event,
		Metadata:  metadata,
	}

	data, _ := json.Marshal(entry)

	req, err := http.NewRequestWithContext(ctx, "POST", l.url, bytes.NewReader(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
