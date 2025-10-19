package repository

import (
	"context"
	"secure-image-service/backend/internal/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
}

