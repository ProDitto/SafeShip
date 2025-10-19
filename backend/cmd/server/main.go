package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	http_adapter "secure-image-service/internal/adapter/handler/http"
	"secure-image-service/internal/adapter/postgres"
	"secure-image-service/internal/adapter/simulator"
	"secure-image-service/internal/usecase"
	"secure-image-service/pkg/config"
	"secure-image-service/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file for local development
	_ = godotenv.Load()

	// Initialize logger
	log := logger.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Set up context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Establish database connection
	dbPool, err := postgres.NewConnection(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer dbPool.Close()
	log.Info().Msg("Database connection established")

	// Initialize repositories
	imageRepo := postgres.NewImageRepository(dbPool)
	customerRepo := postgres.NewCustomerRepository(dbPool)
	buildEventRepo := postgres.NewBuildEventRepository(dbPool)
	sbomRepo := postgres.NewSBOMRecordRepository(dbPool)
	cveRepo := postgres.NewCVEFindingRepository(dbPool)
	auditRepo := postgres.NewAuditLogRepository(dbPool)

	// Initialize simulators
	orchestrator := simulator.NewMockBuildOrchestrator()

	// Initialize use cases
	auditUsecase := usecase.NewAuditUsecase(auditRepo)
	imageUsecase := usecase.NewImageUsecase(imageRepo, buildEventRepo, orchestrator, auditUsecase)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	buildUsecase := usecase.NewBuildUsecase(buildEventRepo, imageRepo, sbomRepo, cveRepo)

	// Initialize HTTP server
	server := http_adapter.NewServer(imageUsecase, customerUsecase, buildUsecase, log)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.APIPort),
		Handler: server,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Msgf("Server starting on port %s", cfg.APIPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown
	log.Info().Msg("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Server shutdown failed")
	}

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Error during server custom shutdown")
	}

	log.Info().Msg("Server gracefully stopped")
}
