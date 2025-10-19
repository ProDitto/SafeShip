package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"secure-image-service/internal/usecase"
)

type Server struct {
	Router          *chi.Mux
	ImageHandler    *ImageHandler
	CustomerHandler *CustomerHandler
	Logger          zerolog.Logger
}

func NewServer(imageUsecase *usecase.ImageUsecase, customerUsecase *usecase.CustomerUsecase, logger zerolog.Logger) *Server {
	s := &Server{
		Router:          chi.NewRouter(),
		ImageHandler:    NewImageHandler(imageUsecase),
		CustomerHandler: NewCustomerHandler(customerUsecase),
		Logger:          logger,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(60 * time.Second))

	s.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	s.Router.Route("/v1", func(r chi.Router) {
		r.Route("/images", func(r chi.Router) {
			r.Get("/", s.ImageHandler.ListImages)
			r.Get("/{id}", s.ImageHandler.GetImage)
			r.Get("/{id}/sbom", s.ImageHandler.GetImageSBOMs)
			r.Get("/{id}/cves", s.ImageHandler.GetImageCVEs)
			r.Get("/{id}/verify", s.ImageHandler.GetImageVerification)
		})
		r.Route("/customers", func(r chi.Router) {
			r.Get("/", s.CustomerHandler.ListCustomers)
			r.Get("/{namespace}", s.CustomerHandler.GetCustomer)
		})
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) Shutdown(ctx context.Context) error {
	// In a real app, you'd have a http.Server instance to call Shutdown on.
	// For this structure, we don't have one to manage, but this is where it would go.
	s.Logger.Info().Msg("HTTP server shutting down")
	return nil
}
