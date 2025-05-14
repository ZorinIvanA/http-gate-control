package client

import (
	"context"
	"net/http"
	"time"
)

type RelayClient interface {
	SendOpen() error
}

type relayClient struct {
	url     string
	client  *http.Client
	circuit *CircuitBreaker
}

func NewRelayClient(url string) RelayClient {
	return &relayClient{
		url: url,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
		circuit: NewCircuitBreaker(3, 5*time.Second),
	}
}

func (c *relayClient) SendOpen() error {
	return c.circuit.Execute(func() error {
		req, _ := http.NewRequestWithContext(context.Background(), "GET", c.url, nil)
		resp, err := c.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil
	})
}
