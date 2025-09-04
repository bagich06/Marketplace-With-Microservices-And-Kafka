package main

import (
	"Product_Service/internal/api"
	"Product_Service/internal/repository"
	"github.com/gorilla/mux"
	"log"
)

func main() {
	db, err := repository.NewPGRepo("postgres://postgres:postgres@localhost:5432/product_db")
	if err != nil {
		log.Fatal(err)
	}
	api := api.NewAPI(mux.NewRouter(), db)
	api.Handle()
	log.Fatal(api.ListenAndServe("localhost:8082"))
}
