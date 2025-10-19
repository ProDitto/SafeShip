package http

import (
	"context"
	"net/http"
	"secure-image-service/internal/usecase"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// Server represents the HTTP server.
type Server struct {
	Router          *chi.Mux
	ImageHandler    *ImageHandler
	CustomerHandler *CustomerHandler
	BuildHandler    *BuildHandler
	Logger          zerolog.Logger
}

// NewServer creates and configures a new Server instance.
func NewServer(
	imageUsecase *usecase.ImageUsecase,
	customerUsecase *usecase.CustomerUsecase,
	buildUsecase *usecase.BuildUsecase,
	logger zerolog.ologer,
) *Server {
	s := &Server{
		Router:          chi.NewRouter(),
		ImageHandler:    NewImageHandler(imageUsecase),
		CustomerHandler: NewCustomerHandler(customerUsecase),
		BuildHandler:    NewBuildHandler(buildUsecase),
		Logger:          logger,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Middleware
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(60 * time.Second))

	// Health check
	s.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// API v1 routes
	s.Router.Route("/v1", func(r chi.Router) {
		// Image routes
		r.Get("/images", s.ImageHandler.ListImages)
		r.Post("/images", s.ImageHandler.CreateBuild)
		r.Get("/images/{id}", s.ImageHandler.GetImage)
		r.Get("/images/{id}/sboms", s.ImageHandler.GetImageSBOMs)
		r.Get("/images/{id}/cves", s.ImageHandler.GetImageCVEs)
		r.Get("/images/{id}/verification", s.ImageHandler.GetImageVerification)

		// Customer routes
		r.Get("/customers", s.CustomerHandler.ListCustomers)
		r.Get("/customers/{namespace}", s.CustomerHandler.GetCustomer)

		// Build routes
		r.Post("/builds/{buildID}/complete", s.BuildHandler.CompleteBuild)
	})
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// Shutdown is a placeholder for graceful server shutdown.
func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info().Msg("Server shutting down")
	return nil
}
