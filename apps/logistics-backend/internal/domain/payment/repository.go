package payment

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByID(cxt context.Context, id uuid.UUID) (*Payment, error)
	GetByOrder(ctx context.Context, id uuid.UUID) (*Payment, error)
	List(ctx context.Context) ([]*Payment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
