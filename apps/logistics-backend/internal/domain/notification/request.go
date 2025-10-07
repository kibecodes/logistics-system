package notification

import (
	"time"

	"github.com/google/uuid"
)

type CreateNotificationRequest struct {
	UserID  uuid.UUID        `json:"user_id"`
	Message string           `json:"message"`
	Type    NotificationType `json:"type"`
}

type UpdateNotificationStatusRequest struct {
	Status NotificationStatus `json:"status"`
}

func (r *CreateNotificationRequest) ToNotification() *Notification {
	return &Notification{
		UserID:  r.UserID,
		Message: r.Message,
		Type:    r.Type,
		SentAt:  time.Now(),
	}
}
