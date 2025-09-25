package postgres

import (
	"context"
	"fmt"
	"logistics-backend/internal/application"
	"logistics-backend/internal/domain/feedback"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FeedbackRepository struct {
	exec sqlx.ExtContext
}

func NewFeedbackRepository(db *sqlx.DB) *FeedbackRepository {
	return &FeedbackRepository{exec: db}
}

func (r *FeedbackRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *FeedbackRepository) Create(ctx context.Context, f *feedback.Feedback) error {
	query := `
		INSERT INTO feedbacks (order_id, customer_id, rating, comments)
		VALUES (:order_id, :customer_id, :rating, :comments)
		RETURNING id
	`

	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, f)
	if err != nil {
		return fmt.Errorf("insert feedback: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&f.ID); err != nil {
			return fmt.Errorf("scanning new feedback id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after scan")
	}

	return nil
}

func (r *FeedbackRepository) GetByID(ctx context.Context, id uuid.UUID) (*feedback.Feedback, error) {
	query := `
		SELECT id, order_id, customer_id, rating, comments FROM feedbacks 
		WHERE id = $1
	`

	var f feedback.Feedback
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &f, query, id)
	return &f, err
}

func (r *FeedbackRepository) List(ctx context.Context) ([]*feedback.Feedback, error) {
	query := `
		SELECT id, order_id, customer_id, rating, comments 
		FROM feedbacks
	`

	var feedbacks []*feedback.Feedback
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &feedbacks, query)
	return feedbacks, err
}

func (r *FeedbackRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM feedbacks
		WHERE id = $1
	`

	res, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete feedback: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify feedback deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("feedback already deleted or invalid")
	}

	return nil
}
