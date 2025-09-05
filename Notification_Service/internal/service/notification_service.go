package service

import (
	"Notification_Service/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-uuid"
)

type EmailServiceInterface interface {
	SendEmail(notification models.EmailNotification) error
	CreateOrderNotificationEmail(orderEvent models.OrderEvent, userInfo models.UserInfo) models.EmailNotification
	CreatePaymentRequiredNotificationEmail(orderEvent models.OrderEvent, userInfo models.UserInfo) models.EmailNotification
	CreatePaymentCompletedNotificationEmail(orderEvent models.OrderEvent, userInfo models.UserInfo) models.EmailNotification
}

type NotificationService struct {
	emailService EmailServiceInterface
}

func NewNotificationService(emailService EmailServiceInterface) *NotificationService {
	return &NotificationService{
		emailService: emailService,
	}
}

func (ns *NotificationService) HandleOrderEvent(event models.OrderEvent) error {
	log.Printf("Processing order event: %+v", event)

	switch event.EventType {
	case "order_created":
		if err := ns.sendOrderCreatedNotificationToSupplier(event); err != nil {
			log.Printf("Failed to send notification to supplier: %v", err)
		}

		if err := ns.sendOrderCreatedNotificationToClient(event); err != nil {
			log.Printf("Failed to send notification to client: %v", err)
		}

	case "order_status_updated":
		if err := ns.sendOrderStatusUpdateNotificationToClient(event); err != nil {
			log.Printf("Failed to send status update notification to client: %v", err)
		}

	case "payment_completed":
		if err := ns.sendPaymentCompletedNotificationToClient(event); err != nil {
			log.Printf("Failed to send payment completed notification to client: %v", err)
		}

	default:
		log.Printf("Unknown event type: %s", event.EventType)
	}

	return nil
}

func (ns *NotificationService) sendOrderCreatedNotificationToSupplier(event models.OrderEvent) error {
	supplierInfo, err := ns.getUserInfo(event.SupplierID)
	if err != nil {
		return fmt.Errorf("failed to get supplier info: %w", err)
	}

	if supplierInfo.Email != "" {
		emailNotification := ns.emailService.CreateOrderNotificationEmail(event, *supplierInfo)

		notificationID, _ := uuid.GenerateUUID()
		notification := models.Notification{
			ID:        notificationID,
			Type:      models.NotificationTypeEmail,
			OrderID:   event.OrderID,
			UserID:    supplierInfo.ID,
			Recipient: supplierInfo.Email,
			Subject:   emailNotification.Subject,
			Message:   emailNotification.Body,
			Status:    models.NotificationStatusPending,
			CreatedAt: time.Now(),
		}

		if err := ns.emailService.SendEmail(emailNotification); err != nil {
			notification.Status = models.NotificationStatusFailed
			notification.FailureInfo = err.Error()
			log.Printf("Failed to send email to supplier: %v", err)
		} else {
			notification.Status = models.NotificationStatusSent
			now := time.Now()
			notification.SentAt = &now
		}

		log.Printf("Email notification created: %+v", notification)
	}

	return nil
}

func (ns *NotificationService) sendOrderCreatedNotificationToClient(event models.OrderEvent) error {
	clientInfo, err := ns.getUserInfo(event.ClientID)
	if err != nil {
		return fmt.Errorf("failed to get client info: %w", err)
	}

	if clientInfo.Email != "" {
		emailNotification := ns.emailService.CreateOrderNotificationEmail(event, *clientInfo)
		if err := ns.emailService.SendEmail(emailNotification); err != nil {
			log.Printf("Failed to send email to client: %v", err)
		}
	}

	return nil
}

func (ns *NotificationService) sendOrderStatusUpdateNotificationToClient(event models.OrderEvent) error {
	clientInfo, err := ns.getUserInfo(event.ClientID)
	if err != nil {
		return fmt.Errorf("failed to get client info: %w", err)
	}

	if clientInfo.Email != "" {
		emailNotification := ns.emailService.CreateOrderNotificationEmail(event, *clientInfo)
		if err := ns.emailService.SendEmail(emailNotification); err != nil {
			log.Printf("Failed to send email to client: %v", err)
		}
	}

	return nil
}

