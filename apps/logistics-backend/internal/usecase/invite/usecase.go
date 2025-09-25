package invite

import (
	"context"
	"fmt"
	domain "logistics-backend/internal/domain/invite"
	"logistics-backend/internal/usecase/common"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      domain.Repository
	txManager common.TxManager
	// optional user repo - auto-create users on accept
}

func NewUseCase(repo domain.Repository, txm common.TxManager) *UseCase {
	return &UseCase{repo: repo, txManager: txm}
}

// Add a new invited member from the shareable invite link
func (uc *UseCase) InviteMember(ctx context.Context, i *domain.Invite) error {
	// generate token, store invite
	// optionally notify user
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Create(txCtx, i); err != nil {
			return fmt.Errorf("create invite failed: %w", err)
		}

		return nil
	})
}

// Get invited member by token
func (uc *UseCase) GetMemberByToken(ctx context.Context, token string) (*domain.Invite, error) {
	return uc.repo.GetByToken(ctx, token)
}

// List all invited members
func (uc *UseCase) ListPendingMembers(ctx context.Context) ([]*domain.Invite, error) {
	return uc.repo.ListPending(ctx)
}

// Delete selected invited member
func (uc *UseCase) DeleteMember(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete invite failed: %w", err)
		}

		return nil
	})
}
