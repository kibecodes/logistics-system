package delivery

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/delivery"
	"logistics-backend/internal/domain/notification"
	"logistics-backend/internal/usecase/common"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      delivery.Repository
	ordRepo   delivery.OrderReader
	drvRepo   delivery.DriverReader
	txManager common.TxManager
	notfRepo  delivery.NotificationReader
}

func NewUseCase(repo delivery.Repository, ordRepo delivery.OrderReader, drvRepo delivery.DriverReader, txm common.TxManager, notf delivery.NotificationReader) *UseCase {
	return &UseCase{repo: repo, ordRepo: ordRepo, drvRepo: drvRepo, txManager: txm, notfRepo: notf}
}

func (uc *UseCase) GetDeliveryByID(ctx context.Context, deliveryId uuid.UUID) (*delivery.Delivery, error) {
	return uc.repo.GetByID(ctx, deliveryId)
}

func (uc *UseCase) UpdateDelivery(ctx context.Context, deliveryID uuid.UUID, column string, value any) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		d, err := uc.repo.GetByID(txCtx, deliveryID)
		if err != nil {
			return fmt.Errorf("could not fetch delivery: %w", err)
		}

		if err := uc.repo.Update(txCtx, deliveryID, column, value); err != nil {
			return fmt.Errorf("update delivery failed: %w", err)
		}

		go func() {
			msg := fmt.Sprintf("ℹ️ Delivery for order %s updated: '%s' changed.", d.OrderID, column)
			_ = uc.notify(ctx, d.DriverID, msg)
		}()

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

		go func() {
			msgCustomer := fmt.Sprintf("🚚 Your order %s is now in transit with driver %s.", order.ID, driver.FullName)
			_ = uc.notify(ctx, order.CustomerID, msgCustomer)

			msgDriver := fmt.Sprintf("✅ You have accepted delivery for order %s.", order.ID)
			_ = uc.notify(ctx, driver.ID, msgDriver)
		}()

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

func (uc *UseCase) notify(ctx context.Context, userID uuid.UUID, message string) error {
	n := &notification.Notification{
		UserID:  userID,
		Message: message,
		Type:    notification.System,
		Status:  notification.Pending,
	}
	return uc.notfRepo.Create(ctx, n)
}
