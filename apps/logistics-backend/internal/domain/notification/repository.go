package notification

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status NotificationStatus) error
	ListPending(ctx context.Context) ([]*Notification, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*Notification, error)
}

// Sender defines a generic interface for sending notifications.
// Concrete implementations (Twilio, SendGrid, Firebase, etc.)
// will satisfy this interface.
type Sender interface {
	Send(ctx context.Context, n *Notification) error
}

type EmailSender interface {
	SendEmail(ctx context.Context, to string, subject string, body string) error
}

type SMSSender interface {
	SendSMS(ctx context.Context, phone string, message string) error
}

type PushSender interface {
	SendPush(ctx context.Context, deviceToken string, title string, message string) error
}

type MultiChannelSender struct {
	emailSender EmailSender
	smsSender   SMSSender
	pushSender  PushSender
}

func (s *MultiChannelSender) Send(ctx context.Context, n *Notification) error {
	switch n.Type {
	case Email:
		return s.emailSender.SendEmail(ctx, n.UserID.String(), "Notification", n.Message)
	case SMS:
		return s.smsSender.SendSMS(ctx, n.UserID.String(), n.Message)
	case Push:
		return s.pushSender.SendPush(ctx, n.UserID.String(), "Notification", n.Message)
	default:
		return nil // or mark as system
	}
}
