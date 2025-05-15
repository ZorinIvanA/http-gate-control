package service

import (
	"sync"
	"time"
)

// GateServiceInterface — интерфейс для управления шлагбаумом
type GateServiceInterface interface {
	CheckAccess() bool
	GetBlockedCount() int
}

type GateService struct {
	openDelay    time.Duration
	lastOpenTime time.Time
	blockedCount int
	mu           sync.Mutex
}

func NewGateService(openDelay time.Duration) GateServiceInterface {
	return &GateService{
		openDelay: openDelay,
	}
}

func (s *GateService) CheckAccess() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if now.Sub(s.lastOpenTime) < s.openDelay {
		s.blockedCount++
		return false
	}

	s.lastOpenTime = now
	s.blockedCount = 0
	return true
}

func (s *GateService) GetBlockedCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.blockedCount
}
