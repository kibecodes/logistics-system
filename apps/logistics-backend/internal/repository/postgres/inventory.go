package postgres

import (
	"context"
	"fmt"

	// "fmt"
	"logistics-backend/internal/domain/inventory"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRespository(db *sqlx.DB) inventory.Repository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(i *inventory.Inventory) error {
	query := `
		INSERT INTO inventories 
		(admin_id, name, category, stock, price_amount, price_currency, images, unit, packaging, description, location, slug)
		VALUES (:admin_id, :name, :category, :stock, :price.amount, :price.currency, :images, :unit, :packaging, :description, :location, :slug)
		RETURNING id
	`
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&i.ID, i)
}

func (r *InventoryRepository) GetByID(InventoryID uuid.UUID) (*inventory.Inventory, error) {
	query := `
		SELECT id, admin_id, name, category, stock, price_amount AS "price.amount", price_currency AS "price.currency", images, unit, packaging, description, location, slug 
		FROM inventories 
		WHERE id = $1
	`
	var inventory inventory.Inventory
	err := r.db.Get(&inventory, query, InventoryID)
	return &inventory, err
}

func (r *InventoryRepository) GetByName(InventoryName string) (*inventory.Inventory, error) {
	query := `
		SELECT id, admin_id, name, category, stock, price_amount AS "price.amount", price_currency AS "price.currency", images, unit, packaging, description, location, slug
		FROM inventories
		WHERE name = $1
	`
	var inventory inventory.Inventory
	err := r.db.Get(&inventory, query, InventoryName)
	if err != nil {
		return nil, err
	}

	return &inventory, nil
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

	query := fmt.Sprintf(`UPDATE inventories SET %s = $1, updated_at = NOW() WHERE id = $2`, column)
	res, err := r.db.ExecContext(ctx, query, value, inventoryID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no order found with id %s", inventoryID)
	}
	return nil
}

func (r *InventoryRepository) GetByCategory(ctx context.Context, category string) ([]inventory.Inventory, error) {
	query := `
	SELECT id, admin_id, name, category, stock,
       price_amount AS "price.amount",
       price_currency AS "price.currency",
       images, unit, packaging, description, location, slug, created_at, updated_at
	FROM inventories
	WHERE category = $1`
	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []inventory.Inventory
	for rows.Next() {
		var inv inventory.Inventory
		if err := rows.Scan(
			&inv.ID, &inv.AdminID, &inv.Name, &inv.Category,
			&inv.Stock, &inv.Price.Amount, &inv.Price.Currency, &inv.Images, &inv.Unit,
			&inv.Packaging, &inv.Description, &inv.Location,
			&inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, err
		}
		inventories = append(inventories, inv)
	}

	return inventories, nil
}

func (r *InventoryRepository) ListCategories(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT category FROM inventories ORDER BY category`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (r *InventoryRepository) List(limit, offset int) ([]*inventory.Inventory, error) {
	query := `
		SELECT id, admin_id, name, category, stock, price_amount AS "price.amount", price_currency AS "price.currency", images, unit, packaging, description, location, slug, created_at, updated_at
		FROM inventories
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	var inventories []*inventory.Inventory
	err := r.db.Select(&inventories, query, limit, offset)
	return inventories, err
}

func (r *InventoryRepository) GetBySlugs(adminSlug, productSlug string) (*inventory.Inventory, error) {
	query := `
		SELECT i.id, i.admin_id, i.name, i.category, i.stock, i.price_amount AS "price.amount", i.price_currency AS "price.currency", i.images, i.unit, i.packaging, i.description, i.location, i.slug
		FROM inventories i
		JOIN users u ON i.admin_id = u.id
		WHERE i.slug = $1 AND u.slug = $2
	`

	var inv inventory.Inventory
	err := r.db.Get(&inv, query, adminSlug, productSlug)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *InventoryRepository) GetStoreView(adminSlug string) (*inventory.StorePublicView, error) {

	// Getting admin info
	var store inventory.StorePublicView
	adminQuery := `
		SELECT full_name AS admin_name, category, location
		FROM users 
		WHERE slug = $1 
		AND role = 'admin'
	`
	if err := r.db.Get(&store, adminQuery, adminQuery); err != nil {
		return nil, err
	}

	// Getting products for this admin
	productQuery := `
		SELECT name, price_amount AS "price.amount", price_currency AS "price.currency", unit, packaging, stock AS in_stock,
			(split_part(images, ',', 1)) AS image
		FROM inventories i
		JOIN users u ON i.admin_id = u.id
		WHERE u.slug = $1
	`
	err := r.db.Select(&store.Products, productQuery, adminQuery)
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (r *InventoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM inventories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *InventoryRepository) GetAllInventories(ctx context.Context) ([]inventory.AllInventory, error) {
	query := `
        SELECT id, name
        FROM inventories
        ORDER BY name ASC
    `
	var inventories []inventory.AllInventory
	err := r.db.Select(&inventories, query)
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
