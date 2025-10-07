package notification

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string
type NotificationStatus string

const (
	Email  NotificationType = "email"
	SMS    NotificationType = "sms"
	Push   NotificationType = "push"
	System NotificationType = "system" // for in-app or placeholder

	Pending NotificationStatus = "pending"
	Sent    NotificationStatus = "sent"
	Failed  NotificationStatus = "failed"
)

type Notification struct {
	ID        uuid.UUID          `db:"id" json:"id"`
	UserID    uuid.UUID          `db:"user_id" json:"user_id"`
	Message   string             `db:"message" json:"message"`
	Type      NotificationType   `db:"type" json:"type"`
	Status    NotificationStatus `db:"status" json:"status"`
	SentAt    time.Time          `db:"sent_at" json:"sent_at"`
	UpdatedAt time.Time          `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time          `db:"created_at" json:"created_at"`
}
