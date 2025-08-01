package delivery

import (
	"context"
	"logistics-backend/internal/domain/driver"
	"logistics-backend/internal/domain/order"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// cross-domain interface so delivery can access orders

type OrderReader interface {
	// order
	GetOrderByID(ctx context.Context, id uuid.UUID) (*order.Order, error)
	UpdateOrderTx(ctx context.Context, tx *sqlx.Tx, orderID uuid.UUID, column string, value any) error

	// driver
	GetDriverByID(ctx context.Context, id uuid.UUID) (*driver.Driver, error)
	UpdateDriverAvailability(ctx context.Context, driverID uuid.UUID, column string, value bool) error
}
