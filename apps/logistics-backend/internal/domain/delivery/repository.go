package delivery

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(delivery *Delivery) error                              // POST method to create delivery from orders.
	CreateTx(ctx context.Context, tx *sqlx.Tx, d *Delivery) error // POST method for create delivery + update orders operations
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	GetByID(id uuid.UUID) (*Delivery, error)                                          // GET method for fetching delivery by id
	List() ([]*Delivery, error)                                                       // GET method to fetch all deliveries
	Update(ctx context.Context, deliveryID uuid.UUID, column string, value any) error // PUT generic method to update specified column value in orders table
	Accept(ctx context.Context, d *Delivery) error                                    // PATCH method for driver to accept delivery.
	Delete(ctx context.Context, id uuid.UUID) error                                   // DELETE method to remove delivery by ID
}
