package delivery

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/delivery"
	"logistics-backend/internal/usecase/common"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      delivery.Repository
	ordRepo   delivery.OrderReader
	drvRepo   delivery.DriverReader
	txManager common.TxManager
}

func NewUseCase(repo delivery.Repository, ordRepo delivery.OrderReader, drvRepo delivery.DriverReader, txm common.TxManager) *UseCase {
	return &UseCase{repo: repo, ordRepo: ordRepo, drvRepo: drvRepo, txManager: txm}
}

func (uc *UseCase) CreateDelivery(ctx context.Context, d *delivery.Delivery) error {

	return uc.txManager.Do(ctx, func(txCtx context.Context) error {

		// 1. fetch order
		order, err := uc.ordRepo.GetOrderByID(ctx, d.OrderID)
		if err != nil {
			return fmt.Errorf("could not fetch order: %w", err)
		}
		if order.OrderStatus != "pending" {
			return delivery.ErrorNoPendingOrder
		}

		// 2. update order status using tx
		if err := uc.ordRepo.UpdateOrder(ctx, order.ID, "status", "assigned"); err != nil {
			return fmt.Errorf("could not update order status: %w", err)
		}

		// 3. create delivery using tx
		if err := uc.repo.Create(ctx, d); err != nil {
			return fmt.Errorf("could not create delivery: %w", err)
		}

		return nil
	})

}

func (uc *UseCase) GetDeliveryByID(ctx context.Context, deliveryId uuid.UUID) (*delivery.Delivery, error) {
	return uc.repo.GetByID(ctx, deliveryId)
}

func (uc *UseCase) UpdateDelivery(ctx context.Context, deliveryID uuid.UUID, column string, value any) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Update(txCtx, deliveryID, column, value); err != nil {
			return fmt.Errorf("update delivery failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) AcceptDelivery(ctx context.Context, d *delivery.Delivery) error {

	return uc.txManager.Do(ctx, func(txCtx context.Context) error {

		// 1. get driver by id
		driver, err := uc.drvRepo.GetDriverByID(txCtx, d.DriverID)
		if err != nil {
			return fmt.Errorf("could not fetch driver: %w", err)
		}
		if driver == nil {
			return fmt.Errorf("driver not found")
		}

		// 2. check availability - bool
		if !driver.Available {
			return fmt.Errorf("driver not available")
		}

		// 3. get order
		order, err := uc.ordRepo.GetOrderByID(txCtx, d.OrderID)
		if err != nil {
			return fmt.Errorf("could not fetch order: %w", err)
		}
		// 4. update order
		if err := uc.ordRepo.UpdateOrder(txCtx, order.ID, "status", "in_transit"); err != nil {
			return fmt.Errorf("could not update order status: %w", err)
		}

		// 5. accept delivery
		if err := uc.repo.Accept(txCtx, d); err != nil {
			return fmt.Errorf("could not accept delivery: %w", err)
		}

		// 6. update driver availability
		if err := uc.drvRepo.UpdateDriverAvailability(txCtx, driver.ID, "availability", false); err != nil {
			return fmt.Errorf("could not update driver availability: %w", err)
		}

		return nil
	})

}

func (uc *UseCase) ListDeliveries(ctx context.Context) ([]*delivery.Delivery, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) DeleteDelivery(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete delivery failed: %w", err)
		}

		return nil
	})
}
