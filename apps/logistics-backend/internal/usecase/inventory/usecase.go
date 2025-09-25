package inventory

import (
	"context"
	"fmt"
	domain "logistics-backend/internal/domain/inventory"
	"logistics-backend/internal/usecase/common"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      domain.Repository
	txManager common.TxManager
}

func NewUseCase(repo domain.Repository, txm common.TxManager) *UseCase {
	return &UseCase{repo: repo, txManager: txm}
}

func (uc *UseCase) CreateInventory(ctx context.Context, i *domain.Inventory) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Create(txCtx, i); err != nil {
			return fmt.Errorf("could not create inventory: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetInventory(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetInventoryByName(ctx context.Context, name string) (*domain.Inventory, error) {
	return uc.repo.GetByName(ctx, name)
}

func (uc *UseCase) UpdateInventory(ctx context.Context, inventoryId uuid.UUID, column string, value any) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateColumn(txCtx, inventoryId, column, value); err != nil {
			return fmt.Errorf("update inventory failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) List(ctx context.Context, limit, offset int) ([]*domain.Inventory, error) {
	return uc.repo.List(ctx, limit, offset)
}

func (uc *UseCase) GetByCategory(ctx context.Context, category string) ([]domain.Inventory, error) {
	return uc.repo.GetByCategory(ctx, category)
}

func (uc *UseCase) ListCategories(ctx context.Context) ([]string, error) {
	return uc.repo.ListCategories(ctx)
}

func (uc *UseCase) GetBySlugs(ctx context.Context, adminSlug, productSlug string) (*domain.Inventory, error) {
	return uc.repo.GetBySlugs(ctx, adminSlug, productSlug)
}

func (uc *UseCase) GetStorePublicView(ctx context.Context, adminSlug string) (*domain.StorePublicView, error) {
	return uc.repo.GetStoreView(ctx, adminSlug)
}

func (uc *UseCase) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete order failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetAllInventories(ctx context.Context) ([]domain.AllInventory, error) {
	return uc.repo.GetAllInventories(ctx)
}
