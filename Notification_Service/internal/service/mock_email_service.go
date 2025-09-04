package service

import (
	"Notification_Service/internal/models"
	"log"
)

type MockEmailService struct {
	fromEmail string
}

func NewMockEmailService(fromEmail string) *MockEmailService {
	return &MockEmailService{
		fromEmail: fromEmail,
	}
}

func (e *MockEmailService) SendEmail(notification models.EmailNotification) error {
	log.Printf("📧 [MOCK EMAIL] To: %s", notification.To)
	log.Printf("📧 [MOCK EMAIL] Subject: %s", notification.Subject)
	log.Printf("📧 [MOCK EMAIL] Body: %s", notification.Body)
	log.Printf("📧 [MOCK EMAIL] Email sent successfully!")

	return nil
}

func (e *MockEmailService) CreateOrderNotificationEmail(orderEvent models.OrderEvent, userInfo models.UserInfo) models.EmailNotification {
	var subject, body string

	switch orderEvent.EventType {
	case "order_created":
		if userInfo.Role == "supplier" {
			subject = "Новый заказ #" + string(rune(orderEvent.OrderID))
			body = "У вас новый заказ: " + orderEvent.ProductName
		} else {
			subject = "Заказ #" + string(rune(orderEvent.OrderID)) + " оформлен"
			body = "Ваш заказ на " + orderEvent.ProductName + " успешно оформлен"
		}

	case "order_status_updated":
		subject = "Обновление статуса заказа #" + string(rune(orderEvent.OrderID))
		body = "Статус заказа изменился на: " + orderEvent.Status

	default:
		subject = "Уведомление о заказе #" + string(rune(orderEvent.OrderID))
		body = "Информация о заказе обновлена"
	}

	return models.EmailNotification{
		To:      userInfo.Email,
		Subject: subject,
		Body:    body,
		HTML:    false,
	}
}
