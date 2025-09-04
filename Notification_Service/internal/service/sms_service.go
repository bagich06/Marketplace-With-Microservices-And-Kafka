package service

import (
	"Notification_Service/internal/models"
	"fmt"
	"log"
)

// SMSService handles SMS notifications
type SMSService struct {
	apiKey    string
	apiSecret string
	sender    string
}

func NewSMSService(apiKey, apiSecret, sender string) *SMSService {
	return &SMSService{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		sender:    sender,
	}
}

// SendSMS sends an SMS notification
// In a real implementation, you would integrate with an SMS provider like Twilio, Nexmo, etc.
func (s *SMSService) SendSMS(notification models.SMSNotification) error {
	// Mock implementation - in production, integrate with actual SMS service
	log.Printf("Sending SMS to %s: %s", notification.PhoneNumber, notification.Message)

	// Here you would implement actual SMS sending logic
	// For example, using Twilio:
	// client := twilio.NewRestClient()
	// params := &api.CreateMessageParams{}
	// params.SetFrom(s.sender)
	// params.SetTo(notification.PhoneNumber)
	// params.SetBody(notification.Message)
	// resp, err := client.Api.CreateMessage(params)

	// For now, we'll just log the message
	log.Printf("SMS sent successfully to %s", notification.PhoneNumber)
	return nil
}

// CreateOrderNotificationSMS creates SMS notification for order events
func (s *SMSService) CreateOrderNotificationSMS(orderEvent models.OrderEvent, userInfo models.UserInfo) models.SMSNotification {
	var message string

	switch orderEvent.EventType {
	case "order_created":
		if userInfo.Role == "supplier" {
			message = fmt.Sprintf("Новый заказ #%d: %s. Проверьте ваш аккаунт для деталей.",
				orderEvent.OrderID, orderEvent.ProductName)
		} else {
			message = fmt.Sprintf("Заказ #%d оформлен: %s. Следите за обновлениями статуса.",
				orderEvent.OrderID, orderEvent.ProductName)
		}

	case "order_status_updated":
		message = fmt.Sprintf("Заказ #%d: статус изменен на '%s'",
			orderEvent.OrderID, orderEvent.Status)

	default:
		message = fmt.Sprintf("Уведомление о заказе #%d", orderEvent.OrderID)
	}

	return models.SMSNotification{
		PhoneNumber: userInfo.Phone,
		Message:     message,
	}
}
