package repository

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type PGRepo struct {
	db *sql.DB
}

func NewPGRepo(connectionString string) (*PGRepo, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Создаем таблицу payments если не существует
	if err := createPaymentsTable(db); err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL database")
	return &PGRepo{db: db}, nil
}

func (r *PGRepo) Close() error {
	return r.db.Close()
}

func createPaymentsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS payments (
		id VARCHAR(36) PRIMARY KEY,
		order_id INTEGER NOT NULL,
		client_id INTEGER NOT NULL,
		amount DECIMAL(10,2) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'pending',
		payment_method VARCHAR(20) NOT NULL,
		transaction_id VARCHAR(100),
		failure_reason TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		completed_at TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
	CREATE INDEX IF NOT EXISTS idx_payments_client_id ON payments(client_id);
	CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
	`

	_, err := db.Exec(query)
	return err
}