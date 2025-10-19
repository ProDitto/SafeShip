package postgres

import (
	"context"
	"errors"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgCustomerRepository struct {
	db *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) repository.CustomerRepository {
	return &pgCustomerRepository{db: db}
}

func (r *pgCustomerRepository) FindAll(ctx context.Context) ([]*domain.Customer, error) {
	rows, err := r.db.Query(ctx, "SELECT namespace, name, contact_info, sla_tier, created_at, updated_at FROM customers ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[domain.Customer])
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *pgCustomerRepository) FindByNamespace(ctx context.Context, namespace string) (*domain.Customer, error) {
	row := r.db.QueryRow(ctx, "SELECT namespace, name, contact_info, sla_tier, created_at, updated_at FROM customers WHERE namespace = $1", namespace)

	customer, err := pgx.RowToAddrOfStructByPos[domain.Customer](row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}

	return customer, nil
}

