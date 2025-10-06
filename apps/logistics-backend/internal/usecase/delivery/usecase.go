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
		driver, err := uc.drvRepo.GetDriverByID(txCtx, d.DriverID)
		if err != nil || !driver.Available {
			return fmt.Errorf("driver not available")
		}

		order, err := uc.ordRepo.GetOrderByID(txCtx, d.OrderID)
		if err != nil {
			return err
		}

		if err := uc.ordRepo.UpdateOrder(txCtx, order.ID, "status", "in_transit"); err != nil {
			return err
		}

		if err := uc.drvRepo.UpdateDriverAvailability(txCtx, driver.ID, "availability", false); err != nil {
			return err
		}

		// Remove pending assignment
		// if err := uc.ordRepo.MarkAssignmentAccepted(txCtx, order.ID, driver.ID); err != nil {
		//     return err
		// }

		return uc.repo.Create(txCtx, d)
	})

}

func (uc *UseCase) ListDeliveries(ctx context.Context) ([]*delivery.Delivery, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) ListActiveDeliveries(ctx context.Context) ([]*delivery.Delivery, error) {
	ativeStatuses := []delivery.DeliveryStatus{
		delivery.Assigned,
		delivery.PickedUp,
	}

	return uc.repo.ListByStatus(ctx, ativeStatuses)
}

func (uc *UseCase) DeleteDelivery(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete delivery failed: %w", err)
		}

		return nil
	})
}
