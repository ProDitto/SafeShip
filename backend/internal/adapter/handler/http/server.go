package http_adapter

import (
	"context"
	"net/http"
	"time"

	"secure-image-service/backend/internal/adapter/handler/http/middleware"
	"secure-image-service/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// Server represents the HTTP server.
type Server struct {
	Router          *chi.Mux
	ImageHandler    *ImageHandler
	CustomerHandler *CustomerHandler
	BuildHandler    *BuildHandler
	WebhookHandler  *WebhookHandler
	Logger          zerolog.Logger
}

// NewServer creates and configures a new Server instance.
func NewServer(
	imageUsecase *usecase.ImageUsecase,
	customerUsecase *usecase.CustomerUsecase,
	buildUsecase *usecase.BuildUsecase,
	logger zerolog.Logger,
) *Server {
	s := &Server{
		Router:          chi.NewRouter(),
		ImageHandler:    NewImageHandler(imageUsecase),
		CustomerHandler: NewCustomerHandler(customerUsecase),
		BuildHandler:    NewBuildHandler(buildUsecase),
		WebhookHandler:  NewWebhookHandler(imageUsecase),
		Logger:          logger,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Middleware
	s.Router.Use(chi_middleware.RequestID)
	s.Router.Use(chi_middleware.RealIP)
	s.Router.Use(chi_middleware.Logger)
	s.Router.Use(chi_middleware.Recoverer)
	s.Router.Use(chi_middleware.Timeout(60 * time.Second))

	// Public routes
	s.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// API v1 routes with authentication
	s.Router.Route("/v1", func(r chi.Router) {
		r.Use(middleware.Authenticator)

		// Image routes
		r.Get("/images", s.ImageHandler.ListImages)
		r.Post("/images", s.ImageHandler.CreateBuild)
		r.Get("/images/{id}", s.ImageHandler.GetImage)
		r.Get("/images/{id}/sbom", s.ImageHandler.GetImageSBOMs)
		r.Get("/images/{id}/cves", s.ImageHandler.GetImageCVEs)
		r.Get("/images/{id}/verify", s.ImageHandler.GetImageVerification)

		// Customer routes
		r.Get("/customers", s.CustomerHandler.ListCustomers)
		r.Get("/customers/{namespace}", s.CustomerHandler.GetCustomer)

		// Build routes
		r.Post("/builds/{buildID}/complete", s.BuildHandler.CompleteBuild)

		// Webhook routes
		r.Post("/webhooks/upstream", s.WebhookHandler.TriggerUpstreamBuild)
	})
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// Shutdown is a placeholder for graceful server shutdown.
func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info().Msg("Server shutdown complete.")
	return nil
}

