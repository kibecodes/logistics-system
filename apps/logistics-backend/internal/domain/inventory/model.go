package inventory

import (
	"time"

	"logistics-backend/internal/domain/money"

	"github.com/google/uuid"
)

type Inventory struct {
	ID            uuid.UUID `db:"id" json:"id"`
	StoreID       uuid.UUID `db:"store_id" json:"store_id"` // FK to stores.id
	Category      string    `db:"category" json:"category"` // e.g. “Dairy” - "name" field equivalent
	Stock         int       `db:"stock" json:"stock"`
	PriceAmount   int64     `db:"price_amount" json:"price_amount"`
	PriceCurrency string    `db:"price_currency" json:"price_currency"`
	Images        string    `db:"images" json:"images"`       // could be JSON array or URLs
	Unit          string    `db:"unit" json:"unit"`           // "per litre", "per bucket"
	Packaging     string    `db:"packaging" json:"packaging"` // “Bucket/Single”
	Description   string    `db:"description" json:"description"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type StorePublicView struct {
	AdminName string             `json:"admin_name"`
	Category  string             `json:"category"`
	Location  string             `json:"location"`
	Products  []InventorySummary `json:"products"`
}

type InventorySummary struct {
	Name      string      `json:"name"`
	Price     money.Money `json:"price"`
	Image     string      `json:"image"`
	Unit      string      `json:"unit"`
	Packaging string      `json:"packaging"`
	InStock   int         `json:"in_stock"`
}

type AllInventory struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	AdminID  uuid.UUID `db:"admin_id" json:"admin_id"`
	Category string    `db:"category" json:"category"`
}
