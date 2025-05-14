package utils

import (
	"fmt"
	"sync"
	"time"
)

type CircuitBreaker struct {
	maxFailures  int
	resetTimeout time.Duration
	failureCount int
	lastFailure  time.Time
	mu           sync.Mutex
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	if cb.failureCount >= cb.maxFailures && time.Since(cb.lastFailure) < cb.resetTimeout {
		cb.mu.Unlock()
		return fmt.Errorf("circuit breaker is open")
	}
	cb.mu.Unlock()

	err := fn()
	if err != nil {
		cb.mu.Lock()
		cb.failureCount++
		cb.lastFailure = time.Now()
		cb.mu.Unlock()
		return err
	}

	cb.mu.Lock()
	cb.failureCount = 0
	cb.mu.Unlock()
	return nil
}
