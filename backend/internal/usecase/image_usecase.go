package usecase

import (
	"context"
	"secure-image-service/backend/internal/adapter/simulator"
	"secure-image-service/backend/internal/domain"
	"secure-image-service/backend/internal/repository"
)

// ImageUsecase handles business logic for images.
type ImageUsecase struct {
	repo         repository.ImageRepository
	buildRepo    repository.BuildEventRepository
	orchestrator simulator.BuildOrchestrator
	auditUC      *AuditUsecase
}

// NewImageUsecase creates a new ImageUsecase.
func NewImageUsecase(
	repo repository.ImageRepository,
	buildRepo repository.BuildEventRepository,
	orchestrator simulator.BuildOrchestrator,
	auditUC *AuditUsecase,
) *ImageUsecase {
	return &ImageUsecase{
		repo:         repo,
		buildRepo:    buildRepo,
		orchestrator: orchestrator,
		auditUC:      auditUC,
	}
}

// ListImages retrieves all images.
func (uc *ImageUsecase) ListImages(ctx context.Context) ([]*domain.Image, error) {
	return uc.repo.FindAll(ctx)
}

// GetImage retrieves a single image by its ID.
func (uc *ImageUsecase) GetImage(ctx context.Context, id int) (*domain.Image, error) {
	return uc.repo.FindByID(ctx, id)
}

// CreateBuild creates a build event and triggers a simulated build.
func (uc *ImageUsecase) CreateBuild(ctx context.Context, tenantNamespace, triggerType string) (*domain.BuildEvent, error) {
	buildEvent := &domain.BuildEvent{
		TenantNamespace: tenantNamespace,
		TriggerType:     triggerType,
		Status:          "pending",
	}

	id, err := uc.buildRepo.Create(ctx, buildEvent)
	if err != nil {
		return nil, err
	}
	buildEvent.ID = id

	// Trigger the build process asynchronously
	go uc.orchestrator.TriggerBuild(context.Background(), buildEvent)

	// Log the audit event
	auditDetails := map[string]interface{}{
		"build_event_id": id,
		"trigger_type":   triggerType,
	}
	// In a real system, actor would come from the request context (e.g., user ID, API key ID)
	_ = uc.auditUC.Log(ctx, tenantNamespace, "build_triggered", "system", auditDetails)

	return buildEvent, nil
}

