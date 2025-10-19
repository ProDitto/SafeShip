package http

import (
	"net/http"
	"secure-image-service/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type CustomerHandler struct {
	usecase *usecase.CustomerUsecase
}

func NewCustomerHandler(uc *usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{usecase: uc}
}

func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.usecase.ListCustomers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve customers")
		return
	}
	respondWithJSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	customer, err := h.usecase.GetCustomer(r.Context(), namespace)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve customer")
		return
	}
	if customer == nil {
		respondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}
	respondWithJSON(w, http.StatusOK, customer)
}

