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
}

type NotificationService struct {
	emailService EmailServiceInterface
	smsService   *SMSService
	pushService  *PushService
}

func NewNotificationService(emailService EmailServiceInterface, smsService *SMSService, pushService *PushService) *NotificationService {
	return &NotificationService{
		emailService: emailService,
		smsService:   smsService,
		pushService:  pushService,
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

	notificationID, _ := uuid.GenerateUUID()
	notification := models.Notification{
		ID:        notificationID,
		OrderID:   event.OrderID,
		UserID:    supplierInfo.ID,
		Status:    models.NotificationStatusPending,
		CreatedAt: time.Now(),
	}

	if supplierInfo.Email != "" {
		emailNotification := ns.emailService.CreateOrderNotificationEmail(event, *supplierInfo)
		notification.Type = models.NotificationTypeEmail
		notification.Recipient = supplierInfo.Email
		notification.Subject = emailNotification.Subject
		notification.Message = emailNotification.Body

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

	if supplierInfo.Phone != "" {
		smsNotification := ns.smsService.CreateOrderNotificationSMS(event, *supplierInfo)
		notification.Type = models.NotificationTypeSMS
		notification.Recipient = supplierInfo.Phone
		notification.Message = smsNotification.Message

		if err := ns.smsService.SendSMS(smsNotification); err != nil {
			notification.Status = models.NotificationStatusFailed
			notification.FailureInfo = err.Error()
			log.Printf("Failed to send SMS to supplier: %v", err)
		} else {
			notification.Status = models.NotificationStatusSent
			now := time.Now()
			notification.SentAt = &now
		}

		log.Printf("SMS notification created: %+v", notification)
	}

	pushNotification := ns.pushService.CreateOrderNotificationPush(event, *supplierInfo)
	notification.Type = models.NotificationTypePush
	notification.Recipient = fmt.Sprintf("user_%d", supplierInfo.ID)
	notification.Subject = pushNotification.Title
	notification.Message = pushNotification.Message

	if err := ns.pushService.SendPushNotification(pushNotification); err != nil {
		notification.Status = models.NotificationStatusFailed
		notification.FailureInfo = err.Error()
		log.Printf("Failed to send push notification to supplier: %v", err)
	} else {
		notification.Status = models.NotificationStatusSent
		now := time.Now()
		notification.SentAt = &now
	}

	log.Printf("Push notification created: %+v", notification)

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

	pushNotification := ns.pushService.CreateOrderNotificationPush(event, *clientInfo)
	if err := ns.pushService.SendPushNotification(pushNotification); err != nil {
		log.Printf("Failed to send push notification to client: %v", err)
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

	if clientInfo.Phone != "" {
		smsNotification := ns.smsService.CreateOrderNotificationSMS(event, *clientInfo)
		if err := ns.smsService.SendSMS(smsNotification); err != nil {
			log.Printf("Failed to send SMS to client: %v", err)
		}
	}

	pushNotification := ns.pushService.CreateOrderNotificationPush(event, *clientInfo)
	if err := ns.pushService.SendPushNotification(pushNotification); err != nil {
		log.Printf("Failed to send push notification to client: %v", err)
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
			Role:     "client", // Default role
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
