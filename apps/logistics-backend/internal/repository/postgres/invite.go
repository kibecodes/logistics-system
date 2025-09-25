package postgres

import (
	"context"
	"fmt"
	"logistics-backend/internal/application"
	"logistics-backend/internal/domain/invite"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InviteRepository struct {
	exec sqlx.ExtContext
}

func NewInviteRepository(db *sqlx.DB) *InviteRepository {
	return &InviteRepository{exec: db}
}

func (r *InviteRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *InviteRepository) Create(ctx context.Context, i *invite.Invite) error {
	query := `
		INSERT INTO invites (id, email, role, token, expires_at, invited_by, created_at)
		VALUES (:id, :email, :role, :token, :expires_at, :invited_by, :created_at)
	`
	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, i)
	if err != nil {
		return fmt.Errorf("insert invite: %w", err)
	}

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&i.ID); err != nil {
			return fmt.Errorf("scanning new invite id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after scan")
	}

	return nil
}

func (r *InviteRepository) GetByToken(ctx context.Context, token string) (*invite.Invite, error) {
	var i invite.Invite
	query := `
		SELECT id, email, role, token, expires_at, invited_by, created_at 
		FROM invites 
		WHERE token = $1
	`
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &i, query, token)
	return &i, err
}

func (r *InviteRepository) ListPending(ctx context.Context) ([]*invite.Invite, error) {
	var invites []*invite.Invite
	query := `
		SELECT id, email, role, token, expires_at, invited_by, created_at
        FROM invites
        WHERE expires_at > NOW()
	`
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &invites, query)
	return invites, err
}

func (r *InviteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM invites 
		WHERE id = $1
	`

	res, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify invite deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("invite already deleted or invalid")
	}

	return nil
}
