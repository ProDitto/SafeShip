package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"
)

type pgBuildEventRepository struct {
	db *pgxpool.Pool
}

func NewBuildEventRepository(db *pgxpool.Pool) repository.BuildEventRepository {
	return &pgBuildEventRepository{db: db}
}

func (r *pgBuildEventRepository) Create(ctx context.Context, event *domain.BuildEvent) (int, error) {
	query := `INSERT INTO build_events (tenant_namespace, trigger_type, status)
              VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, event.TenantNamespace, event.TriggerType, event.Status).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return event.ID, nil
}

func (r *pgBuildEventRepository) FindByID(ctx context.Context, id int) (*domain.BuildEvent, error) {
	query := `SELECT id, tenant_namespace, image_id, trigger_type, status, created_at, updated_at
              FROM build_events WHERE id = $1`
	event := &domain.BuildEvent{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&event.ID, &event.TenantNamespace, &event.ImageID, &event.TriggerType,
		&event.Status, &event.CreatedAt, &event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *pgBuildEventRepository) Update(ctx context.Context, event *domain.BuildEvent) error {
	query := `UPDATE build_events SET image_id = $1, status = $2, updated_at = NOW()
              WHERE id = $3`
	_, err := r.db.Exec(ctx, query, event.ImageID, event.Status, event.ID)
	return err
}

