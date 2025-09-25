package inventoryadapter

import (
	"context"
	"logistics-backend/internal/domain/inventory"
	"logistics-backend/internal/domain/order"
	inventoryusecase "logistics-backend/internal/usecase/inventory"

	"github.com/google/uuid"
)

type UseCaseAdapter struct {
	UseCase *inventoryusecase.UseCase
}

func (a *UseCaseAdapter) GetInventoryByID(ctx context.Context, id uuid.UUID) (*inventory.Inventory, error) {
	return a.UseCase.GetInventory(ctx, id)
}

func (a *UseCaseAdapter) UpdateInventory(ctx context.Context, inventoryId uuid.UUID, column string, value any) error {
	return a.UseCase.UpdateInventory(ctx, inventoryId, column, value)
}

func (a *UseCaseAdapter) GetAllInventories(ctx context.Context) ([]order.Inventory, error) {
	invs, err := a.UseCase.GetAllInventories(ctx) // returns []inventory.AllInventory
	if err != nil {
		return nil, err
	}

	// Map to []order.Inventory
	res := make([]order.Inventory, len(invs))
	for i, inv := range invs {
		res[i] = order.Inventory{
			ID:       inv.ID,
			Name:     inv.Name,
			AdminID:  inv.AdminID,
			Category: inv.Category,
			// map other fields you need
		}
	}
	return res, nil
}
