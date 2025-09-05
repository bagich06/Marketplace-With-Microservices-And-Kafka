package main

import (
	"Payment_Service/internal/api"
	"Payment_Service/internal/kafka"
	"Payment_Service/internal/repository"
	"Payment_Service/internal/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Payment Service...")

	// Подключение к базе данных
	db, err := repository.NewPGRepo("postgres://postgres:password@localhost:5432/marketplace")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Создание Kafka producer
	kafkaProducer, err := kafka.NewProducer([]string{"localhost:9092"})
	if err != nil {
		log.Printf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Создание сервиса
	paymentService := service.NewPaymentService(db, kafkaProducer)

	// Создание API
	router := mux.NewRouter()
	api := api.NewAPI(router, paymentService)
	api.Handle()

	// Запуск Kafka consumer в отдельной горутине
	brokers := []string{"localhost:9092"}
	topics := []string{"order-events"}

	consumer := kafka.NewConsumer(brokers, topics, paymentService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Fatalf("Error starting Kafka consumer: %v", err)
		}
	}()

	// Запуск HTTP сервера
	go func() {
		log.Println("Payment Service started on :8085")
		if err := api.ListenAndServe("localhost:8085"); err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	log.Println("Payment Service started successfully")
	log.Println("Listening for Kafka messages on topics:", topics)

	// Ожидание сигнала завершения
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	log.Println("Shutting down Payment Service...")
	cancel()
	log.Println("Payment Service stopped")
}