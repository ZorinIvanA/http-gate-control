package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZorinIvanA/http-gate-control/internal/client"
	"github.com/ZorinIvanA/http-gate-control/internal/config"
	"github.com/ZorinIvanA/http-gate-control/internal/handler"
	"github.com/ZorinIvanA/http-gate-control/internal/logger"
	"github.com/ZorinIvanA/http-gate-control/internal/metrics"
	"github.com/ZorinIvanA/http-gate-control/internal/service"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.MustLoad()

	// Инициализация метрик
	metrics.MustRegister()

	// Инициализация клиентов
	relayClient := client.NewRelayClient(cfg.RelayURL)
	loggerClient := logger.NewHTTPLogger(cfg.LoggerURL)

	// Инициализация сервиса
	gateService := service.NewGateService(cfg.OpenDelay)

	// Инициализация хендлера
	gateHandler := handler.NewGateHandler(gateService, relayClient, loggerClient)

	// Инициализация роутера
	router := chi.NewRouter()

	// Middleware
	router.Use(metrics.PrometheusMiddleware)
	router.Use(handler.RateLimiterMiddleware(cfg.RateLimit))

	// Маршруты
	router.Post("/open", gateHandler.HandleOpen)
	router.Mount("/metrics", promhttp.Handler())

	// Сервер
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Горутина запуска сервера
	go func() {
		log.Printf("Server started on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Ожидание сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
