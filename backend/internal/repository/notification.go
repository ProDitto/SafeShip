package repository

import (
	"context"
	"secure-image-service/internal/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
}

