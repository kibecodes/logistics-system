package inventoryadapter

import (
	"context"
	"logistics-backend/internal/domain/inventory"
	inventoryusecase "logistics-backend/internal/usecase/inventory"

	"github.com/google/uuid"
)

type InventoryUseCaseAdapter struct {
	UseCase *inventoryusecase.UseCase
}

func (a *InventoryUseCaseAdapter) GetInventoryByID(ctx context.Context, id uuid.UUID) (*inventory.Inventory, error) {
	return a.UseCase.GetByID(ctx, id)
}

func (a *InventoryUseCaseAdapter) UpdateInventory(ctx context.Context, inventoryId uuid.UUID, column string, value any) error {
	return a.UseCase.UpdateInventory(ctx, inventoryId, column, value)
}
