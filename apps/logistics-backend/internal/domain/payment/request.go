package payment

import (
	"github.com/google/uuid"
)

type CreatePaymentRequest struct {
	OrderID  uuid.UUID     `json:"order_id"`
	Amount   int64         `json:"amount"`
	Currency string        `json:"currency"`
	Method   PaymentMethod `json:"method"`
	Status   PaymentStatus `json:"status"`
}

func (r *CreatePaymentRequest) ToPayment() *Payment {
	return &Payment{
		OrderID:  r.OrderID,
		Amount:   r.Amount,
		Currency: r.Currency,
		Method:   r.Method,
		Status:   StatusPending,
	}
}
