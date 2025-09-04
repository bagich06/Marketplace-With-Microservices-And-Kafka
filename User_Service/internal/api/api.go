package api

import (
	"User_Service/internal/repository"
	"net/http"

	"github.com/gorilla/mux"
)

type api struct {
	r  *mux.Router
	db *repository.PGRepo
}

func NewAPI(r *mux.Router, db *repository.PGRepo) *api {
	return &api{r: r, db: db}
}

func (api *api) Handle() {
	api.r.HandleFunc("/api/login", api.LoginHandler)
	api.r.HandleFunc("/api/register", api.RegisterHandler)
	api.r.HandleFunc("/api/user/{id}", api.GetUserByIDHandler).Methods(http.MethodGet)
}

func (api *api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, api.r)
}
