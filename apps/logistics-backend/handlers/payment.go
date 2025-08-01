package handlers

import (
	"encoding/json"
	"logistics-backend/internal/domain/payment"
	usecase "logistics-backend/internal/usecase/payment"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	PH *usecase.UseCase
}

func NewPaymentHandler(ph *usecase.UseCase) *PaymentHandler {
	return &PaymentHandler{PH: ph}
}

// CreatePayment godoc
// @Summary Create new payment
// @Security JWT
// @Description Register new payment with order_id, amount, etc.
// @Tags payments
// @Accept  json
// @Produce  json
// @Param user body payment.CreatePaymentRequest true "User Input"
// @Success 201 {object} payment.Payment
// @Failure 400 {string} handlers.ErrorResponse "Invalid request"
// @Failure 500 {string} handlers.ErrorResponse "Failed to create payment"
// @Router /payments/create [post]
func (ph *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req *payment.CreatePaymentRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	p := req.ToPayment()

	if err := ph.PH.CreatePayment(r.Context(), p); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not create payment", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":       p.ID,
		"order_id": p.OrderID,
		"amount":   p.Amount,
		"method":   p.Method,
		"status":   p.Status,
		"paid_at":  p.PaidAt,
	})
}

// GetPaymentByID godoc
// @Summary Get payment by ID
// @Security JWT
// @Description Fetch a single payment using payment ID
// @Tags payments
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} payment.Payment
// @Failure 400 {string} handlers.ErrorResponse "Invalid ID"
// @Failure 404 {string} handlers.ErrorResponse "Not found"
// @Router /payments/{id} [get]
func (ph *PaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	paymentID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid payment ID", nil)
		return
	}

	p, err := ph.PH.GetPaymentByID(r.Context(), paymentID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "No payment found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// GetPaymentByOrder godoc
// @Summary Get payment by Order ID
// @Security JWT
// @Description Fetch payment(s) using Order ID
// @Tags payments
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} []payment.Payment
// @Failure 400 {string} handlers.ErrorResponse "Invalid Order ID"
// @Failure 404 {string} handlers.ErrorResponse "Not found"
// @Router /payments/{order_id} [get]
func (ph *PaymentHandler) GetPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "order_id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid order ID", nil)
		return
	}

	p, err := ph.PH.GetPaymentByOrderID(r.Context(), orderID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "No payment found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// ListPayments godoc
// @Summary List all payments
// @Security JWT
// @Description Get a list of all payments
// @Tags payments
// @Produce  json
// @Success 200 {array} payment.Payment
// @Router /payments/all_payments [get]
func (ph *PaymentHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := ph.PH.ListPayments(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch payments", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}
