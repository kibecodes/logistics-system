package driver

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, driver *Driver) error                                           // POST
	GetByID(id uuid.UUID) (*Driver, error)                                                      // GET
	GetByEmail(email string) (*Driver, error)                                                   // GET
	List() ([]*Driver, error)                                                                   // GET
	UpdateColumn(ctx context.Context, driverID uuid.UUID, column string, value any) error       // PATCH method for specific driver details update
	UpdateProfile(ctx context.Context, id uuid.UUID, vehicleInfo, currentLocation string) error // PUT method for driver details to be updated after registration
	Delete(ctx context.Context, id uuid.UUID) error                                             // DELETE
}
