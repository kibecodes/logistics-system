package delivery

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, delivery *Delivery) error                             // POST method to create delivery from orders.
	GetByID(ctx context.Context, id uuid.UUID) (*Delivery, error)                     // GET method for fetching delivery by id
	List(ctx context.Context) ([]*Delivery, error)                                    // GET method to fetch all deliveries
	Update(ctx context.Context, deliveryID uuid.UUID, column string, value any) error // PUT generic method to update specified column value in orders table
	Accept(ctx context.Context, d *Delivery) error                                    // PATCH method for driver to accept delivery.
	Delete(ctx context.Context, id uuid.UUID) error                                   // DELETE method to remove delivery by ID

	ListByStatus(ctx context.Context, statuses []DeliveryStatus) ([]*Delivery, error)
}
