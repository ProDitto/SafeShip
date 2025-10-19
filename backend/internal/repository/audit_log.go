package repository

import (
	"context"
	"secure-image-service/backend/internal/domain"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
}

