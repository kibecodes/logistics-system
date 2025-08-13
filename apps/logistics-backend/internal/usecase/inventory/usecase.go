package inventory

import (
	"context"
	domain "logistics-backend/internal/domain/inventory"

	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) CreateInventory(ctx context.Context, i *domain.Inventory) error {
	return uc.repo.Create(i)
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	items, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (uc *UseCase) GetByName(ctx context.Context, name string) (*domain.Inventory, error) {
	return uc.repo.GetByName(name)
}

func (uc *UseCase) UpdateInventory(ctx context.Context, inventoryId uuid.UUID, column string, value any) error {
	return uc.repo.UpdateColumn(ctx, inventoryId, column, value)
}

func (uc *UseCase) List(ctx context.Context, limit, offset int) ([]*domain.Inventory, error) {
	return uc.repo.List(limit, offset)
}

func (uc *UseCase) GetByCategory(ctx context.Context, category string) ([]domain.Inventory, error) {
	return uc.repo.GetByCategory(ctx, category)
}

func (uc *UseCase) ListCategories(ctx context.Context) ([]string, error) {
	return uc.repo.ListCategories(ctx)
}

func (uc *UseCase) GetBySlugs(ctx context.Context, adminSlug, productSlug string) (*domain.Inventory, error) {
	return uc.repo.GetBySlugs(adminSlug, productSlug)
}

func (uc *UseCase) GetStorePublicView(ctx context.Context, adminSlug string) (*domain.StorePublicView, error) {
	return uc.repo.GetStoreView(adminSlug)
}

func (uc *UseCase) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) GetAllInventories(ctx context.Context) ([]domain.AllInventory, error) {
	return uc.repo.GetAllInventories(ctx)
}
