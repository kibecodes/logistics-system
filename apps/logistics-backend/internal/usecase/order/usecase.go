package order

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/order"
	"logistics-backend/internal/usecase/common"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      order.Repository
	invRepo   order.InventoryReader
	usrRepo   order.CustomerReader
	txManager common.TxManager
}

func NewUseCase(repo order.Repository, invRepo order.InventoryReader, usrRepo order.CustomerReader, txm common.TxManager) *UseCase {
	return &UseCase{repo: repo, invRepo: invRepo, usrRepo: usrRepo, txManager: txm}
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
		o.OrderStatus = order.Pending
		if err := uc.repo.Create(txCtx, o); err != nil {
			return fmt.Errorf("could not create order: %w", err)
		}

		return nil

		// TODO: goroutine
		// fire notification for order created
		// fire notification for restocking reminder
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
