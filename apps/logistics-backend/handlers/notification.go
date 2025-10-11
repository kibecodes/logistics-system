package handlers

import (
	"encoding/json"
	"logistics-backend/internal/application"
	"logistics-backend/internal/domain/notification"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	UC *application.OrderService
}

func NewNotificationHandler(uc *application.OrderService) *NotificationHandler {
	return &NotificationHandler{UC: uc}
}

// CreateNotification godoc
// @Summary Create a new notification
// @Security JWT
// @Description Create a new notification with user_id, message, etc.
// @Tags notifications
// @Accept  json
// @Produce  json
// @Param user body notification.CreateNotificationRequest true "User Input"
// @Success 201 {object} notification.Notification
// @Failure 400 {string} handlers.ErrorResponse "Invalid request"
// @Failure 500 {string} handlers.ErrorResponse "Failed to create notification"
// @Router /notifications/create [post]
func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var req *notification.CreateNotificationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	n := req.ToNotification()

	if err := h.UC.Notifications.UseCase.CreateNotification(r.Context(), n); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not create notification", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":      n.ID,
		"user_id": n.UserID,
		"message": n.Message,
		"type":    n.Type,
		"sent_at": n.SentAt,
	})
}

// UpdateNotificationStatus godoc
// @Summary Update a notification's status (e.g. mark as sent or read)
// @Security JWT
// @Description Update notification status by ID
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Param status body notification.UpdateNotificationStatusRequest true "Status update request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handlers.ErrorResponse "Invalid request"
// @Failure 404 {object} handlers.ErrorResponse "Notification not found"
// @Failure 500 {object} handlers.ErrorResponse "Failed to update notification"
// @Router /notifications/{id}/status [patch]
func (h *NotificationHandler) UpdateNotificationStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid notification ID", err)
		return
	}

	var req notification.UpdateNotificationStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.UC.Notifications.UseCase.UpdateNotificationStatus(r.Context(), id, req.Status); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update status", err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Notification status updated successfully",
	})
}

// ListPendingNotifications godoc
// @Summary List all pending notifications
// @Security JWT
// @Description Get a list of notifications not yet sent
// @Tags notifications
// @Produce json
// @Success 200 {array} notification.Notification
// @Failure 500 {object} handlers.ErrorResponse
// @Router /notifications/all_pending_notifications [get]
func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	n, err := h.UC.Notifications.UseCase.ListPendingNotifications(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch notifications", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
}

// ListUserNotifications godoc
// @Summary List notifications by user
// @Security JWT
// @Description Get all notifications belonging to a user (optionally filtered by status)
// @Tags notifications
// @Produce json
// @Param id path string true "User ID"
// @Param status query string false "Notification status (pending, sent, failed, read)"
// @Success 200 {array} notification.Notification
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /notifications/all_my_notifications/{id} [get]
func (h *NotificationHandler) ListUserNotifications(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var n []*notification.Notification
	statusStr := r.URL.Query().Get("status")
	if statusStr != "" {
		status := notification.NotificationStatus(statusStr)
		n, err = h.UC.Notifications.UseCase.ListNotificationsByCustomer(r.Context(), userID, status)
	} else {
		n, err = h.UC.Notifications.UseCase.ListNotificationsByCustomer(r.Context(), userID, "pending")
	}

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch notifications", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
}

// MarkAsRead godoc
// @Summary Mark a single notification as read
// @Security JWT
// @Description Update notification status to "read" by ID
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handlers.ErrorResponse "Invalid ID"
// @Failure 500 {object} handlers.ErrorResponse "Failed to update notification"
// @Router /notifications/{id}/read [patch]
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	if err := h.UC.Notifications.UseCase.MarkAsRead(r.Context(), id); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "could not mark notification as read", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead godoc
// @Summary Mark all user notifications as read
// @Security JWT
// @Description Mark all unread notifications for a given user ID as "read"
// @Tags notifications
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handlers.ErrorResponse "Invalid user ID"
// @Failure 500 {object} handlers.ErrorResponse "Failed to mark notifications as read"
// @Router /notifications/mark_all_as_read/{id} [patch]
func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	if err := h.UC.Notifications.UseCase.MarkAllAsRead(r.Context(), userID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to mark notifications as read", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "All notifications marked as read",
	})
}
