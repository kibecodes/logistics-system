package invite

import (
	"context"

	"github.com/google/uuid"
)

// Invite = pending entry, just an email + role + token + expiration.
type Repository interface {
	Create(ctx context.Context, invite *Invite) error              // POST
	GetByToken(ctx context.Context, token string) (*Invite, error) // GET
	ListPending(ctx context.Context) ([]*Invite, error)            // GET
	Delete(ctx context.Context, id uuid.UUID) error                // DELETE
}
