package main

import (
	"Notification_Service/internal/kafka"
	"Notification_Service/internal/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting Notification Service...")

	// Initialize services
	// OPTION 1: Mock Email Service (для демо без реального SMTP)
	emailService := service.NewMockEmailService("noreply@marketplace.com")

	// OPTION 2: Real Email Service (раскомментируйте и настройте для реального использования)
	// emailService := service.NewEmailService(
	//     "smtp.gmail.com",           // SMTP host
	//     587,                        // SMTP port
	//     "your-email@gmail.com",     // SMTP username (замените на реальный)
	//     "your-app-password",        // SMTP password (замените на реальный App Password)
	//     "noreply@marketplace.com",  // From email
	// )

	smsService := service.NewSMSService(
		"your-sms-api-key",
		"your-sms-api-secret",
		"Marketplace",
	)

	pushService := service.NewPushService(
		"your-fcm-server-key",
		"/path/to/apns.pem",
	)

	notificationService := service.NewNotificationService(emailService, smsService, pushService)

	brokers := []string{"localhost:9092"}
	topics := []string{"order-events"}

	consumer := kafka.NewConsumer(brokers, topics, notificationService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Fatalf("Error starting Kafka consumer: %v", err)
		}
	}()

	log.Println("Notification Service started successfully")
	log.Println("Listening for Kafka messages on topics:", topics)

	// Wait for interrupt signal for graceful shutdown
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	log.Println("Shutting down Notification Service...")
	cancel()
	log.Println("Notification Service stopped")
}
