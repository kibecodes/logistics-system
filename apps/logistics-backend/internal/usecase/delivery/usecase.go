package delivery

import (
	"context"
	"fmt"
	"log"
	"logistics-backend/internal/domain/delivery"

	"github.com/google/uuid"
)

type UseCase struct {
	repo    delivery.Repository
	ordRepo delivery.OrderReader
	drvRepo delivery.DriverReader
}

func NewUseCase(repo delivery.Repository, ordRepo delivery.OrderReader, drvRepo delivery.DriverReader) *UseCase {
	return &UseCase{repo: repo, ordRepo: ordRepo, drvRepo: drvRepo}
}

func (uc *UseCase) CreateDelivery(ctx context.Context, d *delivery.Delivery) error {
	tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. fetch order
	order, err := uc.ordRepo.GetOrderByID(ctx, d.OrderID)
	if err != nil {
		return fmt.Errorf("could not fetch order: %w", err)
	}
	if order.OrderStatus != "pending" {
		return delivery.ErrorNoPendingOrder
	}

	// 2. update order status using tx
	if err := uc.ordRepo.UpdateOrderTx(ctx, tx, order.ID, "status", "assigned"); err != nil {
		return fmt.Errorf("could not update order status: %w", err)
	}

	// 3. create delivery using tx
	if err := uc.repo.CreateTx(ctx, tx, d); err != nil {
		return fmt.Errorf("could not create delivery: %w", err)
	}

	// 4. commit
	return tx.Commit()
}

func (uc *UseCase) GetDeliveryByID(ctx context.Context, deliveryId uuid.UUID) (*delivery.Delivery, error) {
	return uc.repo.GetByID(deliveryId)
}

func (uc *UseCase) UpdateDelivery(ctx context.Context, deliveryID uuid.UUID, column string, value any) error {
	return uc.repo.Update(ctx, deliveryID, column, value)
}

func (uc *UseCase) AcceptDelivery(ctx context.Context, d *delivery.Delivery) error {
	// get driver by id
	log.Printf("delivery usecase driver id: %+v", d.DriverID)
	driver, err := uc.drvRepo.GetDriverByID(ctx, d.DriverID)
	if err != nil {
		log.Printf("ERROR fetching driver: %v", err)
		return fmt.Errorf("could not fetch driver: %w", err)
	}
	if driver == nil {
		log.Printf("Driver is nil for ID: %s", d.DriverID)
		return fmt.Errorf("driver not found")
	}

	// check availability - bool
	if !driver.Available {
		return fmt.Errorf("driver not available")
	}

	tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	// get order
	order, err := uc.ordRepo.GetOrderByID(ctx, d.OrderID)
	if err != nil {
		return fmt.Errorf("could not fetch order: %w", err)
	}
	// update order
	if err := uc.ordRepo.UpdateOrderTx(ctx, tx, order.ID, "status", "in_transit"); err != nil {
		return fmt.Errorf("could not update order status: %w", err)
	}

	// accept delivery
	if err := uc.repo.Accept(ctx, d); err != nil {
		return fmt.Errorf("could not accept delivery: %w", err)
	}

	// update driver availability
	if err := uc.drvRepo.UpdateDriverAvailability(ctx, driver.ID, "availability", false); err != nil {
		return fmt.Errorf("could not update driver availability: %w", err)
	}

	return tx.Commit()
}

func (uc *UseCase) ListDeliveries(ctx context.Context) ([]*delivery.Delivery, error) {
	return uc.repo.List()
}

func (uc *UseCase) DeleteDelivery(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
