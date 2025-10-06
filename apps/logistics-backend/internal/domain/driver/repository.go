package driver

import (
	"context"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, driver *Driver) error                                                          // POST
	GetByID(ctx context.Context, id uuid.UUID) (*Driver, error)                                                // GET
	GetByEmail(ctx context.Context, email string) (*Driver, error)                                             // GET
	List(ctx context.Context) ([]*Driver, error)                                                               // GET all drivers
	UpdateColumn(ctx context.Context, driverID uuid.UUID, column string, value any) error                      // PATCH method for specific driver details update
	UpdateProfile(ctx context.Context, id uuid.UUID, vehicleInfo string, currentLocation postgis.PointS) error // PUT method for driver details to be updated after registration
	Delete(ctx context.Context, id uuid.UUID) error                                                            // DELETE

	GetNearestDriver(ctx context.Context, pickup postgis.PointS, maxDistance float64) (*Driver, error)
	ListAvailableDrivers(ctx context.Context, available bool) ([]*Driver, error)
}
