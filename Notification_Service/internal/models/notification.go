package models

import "time"

type NotificationType string

const (
	NotificationTypeEmail NotificationType = "email"
	NotificationTypeSMS   NotificationType = "sms"
	NotificationTypePush  NotificationType = "push"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)

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

type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Phone    string `json:"phone,omitempty"`
}

type Notification struct {
	ID          string             `json:"id"`
	Type        NotificationType   `json:"type"`
	Recipient   string             `json:"recipient"`
	Subject     string             `json:"subject"`
	Message     string             `json:"message"`
	Status      NotificationStatus `json:"status"`
	OrderID     int                `json:"order_id"`
	UserID      int                `json:"user_id"`
	CreatedAt   time.Time          `json:"created_at"`
	SentAt      *time.Time         `json:"sent_at,omitempty"`
	FailureInfo string             `json:"failure_info,omitempty"`
}

type EmailNotification struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	HTML    bool   `json:"html"`
}

type SMSNotification struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}

type PushNotification struct {
	UserID  int                    `json:"user_id"`
	Title   string                 `json:"title"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
