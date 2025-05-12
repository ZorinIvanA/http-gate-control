package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	totalRequests   uint64
	blockedRequests uint64
	lastReset       time.Time
	mu              sync.Mutex
	flushPeriod     time.Duration
	logPath         string
}

func NewMetrics(flushPeriod time.Duration, logPath string) *Metrics {
	m := &Metrics{
		lastReset:   time.Now(),
		flushPeriod: flushPeriod * time.Second,
		logPath:     logPath,
	}

	// Запускаем периодический сброс метрик
	go m.startDailyReset()

	return m
}

func (m *Metrics) startDailyReset() {
	ticker := time.NewTicker(m.flushPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.ResetDaily()
		}
	}
}

func (m *Metrics) IncTotal() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalRequests++
}

func (m *Metrics) IncBlocked() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.blockedRequests++
}

func (m *Metrics) RPS() float64 {
	elapsed := time.Since(m.lastReset).Seconds()
	if elapsed == 0 {
		return 0
	}
	return float64(m.totalRequests) / elapsed
}

func (m *Metrics) DailyBlocked() uint64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.blockedRequests
}

func (m *Metrics) ResetDaily() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.blockedRequests = 0
	m.totalRequests = 0
	m.lastReset = time.Now()
}

func (m *Metrics) LastReset() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastReset
}
