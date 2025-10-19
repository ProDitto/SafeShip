package usecase

import (
	"context"
	"secure-image-service/internal/domain"
	"secure-image-service/internal/repository"
)

type CustomerUsecase struct {
	repo repository.CustomerRepository
}

func NewCustomerUsecase(repo repository.CustomerRepository) *CustomerUsecase {
	return &CustomerUsecase{repo: repo}
}

func (uc *CustomerUsecase) ListCustomers(ctx context.Context) ([]*domain.Customer, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *CustomerUsecase) GetCustomer(ctx context.Context, namespace string) (*domain.Customer, error) {
	return uc.repo.FindByNamespace(ctx, namespace)
}

