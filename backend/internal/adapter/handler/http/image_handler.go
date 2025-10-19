package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"secure-image-service/internal/usecase"
)

// ImageHandler handles HTTP requests for images.
type ImageHandler struct {
	usecase *usecase.ImageUsecase
}

// NewImageHandler creates a new ImageHandler.
func NewImageHandler(uc *usecase.ImageUsecase) *ImageHandler {
	return &ImageHandler{usecase: uc}
}

type CreateBuildRequest struct {
	TenantNamespace string `json:"tenant_namespace"`
}

// CreateBuild handles POST requests to trigger a new image build.
func (h *ImageHandler) CreateBuild(w http.ResponseWriter, r *http.Request) {
	var req CreateBuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if req.TenantNamespace == "" {
		respondWithError(w, http.StatusBadRequest, "tenant_namespace is required")
		return
	}

	// For now, all builds are triggered via API
	triggerType := "api"

	buildEvent, err := h.usecase.CreateBuild(r.Context(), req.TenantNamespace, triggerType)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusAccepted, buildEvent)
}

// ListImages handles GET requests to retrieve all images.
func (h *ImageHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	images, err := h.usecase.ListImages(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve images")
		return
	}
	respondWithJSON(w, http.StatusOK, images)
}

// GetImage handles GET requests for a specific image.
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid image ID")
		return
	}

	image, err := h.usecase.GetImage(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve image")
		return
	}
	if image == nil {
		respondWithError(w, http.StatusNotFound, "Image not found")
		return
	}

	respondWithJSON(w, http.StatusOK, image)
}

// GetImageSBOMs returns mocked SBOM data for a given image ID.
func (h *ImageHandler) GetImageSBOMs(w http.ResponseWriter, r *http.Request) {
	// This is mocked for the MVP. A real implementation would fetch from the db.
	mockSBOMs := []map[string]string{
		{"format": "SPDX", "uri": "minio://sboms/image-1/sbom.spdx.json"},
		{"format": "CycloneDX", "uri": "minio://sboms/image-1/sbom.cyclonedx.json"},
	}
	respondWithJSON(w, http.StatusOK, mockSBOMs)
}

// GetImageCVEs returns mocked CVE data for a given image ID.
func (h *ImageHandler) GetImageCVEs(w http.ResponseWriter, r *http.Request) {
	// This is mocked for the MVP. A real implementation would fetch from the db.
	idStr := chi.URLParam(r, "id")
	var mockCVEs []map[string]interface{}
	if idStr == "1" {
		mockCVEs = []map[string]interface{}{
			{"cve_id": "CVE-2023-12345", "severity": "High", "fix_available": true},
			{"cve_id": "CVE-2023-67890", "severity": "Medium", "fix_available": false},
		}
	} else {
		mockCVEs = []map[string]interface{}{
			{"cve_id": "CVE-2023-54321", "severity": "Critical", "fix_available": true},
		}
	}
	respondWithJSON(w, http.StatusOK, mockCVEs)
}

// GetImageVerification returns mocked verification data for a given image ID.
func (h *ImageHandler) GetImageVerification(w http.ResponseWriter, r *http.Request) {
	// This is mocked for the MVP.
	mockVerification := map[string]interface{}{
		"signature": map[string]string{
			"key_id":    "gcpkms://projects/secure-project/locations/global/keyRings/cosign/cryptoKeys/prod-key",
			"signature": "MEUCIQ...",
		},
		"rekor_entry": "https://rekor.sigstore.dev/api/v1/log/entries/...",
		"attestations": []map[string]string{
			{"type": "provenance", "uri": "minio://attestations/image-1/provenance.json"},
			{"type": "vuln-scan", "uri": "minio://attestations/image-1/scan-report.json"},
		},
	}
	respondWithJSON(w, http.StatusOK, mockVerification)
}

// respondWithError is a helper to write a JSON error response.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON is a helper to write a JSON response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
