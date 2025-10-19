package repository

import (
	"context"
	"secure-image-service/backend/internal/domain"
)

type SLAViolationRepository interface {
	Create(ctx context.Context, violation *domain.SLAViolation) error
}

