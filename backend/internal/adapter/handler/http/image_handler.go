package http

import (
	"net/http"
	"secure-image-service/internal/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ImageHandler struct {
	usecase *usecase.ImageUsecase
}

func NewImageHandler(uc *usecase.ImageUsecase) *ImageHandler {
	return &ImageHandler{usecase: uc}
}

func (h *ImageHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	images, err := h.usecase.ListImages(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve images")
		return
	}
	respondWithJSON(w, http.StatusOK, images)
}

func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid image ID")
		return
	}

	image, err := h.usecase.GetImage(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve image")
		return
	}
	if image == nil {
		respondWithError(w, http.StatusNotFound, "Image not found")
		return
	}
	respondWithJSON(w, http.StatusOK, image)
}

