package models

import "time"

type OrderEvent struct {
	EventType   string    `json:"event_type"`
	OrderID     int       `json:"order_id"`
	ProductName string    `json:"product_name"`
	ProductID   int       `json:"product_id"`
	SupplierID  int       `json:"supplier_id"`
	ClientID    int       `json:"client_id"`
	Status      string    `json:"status,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}
