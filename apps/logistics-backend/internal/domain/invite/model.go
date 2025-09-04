package invite

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	Admin    Role = "admin"
	Driver   Role = "driver"
	Customer Role = "customer"
	Guest    Role = "guest"
)

// new user invite via invite link
type Invite struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Role      Role      `db:"role" json:"role"`
	Token     string    `db:"token" json:"token"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	InvitedBy uuid.UUID `db:"invited_by" json:"invited_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
