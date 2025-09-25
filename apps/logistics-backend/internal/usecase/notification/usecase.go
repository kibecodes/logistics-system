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

func (uc *UseCase) GetNotification(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) ListNotification(ctx context.Context) ([]*domain.Notification, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) DeleteNotification(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete notification failed: %w", err)
		}

		return nil
	})
}
