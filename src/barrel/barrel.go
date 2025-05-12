package barrel

import (
	"sync"
	"time"

	metrics "github.com/ZorinIvanA/http-gate-control/metrics"
)

// Controller управляет логикой открытия шлагбаума
type Controller struct {
	lastOpenTime  time.Time
	blockDuration time.Duration
	mu            sync.Mutex
	metrics       *metrics.Metrics
}

// NewController создает новый контроллер шлагбаума
func NewController(blockDuration time.Duration, metrics *metrics.Metrics) *Controller {
	return &Controller{
		blockDuration: blockDuration * time.Second,
		metrics:       metrics,
	}
}

// ShouldOpen проверяет, можно ли открывать шлагбаум
// Возвращает true, если прошло достаточно времени с последнего открытия
func (c *Controller) ShouldOpen() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	if now.Sub(c.lastOpenTime) >= c.blockDuration {
		c.lastOpenTime = now
		return true
	}
	return false
}
