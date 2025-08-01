package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(order *Order) error                                                                          // POST method to create order.
	GetByID(id uuid.UUID) (*Order, error)                                                               // GET method for fetching order by id.
	ListByCustomer(customerID uuid.UUID) ([]*Order, error)                                              // GET method for fetching all orders by customer id.
	UpdateColumn(ctx context.Context, orderID uuid.UUID, column string, value any) error                // PATCH method to update specified column value in orders table.
	UpdateColumnTx(ctx context.Context, tx *sqlx.Tx, orderID uuid.UUID, column string, value any) error // PUT transaction method for multipe operation - create + update
	List() ([]*Order, error)                                                                            // GET method for fetching all orders
	Delete(ctx context.Context, id uuid.UUID) error                                                     // DELETE method for removing order by id
}
