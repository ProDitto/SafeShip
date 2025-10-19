package usecase

import (
	"context"
	"fmt"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"
)

type BuildCompletionRequest struct {
	ImageDigest string            `json:"image_digest"`
	Tags        []string          `json:"tags"`
	SLSALevel   int               `json:"slsa_level"`
	SBOMs       []SBOMInput       `json:"sboms"`
	CVEs        []CVEFindingInput `json:"cves"`
}

type SBOMInput struct {
	Format string `json:"format"`
	URI    string `json:"uri"`
}

type CVEFindingInput struct {
	CVEID        string `json:"cve_id"`
	Severity     string `json:"severity"`
	Description  string `json:"description"`
	FixAvailable bool   `json:"fix_available"`
}

type BuildUsecase struct {
	buildRepo repository.BuildEventRepository
	imageRepo repository.ImageRepository
	sbomRepo  repository.SBOMRecordRepository
	cveRepo   repository.CVEFindingRepository
}

func NewBuildUsecase(
	buildRepo repository.BuildEventRepository,
	imageRepo repository.ImageRepository,
	sbomRepo repository.SBOMRecordRepository,
	cveRepo repository.CVEFindingRepository,
) *BuildUsecase {
	return &BuildUsecase{
		buildRepo: buildRepo,
		imageRepo: imageRepo,
		sbomRepo:  sbomRepo,
		cveRepo:   cveRepo,
	}
}

func (uc *BuildUsecase) CompleteBuild(ctx context.Context, buildID int, req BuildCompletionRequest) (*domain.Image, error) {
	// 1. Find the build event and validate its state
	buildEvent, err := uc.buildRepo.FindByID(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("build event not found: %w", err)
	}
	if buildEvent.Status != "pending" {
		return nil, fmt.Errorf("build event %d is not in 'pending' state, current state: %s", buildID, buildEvent.Status)
	}

	// 2. Create the new image record
	newImage := &domain.Image{
		TenantNamespace: buildEvent.TenantNamespace,
		Digest:          req.ImageDigest,
		Tags:            req.Tags,
		SLSALevel:       req.SLSALevel,
	}
	imageID, err := uc.imageRepo.Create(ctx, newImage)
	if err != nil {
		return nil, fmt.Errorf("failed to create image record: %w", err)
	}
	newImage.ID = imageID

	// 3. Create SBOM records
	for _, sbomInput := range req.SBOMs {
		sbomRecord := &domain.SBOMRecord{
			ImageID: imageID,
			Format:  sbomInput.Format,
			URI:     sbomInput.URI,
		}
		if err := uc.sbomRepo.Create(ctx, sbomRecord); err != nil {
			// In a real system, we'd want transactional behavior here
			return nil, fmt.Errorf("failed to create sbom record: %w", err)
		}
	}

	// 4. Create CVE finding records
	var cveFindings []*domain.CVEFinding
	for _, cveInput := range req.CVEs {
		cveFindings = append(cveFindings, &domain.CVEFinding{
			ImageID:      imageID,
			CVEID:        cveInput.CVEID,
			Severity:     cveInput.Severity,
			Description:  cveInput.Description,
			FixAvailable: cveInput.FixAvailable,
		})
	}
	if err := uc.cveRepo.CreateBatch(ctx, cveFindings); err != nil {
		return nil, fmt.Errorf("failed to create cve findings: %w", err)
	}

	// 5. Update the build event with the new image_id and set status to 'completed'
	buildEvent.ImageID = &imageID
	buildEvent.Status = "completed"
	if err := uc.buildRepo.Update(ctx, buildEvent); err != nil {
		return nil, fmt.Errorf("failed to update build event: %w", err)
	}

	return newImage, nil
}

