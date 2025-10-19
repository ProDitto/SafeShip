package postgres

import (
	"context"
	"secure-image-service/backend/internal/domain"
	"secure-image-service/backend/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgSLAViolationRepository struct {
	db *pgxpool.Pool
}

func NewSLAViolationRepository(db *pgxpool.Pool) repository.SLAViolationRepository {
	return &pgSLAViolationRepository{db: db}
}

func (r *pgSLAViolationRepository) Create(ctx context.Context, violation *domain.SLAViolation) error {
	query := `
        INSERT INTO sla_violations (tenant_namespace, cve_finding_id, status)
        VALUES ($1, $2, $3)
    `
	_, err := r.db.Exec(ctx, query, violation.TenantNamespace, violation.CVEFindingID, violation.Status)
	return err
}

