package main

import (
	"User_Service/internal/api"
	"User_Service/internal/repository"
	"github.com/gorilla/mux"
	"log"
)

func main() {
	db, err := repository.NewPGRepo("postgres://postgres:postgres@localhost:5432/user_db")
	if err != nil {
		log.Fatal(err)
	}
	api := api.NewAPI(mux.NewRouter(), db)
	api.Handle()
	log.Fatal(api.ListenAndServe("localhost:8081"))
}
