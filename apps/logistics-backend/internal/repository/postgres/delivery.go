package postgres

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/delivery"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DeliveryRepository struct {
	db *sqlx.DB
}

func NewDeliveryRepository(db *sqlx.DB) delivery.Repository {
	return &DeliveryRepository{db: db}
}

// insert a delivery independently
func (r *DeliveryRepository) Create(d *delivery.Delivery) error {
	query := `
		INSERT INTO deliveries (order_id, driver_id, status)
		VALUES (:order_id, :driver_id, :status)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&d.ID, d)
}

func (r *DeliveryRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

// inside a transaction involving multiple steps - Creating delivery + updating order
func (r *DeliveryRepository) CreateTx(ctx context.Context, tx *sqlx.Tx, d *delivery.Delivery) error {
	query := `
		INSERT INTO deliveries (order_id, driver_id, status)
		VALUES (:order_id, :driver_id, :status)
		RETURNING id
	`

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.GetContext(ctx, &d.ID, d)
}

func (r *DeliveryRepository) GetByID(id uuid.UUID) (*delivery.Delivery, error) {
	query := `SELECT id, order_id, driver_id, assigned_at, picked_up_at, delivered_at, status FROM deliveries WHERE id = $1`
	var d delivery.Delivery
	err := r.db.Get(&d, query, id)
	return &d, err
}

func (r *DeliveryRepository) Update(ctx context.Context, deliveryID uuid.UUID, column string, value any) error {
	// whitelist columns
	allowed := map[string]bool{
		"status": true,
	}

	if !allowed[column] {
		return fmt.Errorf("attempted to update disallowed column: %s", column)
	}

	query := fmt.Sprintf(`UPDATE deliveries SET %s = $1, updated_at = NOW() WHERE id = $2`, column)
	res, err := r.db.ExecContext(ctx, query, value, deliveryID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no order found with id %s", deliveryID)
	}
	return nil
}

func (r *DeliveryRepository) Accept(ctx context.Context, d *delivery.Delivery) error {
	query := `
		UPDATE deliveries
		SET status = 'picked_up',
			picked_up_at = $1,
			driver_id = $2,
			updated_at = NOW()
		WHERE id = $3 AND status = 'assigned'
	`
	res, err := r.db.ExecContext(ctx, query, d.PickedUpAt, d.DriverID, d.ID)
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

func (r *DeliveryRepository) List() ([]*delivery.Delivery, error) {
	query := `SELECT id, order_id, driver_id, assigned_at, picked_up_at, delivered_at, status FROM deliveries`
	var deliveries []*delivery.Delivery
	err := r.db.Select(&deliveries, query)
	return deliveries, err
}

func (r *DeliveryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM deliveries WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
