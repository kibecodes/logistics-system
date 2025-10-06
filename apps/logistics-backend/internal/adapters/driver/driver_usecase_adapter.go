package driveradapter

import (
	"context"
	"fmt"
	"log"
	"logistics-backend/internal/domain/driver"
	driverusecase "logistics-backend/internal/usecase/driver"

	"github.com/google/uuid"
)

type UseCaseAdapter struct {
	UseCase *driverusecase.UseCase
}

func (a *UseCaseAdapter) GetDriverByID(ctx context.Context, id uuid.UUID) (*driver.Driver, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("no driver id")
	}

	log.Printf("adapter driver id: %+v", id)
	driver, err := a.UseCase.GetDriver(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repo error: %w", err)
	}
	if driver == nil {
		return nil, fmt.Errorf("driver not found for id %s", id)
	}

	return driver, nil
}

func (a *UseCaseAdapter) UpdateDriverAvailability(ctx context.Context, driverID uuid.UUID, column string, value bool) error {
	return a.UseCase.UpdateDriverAvailability(ctx, driverID, column, value)
}
