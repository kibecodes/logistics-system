package invite

import (
	"time"

	"github.com/google/uuid"
)

type CreateInviteRequest struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Role      Role      `json:"role" binding:"required"`
	Token     string    `json:"token" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
	InvitedBy uuid.UUID `json:"invited_by" binding:"required"`
}

func (r *CreateInviteRequest) ToInvite() *Invite {
	return &Invite{
		ID:        r.ID,
		Email:     r.Email,
		Role:      r.Role,
		Token:     r.Token,
		ExpiresAt: r.ExpiresAt,
		InvitedBy: r.InvitedBy,
	}
}
