package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"logistics-backend/internal/domain/delivery"
	middleware "logistics-backend/internal/middleware"
	usecase "logistics-backend/internal/usecase/delivery"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DeliveryHandler struct {
	DH *usecase.UseCase
}

func NewDeliveryHandler(dh *usecase.UseCase) *DeliveryHandler {
	return &DeliveryHandler{DH: dh}
}

// CreateDelivery godoc
// @Summary Create a new delivery
// @Security JWT
// @Description Create a new delivery with order_id, driver_id, etc.
// @Tags deliveries
// @Accept  json
// @Produce  json
// @Param user body delivery.CreateDeliveryRequest true "User Input"
// @Success 201 {object} delivery.Delivery
// @Failure 400 {string} handlers.ErrorResponse "Invalid request"
// @Failure 500 {string} handlers.ErrorResponse "Failed to create delivery"
// @Router /deliveries/create [post]
func (dh *DeliveryHandler) CreateDelivery(w http.ResponseWriter, r *http.Request) {
	var req *delivery.CreateDeliveryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	d := req.ToDelivery()

	if err := dh.DH.CreateDelivery(r.Context(), d); err != nil {
		log.Printf("create delivery failed: %v", err)

		switch err {
		case delivery.ErrorNoPendingOrder:
			writeJSONError(w, http.StatusConflict, "No pending orders", err)
		default:
			writeJSONError(w, http.StatusInternalServerError, "could not create delivery", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{
		"id":           d.ID,
		"order_id":     d.OrderID,
		"driver_id":    d.DriverID,
		"assigned_at":  d.AssignedAt,
		"picked_up_at": d.PickedUpAt,
		"delivered_at": d.DeliveredAt,
		"status":       d.Status,
	})
}

// GetDeliveryByID godoc
// @Summary Get delivery by ID
// @Security JWT
// @Description Retrieve a delivery by their ID
// @Tags deliveries
// @Produce  json
// @Param id path string true "Delivery ID"
// @Success 200 {object} delivery.Delivery
// @Failure 400 {string} handlers.ErrorResponse "Invalid ID"
// @Failure 404 {string} handlers.ErrorResponse "Delivery not found"
// @Router /deliveries/by-id/{id} [get]
func (dh *DeliveryHandler) GetDeliveryByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	deliveryID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid delivery ID", nil)
		return
	}

	d, err := dh.DH.GetDeliveryByID(r.Context(), deliveryID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "No delivery found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// UpdateDelivery godoc
// @Summary Update Delivery
// @Security JWT
// @Description Update any delivery struct field of an existing delivery
// @Tags deliveries
// @Accept json
// @Produce json
// @Param delivery_id path string true "Delivery ID"
// @Param update body delivery.UpdateDeliveryRequest true "Field and value to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handlers.ErrorResponse "Invalid delivery ID or request body"
// @Failure 404 {object} handlers.ErrorResponse "Not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /deliveries/{id}/update [put]
func (dh *DeliveryHandler) UpdateDelivery(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	deliveryID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid delivery ID", nil)
		return
	}

	var req delivery.UpdateDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	column := strings.TrimSpace(strings.ToLower(req.Column))
	if column == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing or invalid column name", err)
		return
	}

	if err := dh.DH.UpdateDelivery(r.Context(), deliveryID, column, req.Value); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update delivery", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("delivery %s updated successfully", column),
	})

}

// ListDeliveries godoc
// @Summary List all deliveries
// @Security JWT
// @Description Get a list of all deliveries
// @Tags deliveries
// @Produce  json
// @Success 200 {array} delivery.Delivery
// @Router /deliveries/all_deliveries [get]
func (dh *DeliveryHandler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	deliveries, err := dh.DH.ListDeliveries(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch deliveries", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deliveries)
}

// DeleteDelivery godoc
// @Summary Delete a delivery
// @Description Permanently deletes a delivery by their ID
// @Tags deliveries
// @Accept json
// @Produce json
// @Security JWT
// @Param id path string true "Delivery ID"
// @Success 200 {object} map[string]string "Delivery deleted"
// @Failure 400 {object} handlers.ErrorResponse "Invalid Delivery ID"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /deliveries/{id} [delete]
func (dh *DeliveryHandler) DeleteDelivery(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	deliveryID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid delivery ID", nil)
		return
	}

	if err := dh.DH.DeleteDelivery(r.Context(), deliveryID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete delivery", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("delivery %s deleted", deliveryID),
	})
}

// AcceptDelivery godoc
// @Summary Accept and mark delivery as picked up
// @Description Updates delivery status to 'picked up' and sets the pickup timestamp. Only callable by authenticated drivers.
// @Tags deliveries
// @Security JWT
// @Produce json
// @Param id query string true "Delivery ID"
// @Success 200 {object} delivery.Delivery
// @Failure 400 {object} handlers.ErrorResponse "Missing or invalid delivery ID"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 500 {object} handlers.ErrorResponse "Failed to accept delivery"
// @Router /deliveries/{id}/accept [put]
func (dh *DeliveryHandler) AcceptDelivery(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing delivery ID", nil)
		return
	}

	deliveryID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid delivery ID", err)
		return
	}

	driverID, err := middleware.GetDriverIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	d := &delivery.Delivery{
		ID:         deliveryID,
		DriverID:   driverID,
		PickedUpAt: ptrTime(time.Now()),
		Status:     delivery.DeliveryPickedUp,
	}

	log.Printf("delivery d: %+v", d)
	if err := dh.DH.AcceptDelivery(r.Context(), d); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to accept delivery", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message":  fmt.Sprintf("delivery %s accepted", deliveryID),
		"delivery": d,
	})
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
