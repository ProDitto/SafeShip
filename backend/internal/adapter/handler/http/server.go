package http

import (
	"encoding/json"
	"net/http"
	"secure-image-service/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
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

	s.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	s.Router.Route("/v1", func(r chi.Router) {
		r.Get("/images", s.ImageHandler.ListImages)
		r.Get("/images/{id}", s.ImageHandler.GetImage)

		r.Get("/customers", s.CustomerHandler.ListCustomers)
		r.Get("/customers/{namespace}", s.CustomerHandler.GetCustomer)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

