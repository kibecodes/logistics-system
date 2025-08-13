package payment

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string
type PaymentStatus string

const (
	MethodStripe         PaymentMethod = "stripe"
	MethodPayPal         PaymentMethod = "paypal"
	MethodMobileMoney    PaymentMethod = "mobile_money"
	MethodCashOnDelivery PaymentMethod = "cash_on_delivery"

	StatusPending   PaymentStatus = "pending"
	StatusCompleted PaymentStatus = "completed"
	StatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID       uuid.UUID     `db:"id" json:"id"`
	OrderID  uuid.UUID     `db:"order_id" json:"order_id"`
	Amount   int64         `db:"amount" json:"amount"`
	Currency string        `db:"currency" json:"currency"`
	Method   PaymentMethod `db:"method" json:"method"`
	Status   PaymentStatus `db:"status" json:"status"`
	PaidAt   time.Time     `db:"paid_at" json:"paid_at"`
}
