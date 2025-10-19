package repository

import (
	"context"
	"secure-image-service/internal/domain"
)

type SBOMRecordRepository interface {
	Create(ctx context.Context, sbom *domain.SBOMRecord) error
}

