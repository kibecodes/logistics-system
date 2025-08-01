package order

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Assigned  OrderStatus = "assigned"
	InTransit OrderStatus = "in_transit"
	Delivered OrderStatus = "delivered"
	Cancelled OrderStatus = "cancelled"
)

type Order struct {
	ID               uuid.UUID   `db:"id" json:"id"`
	Quantity         int         `db:"quantity" json:"quantity"`
	InventoryID      uuid.UUID   `db:"inventory_id" json:"inventory_id"`
	CustomerID       uuid.UUID   `db:"user_id" json:"customer_id"`
	PickupLocation   string      `db:"pickup_address" json:"pickup_location"`
	DeliveryLocation string      `db:"delivery_address" json:"delivery_location"`
	OrderStatus      OrderStatus `db:"status" json:"order_status"`
	CreatedAt        time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time   `db:"updated_at" json:"updated_at"`
}
