package repository

import (
	"context"
	"secure-image-service/internal/domain"
)

type ImageRepository interface {
	FindAll(ctx context.Context) ([]*domain.Image, error)
	FindByID(ctx context.Context, id int) (*domain.Image, error)
}

