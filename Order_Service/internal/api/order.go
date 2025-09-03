package api

import (
	"Order_Service/internal/models"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (api *api) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "client" {
		http.Error(w, "Only clients can make an order", http.StatusForbidden)
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.ClientID = user.ID

	orderID, err := api.db.CreateOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderID)
}

func (api *api) GetAllOrdersForClientHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "client" {
		http.Error(w, "Only clients can get a list of orders", http.StatusForbidden)
		return
	}

	orders, err := api.db.GetAllOrdersByClientID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (api *api) GetAllOrdersForSupplierHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "supplier" {
		http.Error(w, "Only suppliers can get a list of orders", http.StatusForbidden)
		return
	}

	orders, err := api.db.GetAllOrdersBySupplierID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (api *api) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "client" {
		http.Error(w, "Only clients can delete an order", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	orderID, ok := vars["id"]
	if !ok {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(orderID)
	if err != nil {
		http.Error(w, "Invalid order id", http.StatusBadRequest)
		return
	}
	err = api.db.DeleteOrderByID(id, user.ID)
	if err != nil {
		http.Error(w, "Error deleting order", http.StatusInternalServerError)
		return
	}
}

func (api *api) validateUserToken(r *http.Request) (*models.User, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, errors.New("authorization header required")
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/api/validate", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid token")
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
