package models

import "time"

type PaymentEvent struct {
	EventType   string    `json:"event_type"`
	PaymentID   string    `json:"payment_id"`
	OrderID     int       `json:"order_id"`
	ClientID    int       `json:"client_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}

type OrderEvent struct {
	EventType   string    `json:"event_type"`
	OrderID     int       `json:"order_id"`
	ProductName string    `json:"product_name"`
	ProductID   int       `json:"product_id"`
	SupplierID  int       `json:"supplier_id"`
	ClientID    int       `json:"client_id"`
	Status      string    `json:"status,omitempty"`
	Amount      float64   `json:"amount,omitempty"`
	PaymentID   string    `json:"payment_id,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Phone    string `json:"phone,omitempty"`
}
