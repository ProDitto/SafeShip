package repository

import (
	"context"
	"secure-image-service/internal/domain"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
}
