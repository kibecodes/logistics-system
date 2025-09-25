package postgres

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/delivery"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DeliveryRepository struct {
	exec sqlx.ExtContext
}

func NewDeliveryRepository(db *sqlx.DB) *DeliveryRepository {
	return &DeliveryRepository{exec: db}
	//flexible — can pass in either a *sqlx.DB or a *sqlx.Tx
}

func (r *DeliveryRepository) Create(ctx context.Context, d *delivery.Delivery) error {
	query := `
		INSERT INTO deliveries (order_id, driver_id, status)
		VALUES (:order_id, :driver_id, :status)
		RETURNING id
	`

	rows, err := sqlx.NamedQueryContext(ctx, r.exec, query, d)
	if err != nil {
		return fmt.Errorf("insert delivery: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&d.ID); err != nil {
			return fmt.Errorf("scanning new delivery id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after scan")
	}

	return nil
}

func (r *DeliveryRepository) GetByID(ctx context.Context, id uuid.UUID) (*delivery.Delivery, error) {
	query := `
		SELECT id, order_id, driver_id, assigned_at, picked_up_at, delivered_at, status 
		FROM deliveries 
		WHERE id = $1
	`

	var d delivery.Delivery
	if err := sqlx.GetContext(ctx, r.exec, &d, query, id); err != nil {
		return nil, fmt.Errorf("get delivery by id: %w", err)
	}

	return &d, nil
}

func (r *DeliveryRepository) Update(ctx context.Context, deliveryID uuid.UUID, column string, value any) error {
	// whitelist columns
	allowed := map[string]bool{
		"status": true,
	}

	if !allowed[column] {
		return fmt.Errorf("attempted to update disallowed column: %s", column)
	}

	query := fmt.Sprintf(`
		UPDATE deliveries 
		SET %s = :value, updated_at = NOW() 
		WHERE id = :id
	`, column)

	args := map[string]interface{}{
		"value": value,
		"id":    deliveryID,
	}

	res, err := sqlx.NamedExecContext(ctx, r.exec, query, args)
	if err != nil {
		return fmt.Errorf("update delivery: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no delivery found with id %s", deliveryID)
	}

	return nil
}

func (r *DeliveryRepository) Accept(ctx context.Context, d *delivery.Delivery) error {
	query := `
		UPDATE deliveries
		SET status = 'picked_up',
			picked_up_at = :picked_up_at,
			driver_id = :driver_id,
			updated_at = NOW()
		WHERE id = :id AND status = 'assigned'
	`

	res, err := sqlx.NamedExecContext(ctx, r.exec, query, d)
	if err != nil {
		return fmt.Errorf("failed to accept delivery: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify delivery update: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("delivery already picked or invalid")
	}

	return nil
}

func (r *DeliveryRepository) List(ctx context.Context) ([]*delivery.Delivery, error) {
	query := `
		SELECT id, order_id, driver_id, assigned_at, picked_up_at, delivered_at, status 
		FROM deliveries
	`
	var deliveries []*delivery.Delivery

	err := sqlx.SelectContext(ctx, r.exec, &deliveries, query)
	return deliveries, err
}

func (r *DeliveryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM deliveries 
		WHERE id = $1
	`

	res, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete delivery: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify delivery deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("delivery already deleted or invalid")
	}

	return nil
}
