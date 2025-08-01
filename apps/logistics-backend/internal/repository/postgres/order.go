package postgres

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/order"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) order.Repository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(o *order.Order) error {
	query := `
		INSERT INTO orders (user_id, inventory_id, quantity, pickup_address, delivery_address, status)
		VALUES (:user_id, :inventory_id, :quantity, :pickup_address, :delivery_address, :status)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&o.ID, o)
}

func (r *OrderRepository) GetByID(id uuid.UUID) (*order.Order, error) {
	query := `SELECT id, user_id, inventory_id, quantity, pickup_address, delivery_address, status, created_at, updated_at FROM orders WHERE id = $1`
	var o order.Order
	err := r.db.Get(&o, query, id)
	return &o, err
}

func (r *OrderRepository) ListByCustomer(customerID uuid.UUID) ([]*order.Order, error) {
	query := `SELECT id, user_id, inventory_id, quantity, pickup_address, delivery_address, status, created_at, updated_at FROM orders WHERE user_id = $1`
	var orders []*order.Order
	err := r.db.Select(&orders, query, customerID)
	return orders, err
}

// non transaction - simple update method
func (r *OrderRepository) UpdateColumn(ctx context.Context, orderID uuid.UUID, column string, value any) error {
	// Validate column name to avoid SQL injection
	allowed := map[string]bool{
		"status":           true,
		"quantity":         true,
		"pickup_address":   true,
		"delivery_address": true,
	}

	if !allowed[column] {
		return fmt.Errorf("attempted to update disallowed column: %s", column)
	}

	query := fmt.Sprintf(`UPDATE orders SET %s = $1, updated_at = NOW() WHERE id = $2`, column)
	res, err := r.db.ExecContext(ctx, query, value, orderID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no order found with id %s", orderID)
	}
	return nil
}

func (r *OrderRepository) UpdateColumnTx(ctx context.Context, tx *sqlx.Tx, orderID uuid.UUID, column string, value any) error {
	// Validate column name to avoid SQL injection
	allowed := map[string]bool{
		"status":           true,
		"quantity":         true,
		"pickup_address":   true,
		"delivery_address": true,
	}

	if !allowed[column] {
		return fmt.Errorf("attempted to update disallowed column: %s", column)
	}

	query := fmt.Sprintf(`UPDATE orders SET %s = $1, updated_at = NOW() WHERE id = $2`, column)
	res, err := tx.ExecContext(ctx, query, value, orderID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no order found with id %s", orderID)
	}
	return nil
}

// func (r *OrderRepository) UpdateStatus(orderID uuid.UUID, status order.OrderStatus) error {
// 	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
// 	res, err := r.db.Exec(query, status, orderID)
// 	if err != nil {
// 		return err
// 	}
// 	rows, _ := res.RowsAffected()
// 	if rows == 0 {
// 		return fmt.Errorf("no order found with id %s", orderID)
// 	}
// 	return nil
// }

func (r *OrderRepository) List() ([]*order.Order, error) {
	query := `SELECT id, user_id, inventory_id, quantity, pickup_address, delivery_address, status, created_at, updated_at FROM orders`
	var orders []*order.Order
	err := r.db.Select(&orders, query)
	return orders, err
}

func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM orders WHERE id = $1 `
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
