package postgres

import (
	"context"
	"encoding/json"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgAuditLogRepository struct {
	db *pgxpool.Pool
}

func NewAuditLogRepository(db *pgxpool.Pool) repository.AuditLogRepository {
	return &pgAuditLogRepository{db: db}
}

func (r *pgAuditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	detailsJSON, err := json.Marshal(log.Details)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO audit_logs (tenant_namespace, action, actor, details)
        VALUES ($1, $2, $3, $4)
    `
	_, err = r.db.Exec(ctx, query, log.TenantNamespace, log.Action, log.Actor, detailsJSON)
	return err
}
