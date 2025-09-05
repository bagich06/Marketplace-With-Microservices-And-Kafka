package service

import (
	"Payment_Service/internal/kafka"
	"Payment_Service/internal/models"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-uuid"
)

type PaymentService struct {
	repo     PaymentRepository
	producer *kafka.Producer
}

type PaymentRepository interface {
	CreatePayment(payment models.Payment) error
	GetPaymentByID(paymentID string) (*models.Payment, error)
	GetPaymentsByClientID(clientID int) ([]models.Payment, error)
	GetPaymentByOrderID(orderID int) (*models.Payment, error)
	UpdatePaymentStatus(paymentID string, status models.PaymentStatus, transactionID, failureReason string) error
}

func NewPaymentService(repo PaymentRepository, producer *kafka.Producer) *PaymentService {
	return &PaymentService{
		repo:     repo,
		producer: producer,
	}
}

func (ps *PaymentService) HandleOrderEvent(event models.OrderEvent) error {
	log.Printf("Processing order event: %+v", event)

	switch event.EventType {
	case "order_created":
		return ps.handleOrderCreated(event)
	default:
		log.Printf("Unknown event type: %s", event.EventType)
	}

	return nil
}

func (ps *PaymentService) handleOrderCreated(event models.OrderEvent) error {
	paymentID, _ := uuid.GenerateUUID()

	payment := models.Payment{
		ID:            paymentID,
		OrderID:       event.OrderID,
		ClientID:      event.ClientID,
		Amount:        event.Amount,
		Status:        models.PaymentStatusPending,
		PaymentMethod: models.PaymentMethodCard,
		CreatedAt:     time.Now(),
	}

	if err := ps.repo.CreatePayment(payment); err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	paymentEvent := models.PaymentEvent{
		EventType: "payment_required",
		PaymentID: paymentID,
		OrderID:   event.OrderID,
		ClientID:  event.ClientID,
		Amount:    event.Amount,
		Status:    string(models.PaymentStatusPending),
		Timestamp: time.Now(),
	}

	if err := ps.producer.PublishMessage("order-events", paymentEvent); err != nil {
		log.Printf("Failed to publish payment_required event: %v", err)
	}

	log.Printf("Payment created for order %d: %s", event.OrderID, paymentID)
	return nil
}

func (ps *PaymentService) CreatePayment(request models.CreatePaymentRequest, clientID int) (*models.PaymentResponse, error) {
	existingPayment, err := ps.repo.GetPaymentByOrderID(request.OrderID)
	if err == nil && existingPayment != nil {
		return &models.PaymentResponse{
			ID:            existingPayment.ID,
			OrderID:       existingPayment.OrderID,
			Amount:        existingPayment.Amount,
			Status:        existingPayment.Status,
			PaymentMethod: existingPayment.PaymentMethod,
			CreatedAt:     existingPayment.CreatedAt,
			PaymentURL:    fmt.Sprintf("/api/payments/%s/pay", existingPayment.ID),
		}, nil
	}

	paymentID, _ := uuid.GenerateUUID()

	payment := models.Payment{
		ID:            paymentID,
		OrderID:       request.OrderID,
		ClientID:      clientID,
		Amount:        request.Amount,
		Status:        models.PaymentStatusPending,
		PaymentMethod: request.PaymentMethod,
		CreatedAt:     time.Now(),
	}

	if err := ps.repo.CreatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return &models.PaymentResponse{
		ID:            payment.ID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Status:        payment.Status,
		PaymentMethod: payment.PaymentMethod,
		CreatedAt:     payment.CreatedAt,
		PaymentURL:    fmt.Sprintf("/api/payments/%s/pay", payment.ID),
	}, nil
}

func (ps *PaymentService) GetPayment(paymentID string) (*models.Payment, error) {
	return ps.repo.GetPaymentByID(paymentID)
}

func (ps *PaymentService) ProcessPayment(paymentID string, request models.ProcessPaymentRequest) error {
	payment, err := ps.repo.GetPaymentByID(paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	if payment.Status != models.PaymentStatusPending {
		return fmt.Errorf("payment is not in pending status")
	}

	success := ps.mockProcessPayment(request)

	var newStatus models.PaymentStatus
	var transactionID string
	var failureReason string

	if success {
		newStatus = models.PaymentStatusCompleted
		transactionID = fmt.Sprintf("txn_%d", time.Now().Unix())
	} else {
		newStatus = models.PaymentStatusFailed
		failureReason = "Payment processing failed"
	}

	if err := ps.repo.UpdatePaymentStatus(paymentID, newStatus, transactionID, failureReason); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	paymentEvent := models.PaymentEvent{
		EventType: "payment_completed",
		PaymentID: paymentID,
		OrderID:   payment.OrderID,
		ClientID:  payment.ClientID,
		Amount:    payment.Amount,
		Status:    string(newStatus),
		Timestamp: time.Now(),
	}

	if err := ps.producer.PublishMessage("order-events", paymentEvent); err != nil {
		log.Printf("Failed to publish payment_completed event: %v", err)
	}

	return nil
}

func (ps *PaymentService) GetPaymentsByClient(clientID int) ([]models.Payment, error) {
	return ps.repo.GetPaymentsByClientID(clientID)
}

func (ps *PaymentService) mockProcessPayment(request models.ProcessPaymentRequest) bool {
	log.Printf("Processing payment with method: %s", request.PaymentMethod)

	if request.CardNumber != "" && len(request.CardNumber) >= 4 {
		return request.CardNumber[0:1] == "4"
	}

	return true
}
