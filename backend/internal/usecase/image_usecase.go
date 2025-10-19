package usecase

import (
	"context"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"
)

type ImageUsecase struct {
	repo repository.ImageRepository
}

func NewImageUsecase(repo repository.ImageRepository) *ImageUsecase {
	return &ImageUsecase{repo: repo}
}

func (uc *ImageUsecase) ListImages(ctx context.Context) ([]*domain.Image, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *ImageUsecase) GetImage(ctx context.Context, id int) (*domain.Image, error) {
	return uc.repo.FindByID(ctx, id)
}

