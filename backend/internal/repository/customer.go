package repository

import (
	"context"
	"secure-image-service/internal/domain"
)

type CustomerRepository interface {
	FindAll(ctx context.Context) ([]*domain.Customer, error)
	FindByNamespace(ctx context.Context, namespace string) (*domain.Customer, error)
}

