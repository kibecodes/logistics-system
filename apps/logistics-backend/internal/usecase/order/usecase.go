package order

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/driver"
	"logistics-backend/internal/domain/notification"
	"logistics-backend/internal/domain/order"
	"logistics-backend/internal/usecase/common"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
)

type UseCase struct {
	repo      order.Repository
	invRepo   order.InventoryReader
	usrRepo   order.CustomerReader
	drvRepo   order.DriverReader
	txManager common.TxManager
	notfRepo  order.NotificationReader
}

func NewUseCase(repo order.Repository, invRepo order.InventoryReader, usrRepo order.CustomerReader, txm common.TxManager, notf order.NotificationReader) *UseCase {
	return &UseCase{repo: repo, invRepo: invRepo, usrRepo: usrRepo, txManager: txm, notfRepo: notf}
}

func (uc *UseCase) CreateOrder(ctx context.Context, o *order.Order) (err error) {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		// 1. check order quantity
		if o.Quantity <= 0 {
			return order.ErrorInvalidQuantity
		}

		// 2. get inventory
		inv, err := uc.invRepo.GetInventoryByID(txCtx, o.InventoryID)
		if err != nil {
			return fmt.Errorf("could not fetch inventory: %w", err)
		}

		// 3. check available stock
		if inv.Stock < o.Quantity {
			return order.ErrorOutOfStock
		}

		// 4. updating inventory table
		newStock := inv.Stock - o.Quantity
		if err := uc.invRepo.UpdateInventory(txCtx, inv.ID, "stock", newStock); err != nil {
			return fmt.Errorf("could not update inventory stock: %w", err)
		}

		// 5. create order with status = pending
		o.Status = order.Pending
		if err := uc.repo.Create(txCtx, o); err != nil {
			return fmt.Errorf("could not create order: %w", err)
		}

		// 6. Fire notifications (async, after commit)
		go func() {
			// a. Notify customer
			msgCustomer := fmt.Sprintf("Your order %s has been created successfully.", o.ID)
			_ = uc.notify(ctx, o.CustomerID, msgCustomer)

			// b. Notify admin if stock low
			if newStock <= 5 { // example threshold
				msgAdmin := fmt.Sprintf("⚠️ Inventory '%s' stock is low: only %d left.", inv.Name, newStock)
				_ = uc.notify(ctx, inv.AdminID, msgAdmin)
			}
		}()

		return nil
	})
}

func (uc *UseCase) GetOrder(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetOrderByCustomer(ctx context.Context, customerID uuid.UUID) ([]*order.Order, error) {
	return uc.repo.ListByCustomer(ctx, customerID)
}

func (uc *UseCase) UpdateOrder(ctx context.Context, orderID uuid.UUID, column string, value any) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Update(txCtx, orderID, column, value); err != nil {
			return fmt.Errorf("update order failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) ListOrders(ctx context.Context) ([]*order.Order, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete order failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetAllInventories(ctx context.Context) ([]order.Inventory, error) {
	return uc.invRepo.GetAllInventories(ctx)
}

func (uc *UseCase) GetAllCustomers(ctx context.Context) ([]order.Customer, error) {
	return uc.usrRepo.GetAllCustomers(ctx)
}

func (uc *UseCase) GetOrderPickupPoint(ctx context.Context, orderID uuid.UUID) (postgis.PointS, error) {
	return uc.repo.GetPickupPoint(ctx, orderID)
}

func (uc *UseCase) GetOrderDeliveryPoint(ctx context.Context, orderID uuid.UUID) (postgis.PointS, error) {
	return uc.repo.GetDeliveryPoint(ctx, orderID)
}

func (uc *UseCase) AssignOrderToDriver(ctx context.Context, orderID, driverID uuid.UUID, maxDistance float64) (*driver.Driver, error) {
	// 1. Fetch the order
	o, err := uc.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("fetch order: %w", err)
	}

	if o.Status != order.Pending {
		return nil, fmt.Errorf("order not pending")
	}

	// 2. Get pickup point
	pickupPoint, err := uc.GetOrderPickupPoint(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get pickup point: %w", err)
	}

	// 3. Find nearest available driver within maxDistance
	nearestDriver, err := uc.drvRepo.GetNearestDriver(ctx, pickupPoint, maxDistance)
	if err != nil || nearestDriver == nil {
		return nil, fmt.Errorf("no available driver within %.2f meters", maxDistance)
	}

	// 4. Update order status to assigned
	if err := uc.UpdateOrder(ctx, orderID, "status", order.Assigned); err != nil {
		return nil, fmt.Errorf("update order status: %w", err)
	}

	// 5. Optionally: create a pending assignment record
	// This could be a lightweight table: order_id, driver_id, assigned_at
	// You can also push a websocket / notification here to the driver
	// if err := uc.repo.CreatePendingAssignment(ctx, orderID, nearestDriver.ID); err != nil {
	//     return nil, fmt.Errorf("create pending assignment: %w", err)
	// }

	// 6. Return assigned driver for confirmation / logging
	return nearestDriver, nil
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
