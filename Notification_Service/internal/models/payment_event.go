package models

import "time"

type PaymentEvent struct {
	EventType string    `json:"event_type"`
	PaymentID string    `json:"payment_id"`
	OrderID   int       `json:"order_id"`
	ClientID  int       `json:"client_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
