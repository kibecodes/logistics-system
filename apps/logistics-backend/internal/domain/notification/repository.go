package notification

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*Notification, error)
	List(ctx context.Context) ([]*Notification, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
