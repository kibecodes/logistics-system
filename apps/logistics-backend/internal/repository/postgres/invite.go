package postgres

import (
	"logistics-backend/internal/domain/invite"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InviteRepository struct {
	db *sqlx.DB
}

func NewInviteRepository(db *sqlx.DB) invite.Repository {
	return &InviteRepository{db: db}
}

func (r *InviteRepository) Create(i *invite.Invite) error {
	query := `
		INSERT INTO invites (id, email, role, token, expires_at, invited_by, created_at)
		VALUES (:id, :email, :role, :token, :expires_at, :invited_by, :created_at)
	`

	_, err := r.db.NamedExec(query, i)
	return err
}

func (r *InviteRepository) GetByToken(token string) (*invite.Invite, error) {
	var i invite.Invite
	query := `
		SELECT id, email, role, token, expires_at, invited_by, created_at 
		FROM invites 
		WHERE token = $1
	`
	err := r.db.Get(&i, query, token)
	return &i, err
}

func (r *InviteRepository) ListPending() ([]*invite.Invite, error) {
	var invites []*invite.Invite
	query := `
		SELECT id, email, role, token, expires_at, invited_by, created_at
        FROM invites
        WHERE expires_at > NOW()
	`
	err := r.db.Select(&invites, query)
	return invites, err
}

func (r *InviteRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM invites WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
