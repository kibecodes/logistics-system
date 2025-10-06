package order

import (
	"context"
	"logistics-backend/internal/domain/driver"
	"logistics-backend/internal/domain/inventory"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

// cross-domain DI using necessary interface

// Access the inventory domain usecase methods.
type InventoryReader interface {
	GetInventoryByID(ctx context.Context, id uuid.UUID) (*inventory.Inventory, error)
	UpdateInventory(ctx context.Context, inventoryId uuid.UUID, column string, value any) error
	GetAllInventories(ctx context.Context) ([]Inventory, error)
}

// Access the user domain usecase method for getting users of role customers.
type CustomerReader interface {
	GetAllCustomers(ctx context.Context) ([]Customer, error)
}

type DriverReader interface {
	GetNearestDriver(ctx context.Context, pickup postgis.PointS, maxDistance float64) (*driver.Driver, error)
}
