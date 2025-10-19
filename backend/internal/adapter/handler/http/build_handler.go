package http

import (
	"encoding/json"
	"net/http"
	"secure-image-service/internal/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BuildHandler struct {
	usecase *usecase.BuildUsecase
}

func NewBuildHandler(uc *usecase.BuildUsecase) *BuildHandler {
	return &BuildHandler{usecase: uc}
}

func (h *BuildHandler) CompleteBuild(w http.ResponseWriter, r *http.Request) {
	buildIDStr := chi.URLParam(r, "buildID")
	buildID, err := strconv.Atoi(buildIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid build ID")
		return
	}

	var req usecase.BuildCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	image, err := h.usecase.CompleteBuild(r.Context(), buildID, req)
	if err != nil {
		// Consider more specific error codes based on err type
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, image)
}

