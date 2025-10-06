package handlers

import (
	"encoding/json"
	"fmt"
	"logistics-backend/internal/application"
	"logistics-backend/internal/domain/driver"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DriverHandler struct {
	UC *application.OrderService
}

func NewDriverHandler(uc *application.OrderService) *DriverHandler {
	return &DriverHandler{UC: uc}
}

// UpdateDriverProfile godoc
// @Summary Update driver profile
// @Description Updates the vehicle information and current location of a driver
// @Tags drivers
// @Security JWT
// @Accept json
// @Produce json
// @Param driver_id path string true "Driver ID"
// @Param body body driver.UpdateDriverRequest true "Driver profile fields to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /drivers/{id}/profile [patch]
func (h *DriverHandler) UpdateDriverProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	driverID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid driver ID", nil)
		return
	}

	var req driver.UpdateDriverProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.UC.Drivers.UseCase.UpdateDriverProfile(r.Context(), driverID, &req); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update driver profile", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Driver profile updated"})

}

// UpdateDriver godoc
// @Summary Update a specific driver field
// @Description Updates a driver's specific field (e.g., VehicleInfo, CurrentLocation) based on driver ID
// @Tags drivers
// @Accept json
// @Produce json
// @Security JWT
// @Param id path string true "Driver ID"
// @Param data body driver.UpdateDriverRequest true "Field and value to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handlers.ErrorResponse "Invalid driver ID or request body"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /drivers/{id}/update [put]
func (h *DriverHandler) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	driverID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid driver ID", nil)
		return
	}

	var req driver.UpdateDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.UC.Drivers.UseCase.UpdateDriver(r.Context(), driverID, &req); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update driver", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("driver %s updated successfully", req.Column),
	})
}

// GetDriverByID godoc
// @Summary Get driver by ID
// @Security JWT
// @Description Retrieve a driver by their ID
// @Tags drivers
// @Produce  json
// @Param id path string true "Driver ID"
// @Success 200 {object} driver.Driver
// @Failure 400 {string} handlers.ErrorResponse "Invalid ID"
// @Failure 404 {string} handlers.ErrorResponse "Driver not found"
// @Router /drivers/by-id/{id} [get]
func (h *DriverHandler) GetDriverByID(w http.ResponseWriter, r *http.Request) {
	driverID := chi.URLParam(r, "id")
	id, err := uuid.Parse(driverID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid driver ID", nil)
		return
	}

	d, err := h.UC.Drivers.UseCase.GetDriver(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Driver not found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// GetUserByDriver godoc
// @Summary Get driver by Email
// @Security JWT
// @Description Retrieve a driver by their Email
// @Tags drivers
// @Produce  json
// @Param email path string true "Driver Email"
// @Success 200 {object} driver.Driver
// @Failure 400 {string} handlers.ErrorResponse "Invalid Email"
// @Failure 404 {string} handlers.ErrorResponse "Driver not found"
// @Router /drivers/by-email/{email} [get]
func (h *DriverHandler) GetDriverByEmail(w http.ResponseWriter, r *http.Request) {
	emailParam := chi.URLParam(r, "email")
	email, err := url.PathUnescape(emailParam)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid email format", nil)
		return
	}

	d, err := h.UC.Drivers.UseCase.GetDriverByEmail(r.Context(), email)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Driver not found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// ListDrivers godoc
// @Summary List all drivers
// @Security JWT
// @Description Get a list of all registered drivers
// @Tags drivers
// @Produce  json
// @Success 200 {array} driver.Driver
// @Router /drivers/all_drivers [get]
func (h *DriverHandler) ListDrivers(w http.ResponseWriter, r *http.Request) {
	drivers, err := h.UC.Drivers.UseCase.ListDrivers(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch drivers", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}

// DeleteDriver godoc
// @Summary Delete a driver
// @Description Permanently deletes a driver by their ID
// @Tags drivers
// @Accept json
// @Produce json
// @Security JWT
// @Param id path string true "Driver ID"
// @Success 200 {object} map[string]string "Driver deleted"
// @Failure 400 {object} handlers.ErrorResponse "Invalid driver ID"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /drivers/{id} [delete]
func (h *DriverHandler) DeleteDriver(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	driverID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid driver ID", nil)
		return
	}

	if err := h.UC.Drivers.UseCase.DeleteDriver(r.Context(), driverID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("driver %s deleted", driverID),
	})
}
