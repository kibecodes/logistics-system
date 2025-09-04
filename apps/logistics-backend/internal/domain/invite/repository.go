package invite

import "github.com/google/uuid"

// Invite = pending entry, just an email + role + token + expiration.
type Repository interface {
	Create(invite *Invite) error              // POST
	GetByToken(token string) (*Invite, error) // GET
	ListPending() ([]*Invite, error)          // GET
	Delete(id uuid.UUID) error                // DELETE
}
