package postgres

import (
	"context"
	"fmt"
	"logistics-backend/internal/application"
	"logistics-backend/internal/domain/payment"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PaymentRepository struct {
	exec sqlx.ExtContext
}

func NewPaymentRepository(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{exec: db}
}

func (r *PaymentRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *PaymentRepository) Create(ctx context.Context, p *payment.Payment) error {
	query := `
		INSERT INTO payments (order_id, amount, method, status)
		VALUES (:order_id, :amount, :method, :status)
		RETURNING id
	`

	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, p)
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&p.ID); err != nil {
			return fmt.Errorf("scanning new payment id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after scan")
	}

	return nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*payment.Payment, error) {
	query := `
		SELECT id, order_id, amount, method, status, paid_at 
		FROM payments 
		WHERE id = $1
	`

	var p payment.Payment
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &p, query, id)
	return &p, err
}

func (r *PaymentRepository) GetByOrder(ctx context.Context, orderID uuid.UUID) (*payment.Payment, error) {
	query := `
		SELECT id, order_id, amount, status, paid_at FROM payments 
		WHERE order_id = $1
	`

	var p payment.Payment
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &p, query, orderID)
	return &p, err
}

func (r *PaymentRepository) List(ctx context.Context) ([]*payment.Payment, error) {
	query := `
		SELECT id, order_id, amount, status 
		FROM payments
	`

	var payments []*payment.Payment
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &payments, query)
	return payments, err
}

func (r *PaymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM payments
		WHERE id = $1
	`
	res, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify payment deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("payment already deleted or invalid")
	}

	return nil
}
