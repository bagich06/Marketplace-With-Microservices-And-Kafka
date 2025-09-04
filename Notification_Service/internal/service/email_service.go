package service

import (
	"Notification_Service/internal/models"
	"crypto/tls"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
}

func NewEmailService(smtpHost string, smtpPort int, smtpUsername, smtpPassword, fromEmail string) *EmailService {
	return &EmailService{
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
		fromEmail:    fromEmail,
	}
}

func (e *EmailService) SendEmail(notification models.EmailNotification) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.fromEmail)
	m.SetHeader("To", notification.To)
	m.SetHeader("Subject", notification.Subject)

	if notification.HTML {
		m.SetBody("text/html", notification.Body)
	} else {
		m.SetBody("text/plain", notification.Body)
	}

	d := gomail.NewDialer(e.smtpHost, e.smtpPort, e.smtpUsername, e.smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", notification.To, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", notification.To)
	return nil
}

func (e *EmailService) CreateOrderNotificationEmail(orderEvent models.OrderEvent, userInfo models.UserInfo) models.EmailNotification {
	var subject, body string

	switch orderEvent.EventType {
	case "order_created":
		if userInfo.Role == "supplier" {
			subject = fmt.Sprintf("Новый заказ #%d", orderEvent.OrderID)
			body = fmt.Sprintf(`
Здравствуйте, %s!

У вас новый заказ:
- Номер заказа: #%d
- Товар: %s
- ID товара: %d

Пожалуйста, обработайте заказ в кратчайшие сроки.

С уважением,
Команда Marketplace
			`, userInfo.Username, orderEvent.OrderID, orderEvent.ProductName, orderEvent.ProductID)
		} else {
			subject = fmt.Sprintf("Заказ #%d оформлен", orderEvent.OrderID)
			body = fmt.Sprintf(`
Здравствуйте, %s!

Ваш заказ успешно оформлен:
- Номер заказа: #%d
- Товар: %s
- ID товара: %d

Мы уведомим вас о изменении статуса заказа.

С уважением,
Команда Marketplace
			`, userInfo.Username, orderEvent.OrderID, orderEvent.ProductName, orderEvent.ProductID)
		}

	case "order_status_updated":
		subject = fmt.Sprintf("Обновление статуса заказа #%d", orderEvent.OrderID)
		body = fmt.Sprintf(`
Здравствуйте, %s!

Статус вашего заказа изменился:
- Номер заказа: #%d
- Товар: %s
- Новый статус: %s

С уважением,
Команда Marketplace
		`, userInfo.Username, orderEvent.OrderID, orderEvent.ProductName, orderEvent.Status)

	default:
		subject = fmt.Sprintf("Уведомление о заказе #%d", orderEvent.OrderID)
		body = fmt.Sprintf("Информация о заказе #%d обновлена", orderEvent.OrderID)
	}

	return models.EmailNotification{
		To:      userInfo.Email,
		Subject: subject,
		Body:    body,
		HTML:    false,
	}
}
