package postgres

import (
	"context"
	"errors"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgImageRepository struct {
	db *pgxpool.Pool
}

func NewImageRepository(db *pgxpool.Pool) repository.ImageRepository {
	return &pgImageRepository{db: db}
}

func (r *pgImageRepository) FindAll(ctx context.Context) ([]*domain.Image, error) {
	rows, err := r.db.Query(ctx, "SELECT id, tenant_namespace, digest, tags, slsa_level, created_at, updated_at FROM images ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[domain.Image])
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (r *pgImageRepository) FindByID(ctx context.Context, id int) (*domain.Image, error) {
	row := r.db.QueryRow(ctx, "SELECT id, tenant_namespace, digest, tags, slsa_level, created_at, updated_at FROM images WHERE id = $1", id)

	image, err := pgx.RowToAddrOfStructByPos[domain.Image](row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}

	return image, nil
}

