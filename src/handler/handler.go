package handler

import (
	"encoding/json"
	"net/http"
	"time"

	httpgatecontrol "github.com/ZorinIvanA/http-gate-control"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ZorinIvanA/http-gate-control/barrel"
	"github.com/ZorinIvanA/http-gate-control/metrics"
	"github.com/ZorinIvanA/http-gate-control/relay"
)

type Handler struct {
	controller *barrel.Controller
	relay      relay.Client
	logger     *zap.Logger
	metrics    *metrics.Metrics
}

func NewHandler(controller *barrel.Controller, relay relay.Client, logger *zap.Logger, metrics *metrics.Metrics) *Handler {
	return &Handler{
		controller: controller,
		relay:      relay,
		logger:     logger,
		metrics:    metrics,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/open", h.HandleOpen)
	r.Get("/metrics", h.HandleMetrics)
	r.Get("/swagger/*", httpgatecontrol.SwaggerHandler())
}

// @Summary Open the barrier
// @Description Opens the barrier if cooldown period has passed
// @Tags barrier
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /open [post]
func (h *Handler) HandleOpen(w http.ResponseWriter, r *http.Request) {
	h.metrics.IncTotal()

	if !h.controller.ShouldOpen() {
		h.metrics.IncBlocked()
		h.logger.Info("Gate opening blocked due to cooldown")
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := h.relay.OpenGate(); err != nil {
		h.logger.Error("Failed to open gate", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info("Gate opened successfully")
	w.WriteHeader(http.StatusOK)
}

// @Summary Get service metrics
// @Description Returns service metrics including RPS and blocked requests
// @Tags metrics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /metrics [get]
func (h *Handler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]interface{}{
		"rps":           h.metrics.RPS(),
		"daily_blocked": h.metrics.DailyBlocked(),
		"last_reset":    h.metrics.LastReset().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
