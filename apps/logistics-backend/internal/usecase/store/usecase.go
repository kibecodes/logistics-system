package store

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/store"
	"logistics-backend/internal/usecase/common"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      store.Repository
	txManager common.TxManager
}

func NewUseCase(repo store.Repository, txm common.TxManager) *UseCase {
	return &UseCase{repo: repo, txManager: txm}
}

func (uc *UseCase) CreateStore(ctx context.Context, s *store.Store) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Create(txCtx, s); err != nil {
			return fmt.Errorf("could not create order: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetStoreByID(ctx context.Context, id uuid.UUID) (*store.Store, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetStoreBySlug(ctx context.Context, slug string) (*store.Store, error) {
	return uc.repo.GetBySlug(ctx, slug)
}

func (uc *UseCase) UpdateStore(ctx context.Context, storeID uuid.UUID, column string, value any) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Update(txCtx, storeID, column, value); err != nil {
			return fmt.Errorf("update store failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetPublicStores(ctx context.Context) ([]*store.Store, error) {
	return uc.repo.ListPublic(ctx)
}

func (uc *UseCase) DeleteStore(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete store failed: %w", err)
		}

		return nil
	})
}
