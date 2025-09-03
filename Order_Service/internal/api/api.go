package api

import (
	"Order_Service/internal/repository"
	"github.com/gorilla/mux"
	"net/http"
)

type api struct {
	r  *mux.Router
	db *repository.PGRepo
}

func NewAPI(r *mux.Router, db *repository.PGRepo) *api {
	return &api{r: r, db: db}
}

func (api *api) Handle() {
	api.r.HandleFunc("/api/order/create", api.CreateOrderHandler)
	api.r.HandleFunc("/api/orders/supplier", api.GetAllOrdersForSupplierHandler)
	api.r.HandleFunc("/api/order/client", api.GetAllOrdersForClientHandler)
	api.r.HandleFunc("/api/order/delete", api.DeleteOrderHandler).Queries("id", "{id}")
}

func (api *api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, api.r)
}
