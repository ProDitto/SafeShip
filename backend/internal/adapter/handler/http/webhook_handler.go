package http_adapter

import (
	"encoding/json"
	"net/http"
	"secure-image-service/backend/internal/usecase"
)

type WebhookHandler struct {
	usecase *usecase.ImageUsecase
}

func NewWebhookHandler(uc *usecase.ImageUsecase) *WebhookHandler {
	return &WebhookHandler{usecase: uc}
}

type UpstreamWebhookRequest struct {
	TenantNamespace string `json:"tenant_namespace"`
}

func (h *WebhookHandler) TriggerUpstreamBuild(w http.ResponseWriter, r *http.Request) {
	var req UpstreamWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.TenantNamespace == "" {
		respondWithError(w, http.StatusBadRequest, "tenant_namespace is required")
		return
	}

	buildEvent, err := h.usecase.CreateBuild(r.Context(), req.TenantNamespace, "webhook")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusAccepted, buildEvent)
}

