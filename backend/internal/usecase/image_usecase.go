package usecase

import (
	"context"
	"secure-image-service/internal/adapter/simulator"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"
)

// ImageUsecase handles business logic for images.
type ImageUsecase struct {
	repo         repository.ImageRepository
	buildRepo    repository.BuildEventRepository
	orchestrator simulator.BuildOrchestrator
}

// NewImageUsecase creates a new ImageUsecase.
func NewImageUsecase(repo repository.ImageRepository, buildRepo repository.BuildEventRepository, orchestrator simulator.BuildOrchestrator) *ImageUsecase {
	return &ImageUsecase{
		repo:         repo,
		buildRepo:    buildRepo,
		orchestrator: orchestrator,
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
	// 1. Create a build event record in the database
	event := &domain.BuildEvent{
		TenantNamespace: tenantNamespace,
		TriggerType:     triggerType,
		Status:          "pending",
	}
	_, err := uc.buildRepo.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	// 2. Trigger the build orchestrator
	if err := uc.orchestrator.TriggerBuild(ctx, event); err != nil {
		// In a real system, you might want to mark the build event as failed here.
		return nil, err
	}

	return event, nil
}
