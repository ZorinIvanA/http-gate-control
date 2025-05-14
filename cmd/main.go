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
	"github.com/ZorinIvanA/http-gate-control/internal/metrics"
	"github.com/ZorinIvanA/http-gate-control/internal/middleware"
	"github.com/ZorinIvanA/http-gate-control/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	relayClient := client.NewRelayClient(cfg.RelayURL)
	loggerClient := client.NewLoggerClient(cfg.LoggerURL)
	metrics := metrics.NewMetrics()
	gateService := service.NewGateService(relayClient, loggerClient, metrics, cfg.OpenDelay)

	r := gin.Default()
	r.Use(middleware.RateLimit(100, 1000))

	r.GET("/open", handler.OpenHandler(gateService))
	r.GET("/metrics", gin.WrapH(metrics.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
