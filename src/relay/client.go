package relay

import (
    "fmt"
    "net/http"
)

// Client представляет клиент для взаимодействия с реле
type Client struct {
    relayIP string
}

// NewClient создает новый клиент реле
func NewClient(relayIP string) *Client {
    return &Client{relayIP: relayIP}
}

// OpenGate отправляет GET-запрос на реле для открытия шлагбаума
func (c *Client) OpenGate() error {
    url := fmt.Sprintf("http://%s/relay", c.relayIP)
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}