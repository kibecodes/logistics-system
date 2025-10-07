package inventory

import (
	"context"
	"fmt"
	domain "logistics-backend/internal/domain/inventory"
	"logistics-backend/internal/domain/notification"
	"logistics-backend/internal/usecase/common"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      domain.Repository
	txManager common.TxManager
	notfRepo  domain.NotificationReader
}

func NewUseCase(repo domain.Repository, txm common.TxManager, notf domain.NotificationReader) *UseCase {
	return &UseCase{repo: repo, txManager: txm, notfRepo: notf}
}

func (uc *UseCase) CreateInventory(ctx context.Context, i *domain.Inventory) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Create(txCtx, i); err != nil {
			return fmt.Errorf("could not create inventory: %w", err)
		}

		// After successful creation, fire notification (async)
		go func() {
			msg := fmt.Sprintf("✅ New inventory '%s' has been added with stock %d.", i.Name, i.Stock)
			_ = uc.notify(ctx, i.AdminID, msg)

			// Optional: immediately alert if created with low stock
			if i.Stock <= 5 {
				lowMsg := fmt.Sprintf("⚠️ Inventory '%s' was created with low stock (%d).", i.Name, i.Stock)
				_ = uc.notify(ctx, i.AdminID, lowMsg)
			}
		}()

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
		// 1. Fetch inventory to get AdminID (inside transaction)
		inv, err := uc.repo.GetByID(txCtx, inventoryId)
		if err != nil {
			return fmt.Errorf("could not fetch inventory: %w", err)
		}

		// 2. Update column
		if err := uc.repo.UpdateColumn(txCtx, inventoryId, column, value); err != nil {
			return fmt.Errorf("update inventory failed: %w", err)
		}

		// 3. Fire notification async (after commit)
		go func() {
			msg := fmt.Sprintf("ℹ️ Inventory %s updated: column '%s' changed.", inv.Name, column)
			_ = uc.notify(ctx, inv.AdminID, msg) // you can use AdminID if available
		}()

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
		// 1. Fetch inventory to get AdminID and Name
		inv, err := uc.repo.GetByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("could not fetch inventory: %w", err)
		}

		// 2. Delete
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete inventory failed: %w", err)
		}

		// 3. Fire notification async
		go func() {
			msg := fmt.Sprintf("🗑️ Inventory '%s' has been deleted.", inv.Name)
			_ = uc.notify(ctx, inv.AdminID, msg)
		}()

		return nil
	})
}

func (uc *UseCase) GetAllInventories(ctx context.Context) ([]domain.AllInventory, error) {
	return uc.repo.GetAllInventories(ctx)
}

func (uc *UseCase) notify(ctx context.Context, userID uuid.UUID, message string) error {
	n := &notification.Notification{
		UserID:  userID,
		Message: message,
		Type:    notification.System,
		Status:  notification.Pending,
	}
	return uc.notfRepo.Create(ctx, n)
}
