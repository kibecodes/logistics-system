package inventory

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(inventory *Inventory) error                             // POST method for creating new inventory.
	GetByID(id uuid.UUID) (*Inventory, error)                      // GET method for fetching inventory by id.
	GetByName(name string) (*Inventory, error)                     // GET method for fetching inventory by name.
	List(limit, offset int) ([]*Inventory, error)                  // GET method for fetching all inventories - slice.
	GetAllInventories(ctx context.Context) ([]AllInventory, error) // GET all inv name & ID without pagination.
	Delete(ctx context.Context, id uuid.UUID) error                // DELETE method to remove inventory by id.

	GetByCategory(ctx context.Context, category string) ([]Inventory, error)                 // GET method for fetching inventories(slice) by category.
	ListCategories(ctx context.Context) ([]string, error)                                    // GET method for fetching all categories in inventories table.
	UpdateColumn(ctx context.Context, inventoryID uuid.UUID, column string, value any) error // PUT method for updating table column values.

	GetBySlugs(adminSlugs, productSlug string) (*Inventory, error) // GET method for fetching specified inventory products by admin's slug.
	GetStoreView(adminSlug string) (*StorePublicView, error)       // GET method for fetching admin's store view.
}
