package repository

import (
	"Payment_Service/internal/models"
	"database/sql"
	"fmt"
	"time"
)

func (r *PGRepo) CreatePayment(payment models.Payment) error {
	query := `
		INSERT INTO payments (id, order_id, client_id, amount, status, payment_method, transaction_id, failure_reason, created_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(query,
		payment.ID,
		payment.OrderID,
		payment.ClientID,
		payment.Amount,
		payment.Status,
		payment.PaymentMethod,
		payment.TransactionID,
		payment.FailureReason,
		payment.CreatedAt,
		payment.CompletedAt,
	)

	return err
}

func (r *PGRepo) GetPaymentByID(paymentID string) (*models.Payment, error) {
	query := `
		SELECT id, order_id, client_id, amount, status, payment_method, transaction_id, failure_reason, created_at, completed_at
		FROM payments
		WHERE id = $1
	`

	var payment models.Payment
	var completedAt sql.NullTime

	err := r.db.QueryRow(query, paymentID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.ClientID,
		&payment.Amount,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.FailureReason,
		&payment.CreatedAt,
		&completedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, err
	}

	if completedAt.Valid {
		payment.CompletedAt = &completedAt.Time
	}

	return &payment, nil
}

func (r *PGRepo) GetPaymentsByClientID(clientID int) ([]models.Payment, error) {
	query := `
		SELECT id, order_id, client_id, amount, status, payment_method, transaction_id, failure_reason, created_at, completed_at
		FROM payments
		WHERE client_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		var completedAt sql.NullTime

		err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.ClientID,
			&payment.Amount,
			&payment.Status,
			&payment.PaymentMethod,
			&payment.TransactionID,
			&payment.FailureReason,
			&payment.CreatedAt,
			&completedAt,
		)
		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			payment.CompletedAt = &completedAt.Time
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PGRepo) GetPaymentByOrderID(orderID int) (*models.Payment, error) {
	query := `
		SELECT id, order_id, client_id, amount, status, payment_method, transaction_id, failure_reason, created_at, completed_at
		FROM payments
		WHERE order_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var payment models.Payment
	var completedAt sql.NullTime

	err := r.db.QueryRow(query, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.ClientID,
		&payment.Amount,
		&payment.Status,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.FailureReason,
		&payment.CreatedAt,
		&completedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found for order %d", orderID)
		}
		return nil, err
	}

	if completedAt.Valid {
		payment.CompletedAt = &completedAt.Time
	}

	return &payment, nil
}

func (r *PGRepo) UpdatePaymentStatus(paymentID string, status models.PaymentStatus, transactionID, failureReason string) error {
	query := `
		UPDATE payments 
		SET status = $1, transaction_id = $2, failure_reason = $3, completed_at = $4
		WHERE id = $5
	`

	var completedAt *time.Time
	if status == models.PaymentStatusCompleted {
		now := time.Now()
		completedAt = &now
	}

	_, err := r.db.Exec(query, status, transactionID, failureReason, completedAt, paymentID)
	return err
}
