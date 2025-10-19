package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"secure-image-service/internal/adapter/handler/http_handler"
	"secure-image-service/internal/adapter/postgres"
	"secure-image-service/internal/usecase"
	"secure-image-service/pkg/config"
	"secure-image-service/pkg/logger"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file for local development
	_ = godotenv.Load()

	log := logger.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Database connection
	dbpool, err := postgres.NewConnection(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer dbpool.Close()
	log.Info().Msg("Database connection established")

	// Initialize repositories
	imageRepo := postgres.NewImageRepository(dbpool)
	customerRepo := postgres.NewCustomerRepository(dbpool)

	// Initialize use cases
	imageUsecase := usecase.NewImageUsecase(imageRepo)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)

	// Initialize HTTP server
	server := http_handler.NewServer(imageUsecase, customerUsecase, log)
	httpServer := &http.Server{
		Addr:    ":" + cfg.APIPort,
		Handler: server.Router,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Msgf("Starting server on port %s", cfg.APIPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Info().Msg("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Server shutdown failed")
	}

	log.Info().Msg("Server gracefully stopped")
}

