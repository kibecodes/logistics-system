package order

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, order *Order) error                                // POST method to create order.
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)                     // GET method for fetching order by id.
	ListByCustomer(ctx context.Context, customerID uuid.UUID) ([]*Order, error)    // GET method for fetching all orders by customer id.
	Update(ctx context.Context, orderID uuid.UUID, column string, value any) error // PATCH method to update specified column value in orders table.
	List(ctx context.Context) ([]*Order, error)                                    // GET method for fetching all orders
	Delete(ctx context.Context, id uuid.UUID) error                                // DELETE method for removing order by id
}
