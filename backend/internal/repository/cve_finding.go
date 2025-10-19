package repository

import (
	"context"
	"secure-image-service/internal/domain"
)

type CVEFindingRepository interface {
	CreateBatch(ctx context.Context, cves []*domain.CVEFinding) error
}

