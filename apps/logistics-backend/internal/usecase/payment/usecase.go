package payment

import (
	"context"
	"fmt"
	domain "logistics-backend/internal/domain/payment"
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

func (uc *UseCase) CreatePayment(ctx context.Context, p *domain.Payment) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Create(txCtx, p); err != nil {
			return fmt.Errorf("create payment failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetPaymentByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	return uc.repo.GetByOrder(ctx, orderID)
}

func (uc *UseCase) ListPayments(ctx context.Context) ([]*domain.Payment, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) DeletePayment(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete payment failed: %w", err)
		}

		return nil
	})
}
