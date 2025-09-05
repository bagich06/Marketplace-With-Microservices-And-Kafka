package service

import (
	"Notification_Service/internal/models"
	"fmt"
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
			subject = fmt.Sprintf("Новый заказ #%d", orderEvent.OrderID)
			body = "У вас новый заказ: " + orderEvent.ProductName
		} else {
			subject = fmt.Sprintf("Заказ #%d оформлен", orderEvent.OrderID)
			body = "Ваш заказ на " + orderEvent.ProductName + " успешно оформлен"
		}

	case "order_status_updated":
		subject = fmt.Sprintf("Обновление статуса заказа #%d", orderEvent.OrderID)
		body = "Статус заказа изменился на: " + orderEvent.Status

	default:
		subject = fmt.Sprintf("Уведомление о заказе #%d", orderEvent.OrderID)
		body = "Информация о заказе обновлена"
	}

	return models.EmailNotification{
		To:      userInfo.Email,
		Subject: subject,
		Body:    body,
		HTML:    false,
	}
}

func (m *MockEmailService) CreatePaymentCompletedNotificationEmail(orderEvent models.OrderEvent, userInfo models.UserInfo) models.EmailNotification {
	subject := fmt.Sprintf("Оплата заказа #%d успешно завершена", orderEvent.OrderID)
	body := fmt.Sprintf("Уважаемый %s,\n\nВаш заказ #%d успешно оплачен!\n\nСпасибо за покупку!",
		userInfo.Username, orderEvent.OrderID)

	return models.EmailNotification{
		To:      userInfo.Email,
		Subject: subject,
		Body:    body,
		HTML:    false,
	}
}

func (m *MockEmailService) CreatePaymentRequiredNotificationEmail(orderEvent models.OrderEvent, userInfo models.UserInfo) models.EmailNotification {
	subject := fmt.Sprintf("Требуется оплата заказа #%d", orderEvent.OrderID)
	body := fmt.Sprintf("Уважаемый %s,\n\nДля заказа #%d требуется оплата в размере %.2f руб.\n\nПожалуйста, перейдите к оплате.",
		userInfo.Username, orderEvent.OrderID, orderEvent.Amount)

	return models.EmailNotification{
		To:      userInfo.Email,
		Subject: subject,
		Body:    body,
		HTML:    false,
	}
}
