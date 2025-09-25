package postgres

import (
	"context"
	"fmt"
	"logistics-backend/internal/application"
	"logistics-backend/internal/domain/notification"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type NotificationRepository struct {
	exec sqlx.ExtContext
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{exec: db}
}

func (r *NotificationRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *NotificationRepository) Create(ctx context.Context, n *notification.Notification) error {
	query := `
		INSERT INTO notifications (user_id, message, type)
		VALUES (:user_id, :message, :type)
		RETURNING id
	`
	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, n)
	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&n.ID); err != nil {
			return fmt.Errorf("scanning new notification id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after scan")
	}

	return nil
}

func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*notification.Notification, error) {
	query := `
		SELECT id, user_id, message, type 
		FROM notifications 
		WHERE id = $1
	`

	var n notification.Notification
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &n, query, id)
	return &n, err
}

func (r *NotificationRepository) List(ctx context.Context) ([]*notification.Notification, error) {
	query := `
		SELECT id, user_id, message, type 
		FROM notifications
	`
	var notifications []*notification.Notification
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &notifications, query)
	return notifications, err
}

func (r *NotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM notifications 
		WHERE id = $1
	`
	res, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify notification deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("notification already deleted or invalid")
	}

	return nil
}
