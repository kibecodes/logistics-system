package order

import (
	"context"
	"logistics-backend/internal/domain/inventory"

	"github.com/google/uuid"
)

// cross-domain DI using necessary interface

type InventoryReader interface {
	// inventory
	GetInventoryByID(ctx context.Context, id uuid.UUID) (*inventory.Inventory, error)
	UpdateInventory(ctx context.Context, inventoryId uuid.UUID, column string, value any) error
}
