package postgres

import (
	"context"
	"fmt"
	"logistics-backend/internal/application"
	"logistics-backend/internal/domain/driver"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DriverRepository struct {
	exec sqlx.ExtContext
}

func NewDriverRepository(db *sqlx.DB) *DriverRepository {
	return &DriverRepository{exec: db}
}

func (r *DriverRepository) execFromCtx(ctx context.Context) sqlx.ExtContext {
	if tx := application.GetTx(ctx); tx != nil {
		return tx
	}
	return r.exec
}

func (r *DriverRepository) Create(ctx context.Context, d *driver.Driver) error {
	query := `
		INSERT INTO drivers (id, full_name, email, vehicle_info, current_location, available, created_at)
		VALUES (:id, :full_name, :email, :vehicle_info, :current_location, :available, :created_at)
		RETURNING id
	`
	rows, err := sqlx.NamedQueryContext(ctx, r.execFromCtx(ctx), query, d)
	if err != nil {
		return fmt.Errorf("insert driver: %w", err)
	}

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&d.ID); err != nil {
			return fmt.Errorf("scanning new driver id: %w", err)
		}
	} else {
		return fmt.Errorf("no id returned after scan")
	}

	return nil
}

func (r *DriverRepository) UpdateProfile(ctx context.Context, driverID uuid.UUID, vehicleInfo string, currentLocation postgis.PointS) error {
	query := `
		UPDATE drivers 
		SET vehicle_info = :vehicle, current_location = :location
		WHERE id = :id
	`

	// Prepare WKT string
	wkt := fmt.Sprintf("SRID=%d;POINT(%f %f)", currentLocation.SRID, currentLocation.X, currentLocation.Y)

	args := map[string]interface{}{
		"vehicle":  vehicleInfo,
		"location": wkt,
		"id":       driverID,
	}

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, args)
	if err != nil {
		return fmt.Errorf("update driver profile: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no user found with id %s", driverID)
	}

	return nil
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

	query := fmt.Sprintf(`
		UPDATE drivers SET %s = $1, updated_at = NOW() 
		WHERE id = $2
	`, column)

	args := map[string]interface{}{
		"value": value,
		"id":    driverID,
	}

	res, err := sqlx.NamedExecContext(ctx, r.execFromCtx(ctx), query, args)
	if err != nil {
		return fmt.Errorf("update driver: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no user found with id %s", driverID)
	}

	return nil
}

func (r *DriverRepository) GetByID(ctx context.Context, id uuid.UUID) (*driver.Driver, error) {
	query := `
		SELECT id, full_name, email, vehicle_info, current_location, available, created_at 
		FROM drivers 
		WHERE id = $1
	`

	var d driver.Driver
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &d, query, id)
	return &d, err

}

func (r *DriverRepository) GetByEmail(ctx context.Context, email string) (*driver.Driver, error) {
	query := `
		SELECT id, full_name, email, vehicle_info, current_location, available, created_at 
		FROM drivers 
		WHERE email = $1
	`

	var d driver.Driver
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &d, query, email)
	return &d, err
}

func (r *DriverRepository) List(ctx context.Context) ([]*driver.Driver, error) {
	query := `
		SELECT id, full_name, email, vehicle_info, current_location, available, created_at 
		FROM drivers
	`

	var drivers []*driver.Driver
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &drivers, query)
	return drivers, err
}

func (r *DriverRepository) ListAvailableDrivers(ctx context.Context, available bool) ([]*driver.Driver, error) {
	query := `
		SELECT id, full_name, email, vehicle_info, current_location, available, created_at 
		FROM drivers
		WHERE available = $1
	`

	var drivers []*driver.Driver
	err := sqlx.SelectContext(ctx, r.execFromCtx(ctx), &drivers, query, available)
	return drivers, err
}

func (r *DriverRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM drivers 
		WHERE id = $1
	`
	res, err := r.exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not verify driver deletion: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("driver already deleted or invalid")
	}

	return nil
}

func (r *DriverRepository) GetNearestDriver(ctx context.Context, pickup postgis.PointS, maxDistance float64) (*driver.Driver, error) {
	query := `
		SELECT id, full_name, current_location, ST_Distance(current_location, $1) AS dist
		FROM drivers
		WHERE available = true
		AND ST_DWithin(current_location, $1, $2)
		ORDER BY current_location <-> $1
		LIMIT 1
	`

	var d driver.Driver
	err := sqlx.GetContext(ctx, r.execFromCtx(ctx), &d, query, pickup, maxDistance)
	return &d, err
}
