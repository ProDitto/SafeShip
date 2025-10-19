package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	http_handler "secure-image-service/internal/adapter/handler/http"
	"secure-image-service/internal/adapter/postgres"
	"secure-image-service/internal/adapter/simulator"
	"secure-image-service/internal/usecase"
	"secure-image-service/pkg/config"
	"secure-image-service/pkg/logger"
)

func main() {
	// Load .env file for local development
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	// Initialize logger
	appLogger := logger.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Set up context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Establish database connection
	dbPool, err := postgres.NewConnection(ctx, cfg.DatabaseURL)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer dbPool.Close()
	appLogger.Info().Msg("Database connection established")

	// Initialize repositories
	imageRepo := postgres.NewImageRepository(dbPool)
	customerRepo := postgres.NewCustomerRepository(dbPool)
	buildEventRepo := postgres.NewBuildEventRepository(dbPool)
	sbomRepo := postgres.NewSBOMRecordRepository(dbPool)
	cveRepo := postgres.NewCVEFindingRepository(dbPool)

	// Initialize simulators
	orchestrator := simulator.NewMockBuildOrchestrator()

	// Initialize use cases
	imageUsecase := usecase.NewImageUsecase(imageRepo, buildEventRepo, orchestrator)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	buildUsecase := usecase.NewBuildUsecase(buildEventRepo, imageRepo, sbomRepo, cveRepo)

	// Initialize HTTP server
	server := http_handler.NewServer(imageUsecase, customerUsecase, buildUsecase, appLogger)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.APIPort),
		Handler: server,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Info().Msgf("Server starting on port %s", cfg.APIPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown
	appLogger.Info().Msg("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		appLogger.Fatal().Err(err).Msg("Server shutdown failed")
	}

	appLogger.Info().Msg("Server gracefully stopped")
}
