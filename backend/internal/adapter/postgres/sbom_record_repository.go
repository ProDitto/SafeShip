package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"
)

type pgSBOMRecordRepository struct {
	db *pgxpool.Pool
}

func NewSBOMRecordRepository(db *pgxpool.Pool) repository.SBOMRecordRepository {
	return &pgSBOMRecordRepository{db: db}
}

func (r *pgSBOMRecordRepository) Create(ctx context.Context, sbom *domain.SBOMRecord) error {
	query := `INSERT INTO sbom_records (image_id, format, uri) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, sbom.ImageID, sbom.Format, sbom.URI)
	return err
}

