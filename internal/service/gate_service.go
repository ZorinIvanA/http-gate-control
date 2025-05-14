package service

import (
	"sync"
	"time"
)

type GateService interface {
	ProcessOpen() (bool, error)
}

type gateService struct {
	lastOpenTime time.Time
	openDelay    time.Duration
	blockedCount uint64
	mu           sync.Mutex
	relayClient  RelayClient
	loggerClient LoggerClient
	metrics      Metrics
}

func NewGateService(relayClient RelayClient, loggerClient LoggerClient, metrics Metrics, openDelay time.Duration) GateService {
	return &gateService{
		openDelay:    openDelay,
		relayClient:  relayClient,
		loggerClient: loggerClient,
		metrics:      metrics,
	}
}

func (s *gateService) ProcessOpen() (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if now.Sub(s.lastOpenTime) < s.openDelay {
		s.blockedCount++
		s.metrics.IncBlocked()
		return false, nil
	}

	if err := s.relayClient.SendOpen(); err != nil {
		return false, err
	}

	s.lastOpenTime = now
	s.blockedCount = 0
	return true, nil
}
