package invite

import (
	"context"
	domain "logistics-backend/internal/domain/invite"

	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
	// optional user repo - auto-create users on accept
}

func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Add a new invited member from the shareable invite link
func (uc *UseCase) InviteMember(ctx context.Context, i *domain.Invite) error {
	// generate token, store invite
	// optionally notify user
	return uc.repo.Create(i)
}

// Get invited member by token
func (uc *UseCase) GetMemberByToken(ctx context.Context, token string) (*domain.Invite, error) {
	return uc.repo.GetByToken(token)
}

// List all invited members
func (uc *UseCase) ListPendingMembers(ctx context.Context) ([]*domain.Invite, error) {
	return uc.repo.ListPending()
}

// Delete selected invited member
func (uc *UseCase) DeleteMember(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(id)
}
