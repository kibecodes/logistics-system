package inventory

import (
	generate "logistics-backend/internal/utils"

	"github.com/google/uuid"
)

type CreateInventoryRequest struct {
	AdminID       uuid.UUID `json:"admin_id"`                    // Foreign key
	Name          string    `json:"name" binding:"required"`     // e.g. “Fresh Milk”
	Category      string    `json:"category" binding:"required"` // e.g. “Dairy”
	Stock         int       `json:"stock" binding:"required"`
	PriceAmount   int64     `json:"price_amount" binding:"required"`
	PriceCurrency string    `json:"price_currency" binding:"required"`
	Images        string    `json:"images" binding:"required"`    // could be JSON array or URLs
	Unit          string    `json:"unit" binding:"required"`      // "per litre", "per bucket"
	Packaging     string    `json:"packaging" binding:"required"` // “Bucket/Single”
	Description   string    `json:"description" binding:"required"`
	Location      string    `json:"location" binding:"required"` // optional
	Slug          string    `json:"slug" binding:"required"`
}

func (r *CreateInventoryRequest) ToInventory() *Inventory {
	return &Inventory{
		AdminID:       r.AdminID,
		Name:          r.Name,
		Category:      r.Category,
		Stock:         r.Stock,
		PriceAmount:   r.PriceAmount,
		PriceCurrency: r.PriceCurrency,
		Images:        r.Images,
		Unit:          r.Unit,
		Packaging:     r.Packaging,
		Description:   r.Description,
		Location:      r.Location,
		Slug:          generate.GenerateSlug(r.Name),
	}
}
