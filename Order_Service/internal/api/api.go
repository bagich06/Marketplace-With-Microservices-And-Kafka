package api

import (
	"Order_Service/internal/kafka"
	"Order_Service/internal/repository"
	"net/http"

	"github.com/gorilla/mux"
)

type api struct {
	r        *mux.Router
	db       *repository.PGRepo
	producer *kafka.Producer
}

func NewAPI(r *mux.Router, db *repository.PGRepo, producer *kafka.Producer) *api {
	return &api{r: r, db: db, producer: producer}
}

func (api *api) Handle() {
	api.r.HandleFunc("/api/order/create", api.CreateOrderHandler)
	api.r.HandleFunc("/api/orders/supplier", api.GetAllOrdersForSupplierHandler)
	api.r.HandleFunc("/api/order/client", api.GetAllOrdersForClientHandler)
	api.r.HandleFunc("/api/order/delete", api.DeleteOrderHandler).Queries("id", "{id}")
	api.r.HandleFunc("/api/order/status/{id}", api.UpdateOrderStatusHandler).Methods(http.MethodPut)
}

func (api *api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, api.r)
}
