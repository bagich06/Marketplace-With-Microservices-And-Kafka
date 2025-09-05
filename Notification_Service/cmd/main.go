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

	emailService := service.NewMockEmailService("noreply@marketplace.com")

	notificationService := service.NewNotificationService(emailService)

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

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	log.Println("Shutting down Notification Service...")
	cancel()
	log.Println("Notification Service stopped")
}
