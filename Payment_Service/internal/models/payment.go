package models

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type PaymentMethod string

const (
	PaymentMethodCard PaymentMethod = "card"
)

type Payment struct {
	ID            string        `json:"id"`
	OrderID       int           `json:"order_id"`
	ClientID      int           `json:"client_id"`
	Amount        float64       `json:"amount"`
	Status        PaymentStatus `json:"status"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	CreatedAt     time.Time     `json:"created_at"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty"`
	FailureReason string        `json:"failure_reason,omitempty"`
	TransactionID string        `json:"transaction_id,omitempty"`
}

type CreatePaymentRequest struct {
	OrderID       int           `json:"order_id"`
	Amount        float64       `json:"amount"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

type ProcessPaymentRequest struct {
	PaymentID      string        `json:"payment_id"`
	PaymentMethod  PaymentMethod `json:"payment_method"`
	CardNumber     string        `json:"card_number,omitempty"`
	ExpiryMonth    int           `json:"expiry_month,omitempty"`
	ExpiryYear     int           `json:"expiry_year,omitempty"`
	CVV            string        `json:"cvv,omitempty"`
	CardholderName string        `json:"cardholder_name,omitempty"`
}

type PaymentResponse struct {
	ID            string        `json:"id"`
	OrderID       int           `json:"order_id"`
	Amount        float64       `json:"amount"`
	Status        PaymentStatus `json:"status"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	CreatedAt     time.Time     `json:"created_at"`
	PaymentURL    string        `json:"payment_url,omitempty"`
}
