package postgres

import (
	"context"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type pgImageRepository struct {
	db *pgxpool.Pool
}

// NewImageRepository creates a new PostgreSQL-backed ImageRepository.
func NewImageRepository(db *pgxpool.Pool) repository.ImageRepository {
	return &pgImageRepository{db: db}
}

// FindAll retrieves all images from the database.
func (r *pgImageRepository) FindAll(ctx context.Context) ([]*domain.Image, error) {
	query := `SELECT id, tenant_namespace, digest, tags, slsa_level, created_at, updated_at FROM images ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*domain.Image
	for rows.Next() {
		var image domain.Image
		var tags pq.StringArray
		if err := rows.Scan(&image.ID, &image.TenantNamespace, &image.Digest, &tags, &image.SLSALevel, &image.CreatedAt, &image.UpdatedAt); err != nil {
			return nil, err
		}
		image.Tags = tags
		images = append(images, &image)
	}

	return images, nil
}

// FindByID retrieves a single image by its ID.
func (r *pgImageRepository) FindByID(ctx context.Context, id int) (*domain.Image, error) {
	query := `SELECT id, tenant_namespace, digest, tags, slsa_level, created_at, updated_at FROM images WHERE id = $1`
	var image domain.Image
	var tags pq.StringArray
	err := r.db.QueryRow(ctx, query, id).Scan(&image.ID, &image.TenantNamespace, &image.Digest, &tags, &image.SLSALevel, &image.CreatedAt, &image.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	image.Tags = tags
	return &image, nil
}

func (r *pgImageRepository) Create(ctx context.Context, image *domain.Image) (int, error) {
	query := `INSERT INTO images (tenant_namespace, digest, tags, slsa_level)
              VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, image.TenantNamespace, image.Digest, pq.Array(image.Tags), image.SLSALevel).Scan(&image.ID, &image.CreatedAt, &image.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return image.ID, nil
}
