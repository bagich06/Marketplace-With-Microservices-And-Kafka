package main

import (
	"Order_Service/internal/api"
	"Order_Service/internal/kafka"
	"Order_Service/internal/repository"
	"log"

	"github.com/gorilla/mux"
)

func main() {
	db, err := repository.NewPGRepo("postgres://postgres:postgres@localhost:5432/order_db")
	if err != nil {
		log.Fatal(err)
	}

	kafkaProducer, err := kafka.NewProducer([]string{"localhost:9092"})
	if err != nil {
		log.Printf("Failed to create Kafka producer: %v", err)
	}

	api := api.NewAPI(mux.NewRouter(), db, kafkaProducer)
	api.Handle()

	log.Println("Order Service started on :8083")
	log.Fatal(api.ListenAndServe("localhost:8083"))
}
