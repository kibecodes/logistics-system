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
		INSERT INTO notifications (user_id, message, type, status, sent_at)
		VALUES (:user_id, :message, :type, :status, :sent_at)
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

func (r *NotificationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status notification.NotificationStatus) error {
	query := `
		UPDATE notifications 
		SET status = :status, updated_at = NOW() 
		WHERE id = :id
	`

	args := map[string]interface{}{
		"status": status,
		"id":     id,
	}

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, args)
	if err != nil {
		return fmt.Errorf("update notification status: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no notification found with id %s", id)
	}

	return nil
}

func (r *NotificationRepository) ListPending(ctx context.Context) ([]*notification.Notification, error) {
	query := `
		SELECT id, user_id, message, type, status, sent_at, created_at, updated_at
		FROM notifications
		WHERE status = 'pending'
		ORDER BY created_at ASC
	`
	var notifications []*notification.Notification
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &notifications, query)
	if err != nil {
		return nil, fmt.Errorf("list pending notifications: %w", err)
	}
	return notifications, nil
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*notification.Notification, error) {
	query := `
		SELECT id, user_id, message, type, status, sent_at, created_at, updated_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	var notifications []*notification.Notification
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &notifications, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list notifications by user: %w", err)
	}
	return notifications, nil
}
