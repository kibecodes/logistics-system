package order

import (
	"time"

	"github.com/cridenour/go-postgis"
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
	ID          uuid.UUID `db:"id" json:"id"`
	AdminID     uuid.UUID `db:"admin_id" json:"admin_id"`
	Quantity    int       `db:"quantity" json:"quantity"`
	InventoryID uuid.UUID `db:"inventory_id" json:"inventory_id"`
	CustomerID  uuid.UUID `db:"user_id" json:"customer_id"`

	PickupAddress string         `db:"pickup_address" json:"pickup_address"`
	PickupPoint   postgis.PointS `db:"pickup_point" json:"pickup_point"`

	DeliveryAddress string         `db:"delivery_address" json:"delivery_address"`
	DeliveryPoint   postgis.PointS `db:"delivery_point" json:"delivery_point"`

	Status    OrderStatus `db:"status" json:"status"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
}
