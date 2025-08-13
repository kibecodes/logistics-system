package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"logistics-backend/internal/application"
	"logistics-backend/internal/domain/order"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderHandler struct {
	UC *application.OrderService
}

func NewOrderHandler(uc *application.OrderService) *OrderHandler {
	return &OrderHandler{UC: uc}
}

// CreateOrder godoc
// @Summary Create a new order
// @Security JWT
// @Description Creates an order and returns the new object
// @Tags orders
// @Accept json
// @Produce json
// @Param order body order.CreateOrderRequest true "Order input"
// @Success 201 {object} order.Order
// @Failure 400 {string} handlers.ErrorResponse "Bad request"
// @Failure 500 {string} handlers.ErrorResponse "Internal server error"
// @Router /orders/create [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req order.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	o := req.ToOrder()

	if err := h.UC.Orders.UseCase.CreateOrder(r.Context(), o); err != nil {

		switch err {
		case order.ErrorOutOfStock:
			writeJSONError(w, http.StatusConflict, "Product is out of stock", err)
			return
		case order.ErrorInvalidQuantity:
			writeJSONError(w, http.StatusConflict, "Invalid Product Quantity", err)
		default:
			writeJSONError(w, http.StatusInternalServerError, "Could not create order", err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":                o.ID,
		"customer_id":       o.CustomerID,
		"inventory_id":      o.InventoryID,
		"quantity":          o.Quantity,
		"pickup_location":   o.PickupLocation,
		"delivery_location": o.DeliveryLocation,
		"order_status":      o.OrderStatus,
		"created_at":        o.CreatedAt,
		"updated_at":        o.UpdatedAt,
	})
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Security JWT
// @Description Fetch a single order using UUID
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} order.Order
// @Failure 400 {string} handlers.ErrorResponse "Invalid ID"
// @Failure 404 {string} handlers.ErrorResponse "Not found"
// @Router /orders/by-id/{id} [get]
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid order ID", nil)
		return
	}
	o, err := h.UC.Orders.UseCase.GetOrder(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Order not found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}

// GetOrderByCustomer godoc
// @Summary Get order by Customer ID
// @Security JWT
// @Description Fetch order(s) using Customer ID
// @Tags orders
// @Produce json
// @Param customer_id path string true "Customer ID"
// @Success 200 {object} []order.Order
// @Failure 400 {string} handlers.ErrorResponse "Invalid Customer ID"
// @Failure 404 {string} handlers.ErrorResponse "Not found"
// @Router /orders/by-customer/{customer_id} [get]
func (h *OrderHandler) GetOrderByCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "customer_id")
	customerID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid customer ID", nil)
		return
	}

	o, err := h.UC.Orders.UseCase.GetOrderByCustomer(r.Context(), customerID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "No orders found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}

// UpdateOrder godoc
// @Summary Update Order
// @Security JWT
// @Description Update any order struct field of an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param update body order.UpdateOrderRequest true "Field and value to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handlers.ErrorResponse "Invalid order ID or request body"
// @Failure 404 {object} handlers.ErrorResponse "Not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /orders/{id}/update [put]
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid order ID", nil)
		return
	}

	var req order.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	column := strings.TrimSpace(strings.ToLower(req.Column))
	if column == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing or invalid column name", err)
		return
	}

	if err := h.UC.Orders.UseCase.UpdateOrder(r.Context(), orderID, column, req.Value); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update order", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("order %s updated successfully", column),
	})
}

// ListOrders godoc
// @Summary List all orders
// @Security JWT
// @Description Get a list of all orders
// @Tags orders
// @Produce  json
// @Success 200 {array} order.Order
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /orders/all_orders [get]
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.UC.Orders.UseCase.ListOrders(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch orders", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// DeleteOrder godoc
// @Summary Delete a order
// @Description Permanently deletes an order by their ID
// @Tags orders
// @Accept json
// @Produce json
// @Security JWT
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string "Order deleted"
// @Failure 400 {object} handlers.ErrorResponse "Invalid order ID"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid order ID", nil)
		return
	}

	if err := h.UC.Orders.UseCase.DeleteOrder(r.Context(), orderID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete order", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("order %s deleted", orderID),
	})
}

// GetOrderFormData godoc
// @Summary      Get data for order form dropdowns
// @Description  Returns a list of customers and inventories for populating order form dropdowns.
// @Tags         orders
// @Produce      json
// @Success      200  {object}  order.DropdownDataRequest
// @Failure      500  {object}  ErrorResponse "Failed to fetch customers or inventories"
// @Router       /orders/form-data [get]
func (h *OrderHandler) GetOrderFormData(w http.ResponseWriter, r *http.Request) {
	data, err := h.UC.GetCustomersAndInventories(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch form data", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
