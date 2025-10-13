package inventory

import (
	"context"
	"logistics-backend/internal/domain/notification"
	"logistics-backend/internal/domain/store"

	"github.com/google/uuid"
)

type NotificationReader interface {
	Create(ctx context.Context, n *notification.Notification) error
}

type StoreReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (*store.Store, error)
}
