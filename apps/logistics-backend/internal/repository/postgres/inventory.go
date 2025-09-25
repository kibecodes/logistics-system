package postgres

import (
	"context"
	"fmt"

	// "fmt"

	"logistics-backend/internal/application"
	"logistics-backend/internal/domain/inventory"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InventoryRepository struct {
	exec sqlx.ExtContext
}

func NewInventoryRespository(db *sqlx.DB) *InventoryRepository {
	return &InventoryRepository{exec: db}
}

func (r *InventoryRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *InventoryRepository) Create(ctx context.Context, i *inventory.Inventory) error {
	query := `
		INSERT INTO inventories 
		(admin_id, name, category, stock, price_amount, price_currency, images, unit, packaging, description, location, slug)
		VALUES (:admin_id, :name, :category, :stock, :price_amount, :price_currency, :images, :unit, :packaging, :description, :location, :slug)
		RETURNING id
	`
	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, i)
	if err != nil {
		return fmt.Errorf("insert inventory: %w", err)
	}

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&i.ID); err != nil {
			return fmt.Errorf("scanning new order id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after scan")
	}

	return nil
}

func (r *InventoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*inventory.Inventory, error) {
	query := `
		SELECT id, admin_id, name, category, stock, price_amount, price_currency, images, unit, packaging, description, location, slug 
		FROM inventories 
		WHERE id = $1
	`
	var i inventory.Inventory
	if err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &i, query, id); err != nil {
		return nil, fmt.Errorf("get inventory by id: %w", err)
	}

	return &i, nil

}

func (r *InventoryRepository) GetByName(ctx context.Context, name string) (*inventory.Inventory, error) {
	query := `
		SELECT id, admin_id, name, category, stock, price_amount, price_currency, images, unit, packaging, description, location, slug
		FROM inventories
		WHERE name = $1
	`
	var i inventory.Inventory
	if err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &i, query, name); err != nil {
		return nil, fmt.Errorf("get inventory by name: %w", err)
	}

	return &i, nil
}

func (r *InventoryRepository) UpdateColumn(ctx context.Context, inventoryID uuid.UUID, column string, value any) error {
	// Whitelist column names
	allowed := map[string]bool{
		"name":           true,
		"stock":          true,
		"price_amount":   true,
		"price_currency": true,
		"unit":           true,
		"location":       true,
	}

	if !allowed[column] {
		return fmt.Errorf("attempted to update disallowed column: %s", column)
	}

	query := fmt.Sprintf(`
		UPDATE inventories SET %s = :value, updated_at = NOW() 
		WHERE id = :id
	`, column)

	args := map[string]interface{}{
		"value": value,
		"id":    inventoryID,
	}

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, args)
	if err != nil {
		return fmt.Errorf("update inventory: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no inventory found with id %s", inventoryID)
	}

	return nil
}

func (r *InventoryRepository) GetByCategory(ctx context.Context, category string) ([]inventory.Inventory, error) {
	query := `
		SELECT id, admin_id, name, category, stock, price_amount, price_currency, images, unit, packaging, description, location, slug, created_at, updated_at
		FROM inventories
		WHERE category = $1
	`

	var inventories []inventory.Inventory
	if err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &inventories, query, category); err != nil {
		return nil, fmt.Errorf("get inventories by category: %w", err)
	}

	return inventories, nil

}

func (r *InventoryRepository) ListCategories(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT category 
		FROM inventories 
		ORDER BY category
	`
	var categories []string
	if err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &categories, query); err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	return categories, nil
}

func (r *InventoryRepository) List(ctx context.Context, limit, offset int) ([]*inventory.Inventory, error) {
	query := `
		SELECT id, admin_id, name, category, stock, price_amount, price_currency, images, unit, packaging, description, location, slug, created_at, updated_at
		FROM inventories
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	var inventories []*inventory.Inventory
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &inventories, query, limit, offset)
	return inventories, err
}

func (r *InventoryRepository) GetBySlugs(ctx context.Context, adminSlug, productSlug string) (*inventory.Inventory, error) {
	query := `
		SELECT i.id, i.admin_id, i.name, i.category, i.stock, i.price_amount, i.price_currency, i.images, i.unit, i.packaging, i.description, i.location, i.slug
		FROM inventories i
		JOIN users u ON i.admin_id = u.id
		WHERE i.slug = $1 AND u.slug = $2
	`

	var inv inventory.Inventory
	if err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &inv, query, productSlug, adminSlug); err != nil {
		return nil, fmt.Errorf("get inventory by slugs: %w", err)
	}

	return &inv, nil
}

func (r *InventoryRepository) GetStoreView(ctx context.Context, adminSlug string) (*inventory.StorePublicView, error) {

	// Getting admin info
	adminQuery := `
	SELECT full_name AS admin_name, category, location
	FROM users 
	WHERE slug = $1 
	AND role = 'admin'
	`
	var store inventory.StorePublicView
	if err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &store, adminQuery, adminSlug); err != nil {
		return nil, fmt.Errorf("get store admin info: %w", err)
	}

	// Getting products for this admin
	productQuery := `
		SELECT name, price_amount, price_currency, unit, packaging, stock AS in_stock,
			(split_part(images, ',', 1)) AS image
		FROM inventories i
		JOIN users u ON i.admin_id = u.id
		WHERE u.slug = $1
	`
	if err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &store.Products, productQuery, adminSlug); err != nil {
		return nil, fmt.Errorf("get store products: %w", err)
	}

	return &store, nil
}

func (r *InventoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM inventories 
		WHERE id = $1
	`

	res, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete inventory: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify inventory deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("inventory already deleted or invalid")
	}

	return nil
}

func (r *InventoryRepository) GetAllInventories(ctx context.Context) ([]inventory.AllInventory, error) {
	query := `
        SELECT id, name, admin_id, category
        FROM inventories
        ORDER BY name ASC
    `
	var inventories []inventory.AllInventory
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &inventories, query)
	return inventories, err
}

// func GetColumnValues[T any](ctx context.Context, db *sqlx.DB, column string) ([]T, error) {
// 	// Optional: whitelist
// 	allowed := map[string]bool{
// 		"stock":    true,
// 		"category": true,
// 		"location": true,
// 	}

// 	if !allowed[column] {
// 		return nil, fmt.Errorf("invalid column: %s", column)
// 	}

// 	query := fmt.Sprintf(`SELECT %s FROM inventories`, column)

// 	var values []T
// 	err := db.SelectContext(ctx, &values, query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return values, nil
// }

// func UpdateColumnValues[T any](ctx context.Context, db *sqlx.DB, column string, newValue T, filterColumn string, filterValue any) error {
// 	allowed := map[string]bool{
// 		"stock":    true,
// 		"category": true,
// 		"location": true,
// 		"name":     true,
// 	}

// 	if !allowed[column] {
// 		return fmt.Errorf("invalid column to update: %s", column)
// 	}
// 	if !allowed[filterColumn] {
// 		return fmt.Errorf("invalid filter column: %s", filterColumn)
// 	}

// 	query := fmt.Sprintf(`UPDATE inventories SET %s = $1 WHERE %s = $2`, column, filterColumn)

// 	_, err := db.ExecContext(ctx, query, newValue, filterValue)
// 	return err
// }
