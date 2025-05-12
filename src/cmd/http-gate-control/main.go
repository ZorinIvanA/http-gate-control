package main

import (
	"context"
	"net/http"

	httpgatecontrol "github.com/ZorinIvanA/http-gate-control"

	"github.com/ZorinIvanA/http-gate-control/barrel"
	"github.com/ZorinIvanA/http-gate-control/config"
	"github.com/ZorinIvanA/http-gate-control/handler"
	"github.com/ZorinIvanA/http-gate-control/logger"
	"github.com/ZorinIvanA/http-gate-control/metrics"
	"github.com/ZorinIvanA/http-gate-control/relay"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Инициализация логгера
	logger, err := logger.NewLogger(cfg.LogPath)
	if err != nil {
		panic(err)
	}

	// Инициализация метрик
	metrics := metrics.NewMetrics(cfg.MetricsFlushPeriod, cfg.LogPath)

	// Создание клиента реле
	relayClient := relay.NewClient(cfg.RelayIP)

	// Создание контроллера шлагбаума
	controller := barrel.NewController(cfg.BlockDuration, metrics)

	// Создание обработчика
	handler := handler.NewHandler(controller, *relayClient, logger, metrics)

	// Настройка маршрутов
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	handler.RegisterRoutes(r)

	// Добавление Swagger если включен
	if cfg.SwaggerEnabled {
		r.Get("/swagger/*", httpgatecontrol.SwaggerHandler())
	}

	// Запуск сервера
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	logger.Info("Starting server", zap.String("port", cfg.ServerPort))

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Ожидание завершения
	<-context.Background().Done()
	logger.Info("Shutting down server...")
}
