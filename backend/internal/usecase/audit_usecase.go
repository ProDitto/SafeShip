package usecase

import (
	"context"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"
)

type AuditUsecase struct {
	repo repository.AuditLogRepository
}

func NewAuditUsecase(repo repository.AuditLogRepository) *AuditUsecase {
	return &AuditUsecase{repo: repo}
}

func (uc *AuditUsecase) Log(ctx context.Context, namespace, action, actor string, details map[string]interface{}) error {
	logEntry := &domain.AuditLog{
		TenantNamespace: namespace,
		Action:          action,
		Actor:           actor,
		Details:         details,
	}
	return uc.repo.Create(ctx, logEntry)
}
