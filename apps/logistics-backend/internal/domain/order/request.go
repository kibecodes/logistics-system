package order

import (
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	Quantity         int       `json:"quantity" binding:"required"`
	InventoryID      uuid.UUID `json:"inventory_id" binding:"required"`
	CustomerID       uuid.UUID `json:"customer_id" binding:"required"`
	PickupLocation   string    `json:"pickup_location" binding:"required"`
	DeliveryLocation string    `json:"delivery_location" binding:"required"`
}

type UpdateOrderRequest struct {
	Column string      `json:"column" binding:"required"` // e.g. "status", "quantity"
	Value  interface{} `json:"value" binding:"required"`  // Accepts string, int, etc.
}

func (r *CreateOrderRequest) ToOrder() *Order {
	return &Order{
		Quantity:         r.Quantity,
		InventoryID:      r.InventoryID,
		CustomerID:       r.CustomerID,
		PickupLocation:   r.PickupLocation,
		DeliveryLocation: r.DeliveryLocation,
		OrderStatus:      Pending,
	}
}
