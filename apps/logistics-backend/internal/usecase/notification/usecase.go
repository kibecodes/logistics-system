package notification

import (
	"context"
	"fmt"
	domain "logistics-backend/internal/domain/notification"
	"logistics-backend/internal/usecase/common"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      domain.Repository
	txManager common.TxManager
}

func NewUseCase(repo domain.Repository, txm common.TxManager) *UseCase {
	return &UseCase{repo: repo, txManager: txm}
}

func (uc *UseCase) CreateNotification(ctx context.Context, n *domain.Notification) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Create(txCtx, n); err != nil {
			return fmt.Errorf("create notification failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) UpdateNotificationStatus(ctx context.Context, id uuid.UUID, status domain.NotificationStatus) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateStatus(txCtx, id, status); err != nil {
			return fmt.Errorf("update notification failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) ListPendingNotifications(ctx context.Context) ([]*domain.Notification, error) {
	return uc.repo.ListPending(ctx)
}

func (uc *UseCase) ListNotificationsByCustomer(ctx context.Context, userID uuid.UUID) ([]*domain.Notification, error) {
	return uc.repo.ListByUser(ctx, userID)
}

func (uc *UseCase) SendNotification(ctx context.Context, n *domain.Notification, sender domain.Sender) error {
	if err := uc.repo.Create(ctx, n); err != nil {
		return err
	}

	if err := sender.Send(ctx, n); err != nil {
		_ = uc.repo.UpdateStatus(ctx, n.ID, domain.Failed)
		return err
	}

	return uc.repo.UpdateStatus(ctx, n.ID, domain.Sent)
}

// TODO: Implement actual notification sending via external services (e.g., Twilio, SendGrid).
// The system should select the appropriate sender based on Notification.Type (email, sms, push)
// and update the notification status accordingly.
