package main

import (
	"Order_Service/internal/api"
	"Order_Service/internal/repository"
	"github.com/gorilla/mux"
	"log"
)

func main() {
	db, err := repository.NewPGRepo("postgres://postgres:postgres@localhost:5432/order_db")
	if err != nil {
		log.Fatal(err)
	}
	api := api.NewAPI(mux.NewRouter(), db)
	api.Handle()
	log.Fatal(api.ListenAndServe("localhost:8082"))
}
