package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"logistics-backend/internal/domain/driver"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DriverRepository struct {
	db *sqlx.DB
}

func NewDriverRepository(db *sqlx.DB) driver.Repository {
	return &DriverRepository{db: db}
}

func (r *DriverRepository) Create(ctx context.Context, d *driver.Driver) error {
	query := `
		INSERT INTO drivers (id, full_name, email, vehicle_info, current_location, available, created_at)
		VALUES (:id, :full_name, :email, :vehicle_info, :current_location, :available, :created_at)
		RETURNING id
	`
	_, err := r.db.NamedExecContext(ctx, query, d)
	return err
}

func (r *DriverRepository) UpdateProfile(ctx context.Context, id uuid.UUID, vehicleInfo string, currentLocation string) error {
	query := `
		UPDATE drivers 
		SET vehicle_info = $1, current_location = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, vehicleInfo, currentLocation, id)
	return err
}

func (r *DriverRepository) UpdateColumn(ctx context.Context, driverID uuid.UUID, column string, value any) error {
	allowed := map[string]bool{
		"full_name":        true,
		"email":            true,
		"vehicle_info":     true,
		"current_location": true,
	}

	if !allowed[column] {
		return fmt.Errorf("attempted to update disallowed column: %s", column)
	}

	query := fmt.Sprintf(`UPDATE drivers SET %s = $1, updated_at = NOW() WHERE id = $2`, column)
	res, err := r.db.ExecContext(ctx, query, value, driverID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no driver found with id %s", driverID)
	}
	return nil
}

func (r *DriverRepository) GetByID(id uuid.UUID) (*driver.Driver, error) {
	query := `SELECT id, full_name, email, vehicle_info, current_location, available, created_at FROM drivers WHERE id = $1`
	var d driver.Driver
	err := r.db.Get(&d, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("driver with id %s not found", id)
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	return &d, err
}

func (r *DriverRepository) GetByEmail(email string) (*driver.Driver, error) {
	query := `SELECT id, full_name, email, vehicle_info, current_location, available, created_at FROM drivers WHERE email = $1`
	var d driver.Driver
	err := r.db.Get(&d, query, email)
	return &d, err
}

func (r *DriverRepository) List() ([]*driver.Driver, error) {
	query := `SELECT id, full_name, email, vehicle_info, current_location, available, created_at FROM drivers`
	var drivers []*driver.Driver
	err := r.db.Select(&drivers, query)
	return drivers, err
}

func (r *DriverRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM drivers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
