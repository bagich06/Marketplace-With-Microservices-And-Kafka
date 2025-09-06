package api

import (
	"Payment_Service/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)

type api struct {
	router         *mux.Router
	paymentService *service.PaymentService
}

func NewAPI(router *mux.Router, paymentService *service.PaymentService) *api {
	return &api{
		router:         router,
		paymentService: paymentService,
	}
}

func (api *api) Handle() {
	api.router.HandleFunc("/api/payments", api.CreatePaymentHandler).Methods("POST")
	api.router.HandleFunc("/api/payments/{id}", api.GetPaymentHandler).Methods("GET")
	api.router.HandleFunc("/api/payments/{id}/pay", api.ProcessPaymentHandler).Methods("POST")
	api.router.HandleFunc("/api/payments/client/{client_id}", api.GetPaymentsByClientHandler).Methods("GET")
}

func (api *api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, api.router)
}
