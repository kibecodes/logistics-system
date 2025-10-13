package order

import (
	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	AdminID     uuid.UUID `json:"admin_id" binding:"required"`
	Quantity    int       `json:"quantity" binding:"required"`
	InventoryID uuid.UUID `json:"inventory_id" binding:"required"`
	CustomerID  uuid.UUID `json:"customer_id" binding:"required"`

	PickupAddress string         `json:"pickup_address" binding:"required"`
	PickupPoint   postgis.PointS `json:"pickup_point" binding:"required"`

	DeliveryAddress string         `json:"delivery_address" binding:"required"`
	DeliveryPoint   postgis.PointS `json:"delivery_point" binding:"required"`
}

type UpdateOrderRequest struct {
	Column string      `json:"column" binding:"required"` // e.g. "status", "quantity"
	Value  interface{} `json:"value" binding:"required"`  // Accepts string, int, etc.
}

func (r *CreateOrderRequest) ToOrder() *Order {
	return &Order{
		AdminID:         r.AdminID,
		Quantity:        r.Quantity,
		InventoryID:     r.InventoryID,
		CustomerID:      r.CustomerID,
		PickupAddress:   r.PickupAddress,
		DeliveryAddress: r.DeliveryAddress,
		Status:          Pending,
	}
}

// DropdownDataRequest represents the data used to populate order form dropdowns.
// swagger:model
type DropdownDataRequest struct {
	Customers   []Customer  `json:"customers"`
	Inventories []Inventory `json:"inventories"`
}

type Customer struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"full_name" json:"name"`
}

type Inventory struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	AdminID  uuid.UUID `db:"admin_id" json:"admin_id"`
	Category string    `db:"category" json:"category"`
}
