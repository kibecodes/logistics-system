package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	deliveryadapter "logistics-backend/internal/adapters/delivery"
	driveradapter "logistics-backend/internal/adapters/driver"
	inventoryadapter "logistics-backend/internal/adapters/inventory"
	notificationadapter "logistics-backend/internal/adapters/notification"
	orderadapter "logistics-backend/internal/adapters/order"
	useradapter "logistics-backend/internal/adapters/user"
	"sort"

	"logistics-backend/internal/domain/driver"
	"logistics-backend/internal/domain/notification"
	order "logistics-backend/internal/domain/order"

	"github.com/google/uuid"
)

type Assignment struct {
	Order  *order.Order
	Driver *driver.Driver
}

type OrderService struct {
	Users         *useradapter.UseCaseAdapter
	Orders        *orderadapter.UseCaseAdapter
	Drivers       *driveradapter.UseCaseAdapter
	Deliveries    *deliveryadapter.UseCaseAdapter
	Inventories   *inventoryadapter.UseCaseAdapter
	Notifications *notificationadapter.UseCaseAdapter
}

func NewOrderService(
	userUC *useradapter.UseCaseAdapter,
	orderUC *orderadapter.UseCaseAdapter,
	driverUC *driveradapter.UseCaseAdapter,
	deliveryUC *deliveryadapter.UseCaseAdapter,
	inventoryUC *inventoryadapter.UseCaseAdapter,
	notificationUC *notificationadapter.UseCaseAdapter,
) *OrderService {
	return &OrderService{
		Users:         userUC,
		Orders:        orderUC,
		Drivers:       driverUC,
		Deliveries:    deliveryUC,
		Inventories:   inventoryUC,
		Notifications: notificationUC,
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

func (s *OrderService) OrderAssignment(ctx context.Context, maxDistance float64) ([]Assignment, error) {
	// 1. Fetch all pending orders
	allOrders, err := s.Orders.UseCase.ListOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch all orders failed: %w", err)
	}

	pendingOrders := filterPendingOrders(allOrders)
	if len(pendingOrders) == 0 {
		return nil, fmt.Errorf("no pending orders found")
	}

	// sort by oldest first
	sort.Slice(pendingOrders, func(i, j int) bool {
		return pendingOrders[i].CreatedAt.Before(pendingOrders[j].CreatedAt)
	})

	// 2. Get available drivers
	availableDrivers, err := s.Drivers.UseCase.ListAvailableDrivers(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("fetch available drivers failed: %w", err)
	}
	if len(availableDrivers) == 0 {
		return nil, fmt.Errorf("no available drivers at the moment")
	}

	// 3. Get total active deliveries (to gauge load)
	activeDeliveries, err := s.Deliveries.UseCase.ListActiveDeliveries(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch active deliveries failed: %w", err)
	}

	var assignments []Assignment
	driverCount := len(availableDrivers)
	deliveryLoadRatio := float64(len(activeDeliveries)) / float64(driverCount)

	// 4. Adaptive assignment strategy
	for i, o := range pendingOrders {
		var designatedDriver *driver.Driver

		if deliveryLoadRatio < 0.3 {
			// Early stage → assign sequentially to spread out first runs
			designatedDriver = availableDrivers[i%driverCount]
		} else {
			// Later stage → assign by proximity to pickup
			pickupPoint, err := s.Orders.UseCase.GetOrderPickupPoint(ctx, o.ID)
			if err != nil {
				continue
			}

			designatedDriver, err = s.Drivers.UseCase.GetClosestDriver(ctx, pickupPoint, maxDistance)
			if err != nil || designatedDriver == nil {
				continue
			}
		}

		// 5. Perform assignment via order domain
		if designatedDriver != nil {
			_, err := s.Orders.UseCase.AssignOrderToDriver(ctx, o.ID, designatedDriver.ID, maxDistance)
			if err != nil {
				continue
			}

			// Create a notification for both driver and customer
			msgCustomer := fmt.Sprintf("Your order %s has been assigned to driver %s", o.ID.String(), designatedDriver.FullName)
			msgDriver := fmt.Sprintf("You have been assigned a new delivery: order %s", o.ID.String())

			// Customer notification
			_ = s.Notifications.UseCase.CreateNotification(ctx, &notification.Notification{
				UserID:  o.CustomerID,
				Message: msgCustomer,
				Type:    notification.System, // or SMS later
				Status:  notification.Pending,
			})

			// Driver notification
			_ = s.Notifications.UseCase.CreateNotification(ctx, &notification.Notification{
				UserID:  designatedDriver.ID,
				Message: msgDriver,
				Type:    notification.System,
				Status:  notification.Pending,
			})

			// ✅ TODO (Future optimization):
			// Consider grouping or batching assignments by proximity and destination direction.
			// If multiple orders have similar pickup or delivery coordinates, compare them before assigning.
			// The goal: reduce redundant driver trips going to the same zone —
			// effectively load-balancing drivers by route proximity.
			//
			// This feature could support an "optimized dispatch mode" for admin control,
			// minimizing payout per trip by assigning clustered deliveries to fewer drivers.
			// e.g. two drivers can handle 5 nearby orders instead of sending 5 individually.
			//
			// Implementation idea:
			// - Use spatial clustering (K-Means / ST_ClusterWithin in PostGIS)
			// - Compute centroid of grouped destinations
			// - Adjust driver workloads dynamically based on cluster density and availability
			//
			// Feature name idea: "Proximity Load Shift" or "Smart Dispatch Optimization"

			assignments = append(assignments, Assignment{
				Order:  o,
				Driver: designatedDriver,
			})
		}
	}

	if len(assignments) == 0 {
		return nil, fmt.Errorf("no assignments were made")
	}

	return assignments, nil
}

func filterPendingOrders(orders []*order.Order) []*order.Order {
	var pending []*order.Order
	for _, o := range orders {
		if o.Status == order.Pending {
			pending = append(pending, o)
		}
	}
	return pending
}
