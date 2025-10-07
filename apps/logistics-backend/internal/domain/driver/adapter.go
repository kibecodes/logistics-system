package driver

import (
	"context"
	"logistics-backend/internal/domain/notification"
)

type NotificationReader interface {
	Create(ctx context.Context, n *notification.Notification) error
}
