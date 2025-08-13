package user

import (
	"context"
	"logistics-backend/internal/domain/driver"
)

type DriverReader interface {
	RegisterDriver(ctx context.Context, d *driver.Driver) error
}
