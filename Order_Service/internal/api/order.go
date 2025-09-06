package api

import (
	"Order_Service/internal/jwt"
	"Order_Service/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
	order.Status = "pending"

	orderID, err := api.db.CreateOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создали событие о создании заказа
	if api.producer != nil {
		orderEvent := models.OrderEvent{
			EventType:   "order_created",
			OrderID:     orderID,
			ProductName: order.ProductName,
			ProductID:   order.ProductID,
			SupplierID:  order.SupplierID,
			ClientID:    order.ClientID,
			Amount:      order.Amount,
			Status:      order.Status,
			Timestamp:   time.Now(),
		}

		// Отправляем его в кафку
		if err := api.producer.PublishMessage("order-events", orderEvent); err != nil {
			log.Printf("Failed to publish order created event: %v", err)
		}
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

func (api *api) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "supplier" {
		http.Error(w, "Only suppliers can update order status", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	orderIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var updateRequest models.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validStatuses := []string{"pending", "confirmed", "processing", "shipped", "delivered", "cancelled"}
	isValidStatus := false
	for _, status := range validStatuses {
		if updateRequest.Status == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	order, err := api.db.GetOrderByID(orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if order.SupplierID != user.ID {
		http.Error(w, "You can only update status of your own orders", http.StatusForbidden)
		return
	}

	err = api.db.UpdateOrderStatus(orderID, updateRequest.Status)
	if err != nil {
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}

	// Создали новое событие с обновлением статуса
	if api.producer != nil {
		orderEvent := models.OrderEvent{
			EventType:   "order_status_updated",
			OrderID:     orderID,
			ProductName: order.ProductName,
			ProductID:   order.ProductID,
			SupplierID:  order.SupplierID,
			ClientID:    order.ClientID,
			Amount:      order.Amount,
			Status:      updateRequest.Status,
			Timestamp:   time.Now(),
		}

		// Отправили его в кафку
		if err := api.producer.PublishMessage("order-events", orderEvent); err != nil {
			log.Printf("Failed to publish order status updated event: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (api *api) validateUserToken(r *http.Request) (*models.User, error) {
	authHeader := r.Header.Get("Authorization")
	tokenString, err := jwt.ExtractTokenFromHeader(authHeader)
	if err != nil {
		return nil, err
	}

	claims, err := jwt.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	}

	return user, nil
}
