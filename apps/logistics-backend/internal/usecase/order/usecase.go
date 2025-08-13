package order

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/order"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UseCase struct {
	repo    order.Repository
	invRepo order.InventoryReader
	usrRepo order.CustomerReader
}

func NewUseCase(repo order.Repository, invRepo order.InventoryReader, usrRepo order.CustomerReader) *UseCase {
	return &UseCase{repo: repo, invRepo: invRepo, usrRepo: usrRepo}
}

func (uc *UseCase) CreateOrder(ctx context.Context, o *order.Order) error {
	if o.Quantity <= 0 {
		return order.ErrorInvalidQuantity
	}

	inv, err := uc.invRepo.GetInventoryByID(ctx, o.InventoryID)
	if err != nil {
		return fmt.Errorf("could not fetch inventory: %w", err)
	}

	if inv.Stock < o.Quantity {
		return order.ErrorOutOfStock
	}

	// updating inventory table
	newStock := inv.Stock - o.Quantity
	if err := uc.invRepo.UpdateInventory(ctx, inv.ID, "stock", newStock); err != nil {
		return fmt.Errorf("could not update inventory stock: %w", err)
	}

	o.OrderStatus = order.Pending
	return uc.repo.Create(o)

	// TODO: goroutine
	// fire notification for order created
	// fire notification for restocking reminder
}

func (uc *UseCase) GetOrder(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	return uc.repo.GetByID(id)
}

func (uc *UseCase) GetOrderByCustomer(ctx context.Context, customerID uuid.UUID) ([]*order.Order, error) {
	return uc.repo.ListByCustomer(customerID)
}

// put method for simple update operation
func (uc *UseCase) UpdateOrder(ctx context.Context, orderID uuid.UUID, column string, value any) error {
	return uc.repo.UpdateColumn(ctx, orderID, column, value)
}

// put method for multiple operations - create + update
func (uc *UseCase) UpdateOrderTx(ctx context.Context, tx *sqlx.Tx, orderID uuid.UUID, column string, value any) error {
	return uc.repo.UpdateColumnTx(ctx, tx, orderID, column, value)
}

func (uc *UseCase) ListOrders(ctx context.Context) ([]*order.Order, error) {
	return uc.repo.List()
}

func (uc *UseCase) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *UseCase) GetAllInventories(ctx context.Context) ([]order.Inventory, error) {
	return uc.invRepo.GetAllInventories(ctx)
}

func (uc *UseCase) GetAllCustomers(ctx context.Context) ([]order.Customer, error) {
	return uc.usrRepo.GetAllCustomers(ctx)
}
