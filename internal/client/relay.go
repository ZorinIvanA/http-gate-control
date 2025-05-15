package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// RelayClient интерфейс для управления реле
type RelayClient interface {
	OpenGate(ctx context.Context) error
}

// HTTPRelayClient реализация RelayClient через HTTP
type HTTPRelayClient struct {
	url string
}

func NewRelayClient(url string) RelayClient {
	return &HTTPRelayClient{
		url: url,
	}
}

func (c *HTTPRelayClient) OpenGate(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
