package repository

import (
	"context"
	"secure-image-service/internal/domain"
)

// ImageRepository defines the interface for accessing image data.
type ImageRepository interface {
	FindAll(ctx context.Context) ([]*domain.Image, error)
	FindByID(ctx context.Context, id int) (*domain.Image, error)
	Create(ctx context.Context, image *domain.Image) (int, error)
}
