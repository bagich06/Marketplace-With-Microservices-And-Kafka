package service

import (
	"Notification_Service/internal/models"
	"fmt"
	"log"
)

// PushService handles push notifications
type PushService struct {
	fcmServerKey string
	apnsCertPath string
}

func NewPushService(fcmServerKey, apnsCertPath string) *PushService {
	return &PushService{
		fcmServerKey: fcmServerKey,
		apnsCertPath: apnsCertPath,
	}
}

// SendPushNotification sends a push notification
// In a real implementation, you would integrate with FCM (Firebase Cloud Messaging) for Android
// and APNs (Apple Push Notification service) for iOS
func (p *PushService) SendPushNotification(notification models.PushNotification) error {
	// Mock implementation - in production, integrate with FCM/APNs
	log.Printf("Sending push notification to user %d: %s - %s",
		notification.UserID, notification.Title, notification.Message)

	// Here you would implement actual push notification sending
	// For FCM:
	// client, err := fcm.NewClient(p.fcmServerKey)
	// if err != nil {
	//     return err
	// }
	//
	// message := &fcm.Message{
	//     To: userDeviceToken,
	//     Data: notification.Data,
	//     Notification: &fcm.Notification{
	//         Title: notification.Title,
	//         Body:  notification.Message,
	//     },
	// }
	//
	// response, err := client.Send(message)

	// For now, we'll just log the notification
	log.Printf("Push notification sent successfully to user %d", notification.UserID)
	return nil
}

// CreateOrderNotificationPush creates push notification for order events
func (p *PushService) CreateOrderNotificationPush(orderEvent models.OrderEvent, userInfo models.UserInfo) models.PushNotification {
	var title, message string
	data := map[string]interface{}{
		"order_id":   orderEvent.OrderID,
		"product_id": orderEvent.ProductID,
		"event_type": orderEvent.EventType,
		"user_role":  userInfo.Role,
	}

	switch orderEvent.EventType {
	case "order_created":
		if userInfo.Role == "supplier" {
			title = "Новый заказ"
			message = fmt.Sprintf("У вас новый заказ #%d на товар '%s'",
				orderEvent.OrderID, orderEvent.ProductName)
			data["action"] = "view_order"
		} else {
			title = "Заказ оформлен"
			message = fmt.Sprintf("Ваш заказ #%d на товар '%s' успешно оформлен",
				orderEvent.OrderID, orderEvent.ProductName)
			data["action"] = "track_order"
		}

	case "order_status_updated":
		title = "Обновление заказа"
		message = fmt.Sprintf("Статус заказа #%d изменен на '%s'",
			orderEvent.OrderID, orderEvent.Status)
		data["status"] = orderEvent.Status
		data["action"] = "view_order"

	default:
		title = "Уведомление о заказе"
		message = fmt.Sprintf("Обновление по заказу #%d", orderEvent.OrderID)
		data["action"] = "view_order"
	}

	return models.PushNotification{
		UserID:  userInfo.ID,
		Title:   title,
		Message: message,
		Data:    data,
	}
}