func (ns *NotificationService) getUserInfo(userID int) (*models.UserInfo, error) {
	var userInfo models.UserInfo

	resp, err := http.Get(fmt.Sprintf("http://localhost:8081/api/user/%d", userID))
	if err != nil {
		log.Printf("User service unavailable, using mock data for user %d", userID)
		return &models.UserInfo{
			ID:       userID,
			Username: fmt.Sprintf("user%d", userID),
			Email:    fmt.Sprintf("user%d@example.com", userID),
			Role:     "client",
			Phone:    fmt.Sprintf("+79001%06d", userID),
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

func (ns *NotificationService) sendPaymentCompletedNotificationToClient(event models.OrderEvent) error {
	log.Printf("Sending payment completed notification for order %d", event.OrderID)

	userInfo, err := ns.getUserInfo(event.ClientID)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	notification := ns.emailService.CreatePaymentCompletedNotificationEmail(event, *userInfo)
	if err := ns.emailService.SendEmail(notification); err != nil {
		return fmt.Errorf("failed to send payment completed email: %w", err)
	}

	log.Printf("Payment completed notification sent to %s for order %d", userInfo.Email, event.OrderID)
	return nil
}

func (ns *NotificationService) HandlePaymentEvent(event models.PaymentEvent) error {
	log.Printf("Processing payment event: %+v", event)

	switch event.EventType {
	case "payment_required":
		if err := ns.sendPaymentRequiredNotificationToClient(event); err != nil {
			log.Printf("Failed to send payment required notification to client: %v", err)
		}
	case "payment_completed":
		if err := ns.sendPaymentCompletedNotificationToClientFromPaymentEvent(event); err != nil {
			log.Printf("Failed to send payment completed notification to client: %v", err)
		}
	default:
		log.Printf("Unknown payment event type: %s", event.EventType)
	}

	return nil
}

func (ns *NotificationService) sendPaymentCompletedNotificationToClientFromPaymentEvent(event models.PaymentEvent) error {
	log.Printf("Sending payment completed notification for order %d", event.OrderID)

	userInfo, err := ns.getUserInfo(event.ClientID)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	// Создаем OrderEvent из PaymentEvent для совместимости с существующим методом
	orderEvent := models.OrderEvent{
		EventType: "payment_completed",
		OrderID:   event.OrderID,
		ClientID:  event.ClientID,
		Amount:    event.Amount,
		Status:    event.Status,
		Timestamp: event.Timestamp,
	}

	notification := ns.emailService.CreatePaymentCompletedNotificationEmail(orderEvent, *userInfo)
	if err := ns.emailService.SendEmail(notification); err != nil {
		return fmt.Errorf("failed to send payment completed email: %w", err)
	}

	log.Printf("Payment completed notification sent to %s for order %d", userInfo.Email, event.OrderID)
	return nil
}

func (ns *NotificationService) sendPaymentRequiredNotificationToClient(event models.PaymentEvent) error {
	log.Printf("Sending payment required notification for order %d", event.OrderID)

	userInfo, err := ns.getUserInfo(event.ClientID)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	// Создаем OrderEvent из PaymentEvent для совместимости с существующим методом
	orderEvent := models.OrderEvent{
		EventType: "payment_required",
		OrderID:   event.OrderID,
		ClientID:  event.ClientID,
		Amount:    event.Amount,
		Status:    event.Status,
		Timestamp: event.Timestamp,
	}

	notification := ns.emailService.CreatePaymentRequiredNotificationEmail(orderEvent, *userInfo)
	if err := ns.emailService.SendEmail(notification); err != nil {
		return fmt.Errorf("failed to send payment required email: %w", err)
	}

	log.Printf("Payment required notification sent to %s for order %d", userInfo.Email, event.OrderID)
	return nil
}
