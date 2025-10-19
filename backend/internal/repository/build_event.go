package repository

import (
	"context"
	"secure-image-service/internal/domain"
)

type BuildEventRepository interface {
	Create(ctx context.Context, event *domain.BuildEvent) (int, error)
	FindByID(ctx context.Context, id int) (*domain.BuildEvent, error)
	Update(ctx context.Context, event *domain.BuildEvent) error
}

