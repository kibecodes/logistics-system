package inventory

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, inventory *Inventory) error            // POST method for creating new inventory.
	GetByID(ctx context.Context, id uuid.UUID) (*Inventory, error)     // GET method for fetching inventory by id.
	GetByName(ctx context.Context, name string) (*Inventory, error)    // GET method for fetching inventory by name.
	List(ctx context.Context, limit, offset int) ([]*Inventory, error) // GET method for fetching all inventories - slice.
	GetAllInventories(ctx context.Context) ([]AllInventory, error)     // GET all inv ID, Name & AdminID without pagination.
	Delete(ctx context.Context, id uuid.UUID) error                    // DELETE method to remove inventory by id.

	GetByCategory(ctx context.Context, category string) ([]*Inventory, error)       // GET method for fetching inventories(slice) by category.
	ListCategories(ctx context.Context) ([]string, error)                           // GET method for fetching all categories in inventories table.
	UpdateColumn(ctx context.Context, id uuid.UUID, column string, value any) error // PUT method for updating table column values.

	GetByStoreID(ctx context.Context, storeID uuid.UUID) ([]*Inventory, error)
}
