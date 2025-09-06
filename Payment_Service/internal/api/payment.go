package api

import (
	"Payment_Service/internal/jwt"
	"Payment_Service/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *api) CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "client" {
		http.Error(w, "Only clients can create payments", http.StatusForbidden)
		return
	}

	var request models.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.OrderID <= 0 {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	if request.Amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	// Вызываем .CreatePayment из payment_service, который отвечает за отправку уведомлений в кафку
	response, err := api.paymentService.CreatePayment(request, user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *api) GetPaymentHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	paymentID, ok := vars["id"]
	if !ok {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	payment, err := api.paymentService.GetPayment(paymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if payment.ClientID != user.UserID && user.Role != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func (api *api) ProcessPaymentHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	paymentID, ok := vars["id"]
	if !ok {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	payment, err := api.paymentService.GetPayment(paymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if payment.ClientID != user.UserID && user.Role != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var request models.ProcessPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request.PaymentID = paymentID

	if request.PaymentMethod == models.PaymentMethodCard {
		if request.CardNumber == "" || request.ExpiryMonth == 0 || request.ExpiryYear == 0 || request.CVV == "" {
			http.Error(w, "Card details are required for card payment", http.StatusBadRequest)
			return
		}

		if len(request.CardNumber) < 13 || len(request.CardNumber) > 19 {
			http.Error(w, "Invalid card number length", http.StatusBadRequest)
			return
		}

		if len(request.CVV) < 3 || len(request.CVV) > 4 {
			http.Error(w, "Invalid CVV", http.StatusBadRequest)
			return
		}

		if request.ExpiryYear < 2024 || (request.ExpiryYear == 2024 && request.ExpiryMonth < 1) {
			http.Error(w, "Card has expired", http.StatusBadRequest)
			return
		}
	}

	// Вызывваем .ProcessPayment из payment_service, для того, чтобы выполнить платеж и отправить сообщение в кафку
	err = api.paymentService.ProcessPayment(paymentID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedPayment, err := api.paymentService.GetPayment(paymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Payment processed successfully",
		"payment": updatedPayment,
	})
}

func (api *api) GetPaymentsByClientHandler(w http.ResponseWriter, r *http.Request) {
	user, err := api.validateUserToken(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	clientIDStr, ok := vars["client_id"]
	if !ok {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		http.Error(w, "Invalid client ID format", http.StatusBadRequest)
		return
	}

	if clientID != user.UserID && user.Role != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	payments, err := api.paymentService.GetPaymentsByClient(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func (api *api) validateUserToken(r *http.Request) (*jwt.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	tokenString, err := jwt.ExtractTokenFromHeader(authHeader)
	if err != nil {
		return nil, err
	}

	claims, err := jwt.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
