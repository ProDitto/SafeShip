package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"secure-image-service/internal/usecase"
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

// GetImageSBOMs returns mocked SBOM data for a given image ID.
func (h *ImageHandler) GetImageSBOMs(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	// In a real app, you'd query the database based on the ID.
	// Here, we return mock data for the demo.
	sboms := []map[string]string{
		{"format": "SPDX", "uri": "minio://sboms/image-" + idStr + ".spdx.json"},
		{"format": "CycloneDX", "uri": "minio://sboms/image-" + idStr + ".cdx.json"},
	}
	respondWithJSON(w, http.StatusOK, sboms)
}

// GetImageCVEs returns mocked CVE data for a given image ID.
func (h *ImageHandler) GetImageCVEs(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	// Mock data varies by ID for a better demo
	var cves []map[string]interface{}
	if idStr == "1" {
		cves = []map[string]interface{}{
			{"cve_id": "CVE-2023-4567", "severity": "Critical", "fix_available": true},
			{"cve_id": "CVE-2023-8910", "severity": "Medium", "fix_available": false},
		}
	} else if idStr == "3" {
		cves = []map[string]interface{}{
			{"cve_id": "CVE-2024-0001", "severity": "High", "fix_available": true},
		}
	} else {
		cves = []map[string]interface{}{} // No CVEs for other images
	}
	respondWithJSON(w, http.StatusOK, cves)
}

// GetImageVerification returns mocked verification data for a given image ID.
func (h *ImageHandler) GetImageVerification(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	verificationData := map[string]interface{}{
		"signature": map[string]string{
			"keyId":     "gcpkms://projects/secure-project/locations/global/keyRings/cosign/cryptoKeys/prod-key",
			"signature": "MEUCIQ...",
		},
		"attestations": []map[string]string{
			{"type": "provenance", "uri": "minio://attestations/image-" + idStr + "-provenance.json"},
			{"type": "slsa-v1.0", "uri": "minio://attestations/image-" + idStr + "-slsa.json"},
		},
		"rekorEntry": "https://rekor.sigstore.dev/api/v1/log/entries/e2a5...cfa1",
	}
	respondWithJSON(w, http.StatusOK, verificationData)
}

// Helpers from server.go, duplicated here for brevity in this context.
// In a real app, these would be in a shared package.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
