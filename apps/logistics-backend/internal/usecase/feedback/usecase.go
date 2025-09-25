package feedback

import (
	"context"
	"fmt"
	domain "logistics-backend/internal/domain/feedback"
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

func (uc *UseCase) CreateFeedback(ctx context.Context, f *domain.Feedback) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Create(txCtx, f); err != nil {
			return fmt.Errorf("create feecback failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetFeedbackByID(ctx context.Context, id uuid.UUID) (*domain.Feedback, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) ListFeedback(ctx context.Context) ([]*domain.Feedback, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) DeleteFeedback(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete feedback failed: %w", err)
		}

		return nil
	})
}
