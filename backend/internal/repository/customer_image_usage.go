package repository

import (
	"context"
	"secure-image-service/backend/internal/domain"
)

type CustomerImageUsageRepository interface {
	Create(ctx context.Context, usage *domain.CustomerImageUsage) error
	FindByTenant(ctx context.Context, tenantNamespace string) ([]*domain.CustomerImageUsage, error)
}

