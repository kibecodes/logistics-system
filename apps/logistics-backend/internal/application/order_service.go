package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	deliveryadapter "logistics-backend/internal/adapters/delivery"
	driveradapter "logistics-backend/internal/adapters/driver"
	inventoryadapter "logistics-backend/internal/adapters/inventory"
	orderadapter "logistics-backend/internal/adapters/order"
	useradapter "logistics-backend/internal/adapters/user"

	"github.com/google/uuid"
)

type OrderService struct {
	Users       *useradapter.UseCaseAdapter
	Orders      *orderadapter.UseCaseAdapter
	Drivers     *driveradapter.UseCaseAdapter
	Deliveries  *deliveryadapter.UseCaseAdapter
	Inventories *inventoryadapter.UseCaseAdapter
}

func NewOrderService(
	userUC *useradapter.UseCaseAdapter,
	orderUC *orderadapter.UseCaseAdapter,
	driverUC *driveradapter.UseCaseAdapter,
	deliveryUC *deliveryadapter.UseCaseAdapter,
	inventoryUC *inventoryadapter.UseCaseAdapter,
) *OrderService {
	return &OrderService{
		Users:       userUC,
		Orders:      orderUC,
		Drivers:     driverUC,
		Deliveries:  deliveryUC,
		Inventories: inventoryUC,
	}
}

// GetOrderWithDriver looks up delivery to find assigned driver
func (s *OrderService) GetOrderWithDriver(ctx context.Context, orderID uuid.UUID) (any, error) {
	// Step 1: Get the order
	order, err := s.Orders.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	// Step 2: Try to find delivery for this order
	delivery, err := s.Deliveries.GetDeliveryByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No delivery assigned yet
			return struct {
				Order    any
				Delivery any
				Driver   any
			}{
				Order:    order,
				Delivery: nil,
				Driver:   nil,
			}, nil
		}
		return nil, fmt.Errorf("get delivery: %w", err)
	}

	// Step 3: Fetch the driver from delivery
	driver, err := s.Drivers.GetDriverByID(ctx, delivery.DriverID)
	if err != nil {
		return nil, fmt.Errorf("get driver: %w", err)
	}

	return struct {
		Order    any
		Delivery any
		Driver   any
	}{Order: order, Delivery: delivery, Driver: driver}, nil
}

func (s *OrderService) UpdateOrderAndDriver(ctx context.Context, orderID uuid.UUID, driverID uuid.UUID, column string, value any) error {
	if err := s.Orders.UpdateOrder(ctx, orderID, column, value); err != nil {
		return fmt.Errorf("update order: %w", err)
	}
	if err := s.Drivers.UpdateDriverAvailability(ctx, driverID, "available", false); err != nil {
		return fmt.Errorf("update driver: %w", err)
	}
	return nil
}

func (s *OrderService) GetCustomersAndInventories(ctx context.Context) (any, error) {
	customers, err := s.Users.GetAllCustomers(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch customers: %w", err)
	}

	inventories, err := s.Inventories.GetAllInventories(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch inventories: %w", err)
	}

	return struct {
		Customers   any
		Inventories any
	}{
		Customers:   customers,
		Inventories: inventories,
	}, nil
}
