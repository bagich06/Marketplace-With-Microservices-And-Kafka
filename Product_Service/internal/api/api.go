package api

import (
	"Product_Service/internal/repository"
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
	api.r.HandleFunc("/api/product/create", api.CreateProductHandler)
	api.r.HandleFunc("/api/product/delete", api.DeleteProductHandler).Queries("id", "{id}")
	api.r.HandleFunc("/api/product/client", api.GetAllProductsForClientHandler)
	api.r.HandleFunc("/api/product/supplier", api.GetAllProductsForSupplierHandler)
}

func (api *api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, api.r)
}
