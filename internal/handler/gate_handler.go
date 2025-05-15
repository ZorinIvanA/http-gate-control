package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ZorinIvanA/http-gate-control/internal/client"  // Импорт клиента реле
	"github.com/ZorinIvanA/http-gate-control/internal/logger"  // Импорт логгера
	"github.com/ZorinIvanA/http-gate-control/internal/metrics" // Импорт метрик
	"github.com/ZorinIvanA/http-gate-control/internal/service"
	"golang.org/x/time/rate"
)

type GateHandler struct {
	service     service.GateServiceInterface
	relayClient client.RelayClient
	logger      logger.Logger
}

func NewGateHandler(
	service service.GateServiceInterface,
	relayClient client.RelayClient,
	logger logger.Logger,
) *GateHandler {
	return &GateHandler{
		service:     service,
		relayClient: relayClient,
		logger:      logger,
	}
}

func (h *GateHandler) HandleOpen(w http.ResponseWriter, r *http.Request) {
	metrics.RequestCounter.Inc() // Теперь metrics доступен

	if h.service.CheckAccess() {
		// Открываем шлагбаум
		if err := h.relayClient.OpenGate(r.Context()); err != nil {
			h.logger.LogEvent(r.Context(), "Relay failure", map[string]interface{}{
				"error": err.Error(),
			})
			http.Error(w, "Relay failure", http.StatusInternalServerError)
			return
		}

		h.logger.LogEvent(r.Context(), "Gate opened", map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
		})

		metrics.BlockedCounter.Set(0) // Теперь metrics доступен
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Gate opened")
	} else {
		blockedCount := h.service.GetBlockedCount()
		metrics.BlockedCounter.Inc() // Теперь metrics доступен

		h.logger.LogEvent(r.Context(), "Access blocked", map[string]interface{}{
			"blocked_count": blockedCount,
		})

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Access blocked")
	}
}

// RateLimiterMiddleware добавляет ограничение частоты запросов
func RateLimiterMiddleware(maxRequests int) func(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(maxRequests), maxRequests)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				metrics.RateLimitCounter.Inc()
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// BasicAuthMiddleware добавляет защиту для Swagger UI
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "password" {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
