package orderadapter

import (
	"context"
	"fmt"
	"log"
	"logistics-backend/internal/domain/driver"
	"logistics-backend/internal/domain/order"
	driverusecase "logistics-backend/internal/usecase/driver"
	orderusecase "logistics-backend/internal/usecase/order"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// implements delivery.OrderReader
type OrderUseCaseAdapter struct {
	UseCase       *orderusecase.UseCase
	DriverUseCase *driverusecase.UseCase
}

func (a *OrderUseCaseAdapter) GetOrderByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	return a.UseCase.GetOrder(ctx, id)
}

func (a *OrderUseCaseAdapter) UpdateOrderTx(ctx context.Context, tx *sqlx.Tx, orderID uuid.UUID, column string, value any) error {
	return a.UseCase.UpdateOrderTx(ctx, tx, orderID, column, value)
}

func (a *OrderUseCaseAdapter) GetDriverByID(ctx context.Context, id uuid.UUID) (*driver.Driver, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("no driver id")
	}

	log.Printf("adapter driver id: %+v", id)
	driver, err := a.DriverUseCase.GetDriver(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repo error: %w", err)
	}
	if driver == nil {
		return nil, fmt.Errorf("driver not found for id %s", id)
	}

	return driver, nil
}

func (a *OrderUseCaseAdapter) UpdateDriverAvailability(ctx context.Context, driverID uuid.UUID, column string, value bool) error {
	return a.DriverUseCase.UpdateDriverAvailability(ctx, driverID, column, value)
}
