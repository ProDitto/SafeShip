package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"
)

type pgCVEFindingRepository struct {
	db *pgxpool.Pool
}

func NewCVEFindingRepository(db *pgxpool.Pool) repository.CVEFindingRepository {
	return &pgCVEFindingRepository{db: db}
}

func (r *pgCVEFindingRepository) CreateBatch(ctx context.Context, cves []*domain.CVEFinding) error {
	if len(cves) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `INSERT INTO cve_findings (image_id, cve_id, severity, description, fix_available)
              VALUES ($1, $2, $3, $4, $5)`

	for _, cve := range cves {
		batch.Queue(query, cve.ImageID, cve.CVEID, cve.Severity, cve.Description, cve.FixAvailable)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(cves); i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}

