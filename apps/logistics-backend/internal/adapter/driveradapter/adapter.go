package driveradapter

import (
	"context"
	"logistics-backend/internal/domain/driver"
	driverusecase "logistics-backend/internal/usecase/driver"
)

type DriverUseCaseAdapter struct {
	UseCase *driverusecase.UseCase
}

func (a *DriverUseCaseAdapter) RegisterDriver(ctx context.Context, d *driver.Driver) error {
	return a.UseCase.RegisterDriver(ctx, d)
}
