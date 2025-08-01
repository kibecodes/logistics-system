package user

import (
	"context"
	"logistics-backend/internal/domain/driver"
)

type CreateDriver interface {
	RegisterDriver(ctx context.Context, d *driver.Driver) error
}
