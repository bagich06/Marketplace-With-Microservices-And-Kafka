package repository

import (
	"Order_Service/internal/models"
	"context"
)

func (repo *PGRepo) CreateOrder(order models.Order) (int, error) {
	err := repo.pool.QueryRow(context.Background(),
		`INSERT INTO orders (product_name, product_id, supplier_id, client_id, amount, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID`,
		order.ProductName, order.ProductID, order.SupplierID, order.ClientID, order.Amount, "pending").Scan(&order.ID)
	if err != nil {
		return 0, err
	}
	return order.ID, nil
}

func (repo *PGRepo) GetAllOrdersByClientID(clientID int) ([]models.Order, error) {
	var orders []models.Order
	rows, err := repo.pool.Query(context.Background(), `SELECT id, product_name, product_id, supplier_id, client_id, amount, status FROM orders WHERE client_id = $1`, clientID)
	if err != nil {
		return orders, err
	}
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.ProductName, &order.ProductID, &order.SupplierID, &order.ClientID, &order.Amount, &order.Status)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (repo *PGRepo) GetAllOrdersBySupplierID(supplierID int) ([]models.Order, error) {
	var orders []models.Order
	rows, err := repo.pool.Query(context.Background(), `SELECT id, product_name, product_id, supplier_id, client_id, amount, status FROM orders WHERE supplier_id = $1`, supplierID)
	if err != nil {
		return orders, err
	}
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.ProductName, &order.ProductID, &order.SupplierID, &order.ClientID, &order.Amount, &order.Status)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (repo *PGRepo) DeleteOrderByID(id int, clientID int) error {
	_, err := repo.pool.Exec(context.Background(), `DELETE FROM orders WHERE id = $1 AND client_id = $2`, id, clientID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PGRepo) GetOrderByID(id int) (*models.Order, error) {
	var order models.Order
	err := repo.pool.QueryRow(context.Background(),
		`SELECT id, product_name, product_id, supplier_id, client_id, amount, status FROM orders WHERE id = $1`,
		id).Scan(&order.ID, &order.ProductName, &order.ProductID, &order.SupplierID, &order.ClientID, &order.Amount, &order.Status)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (repo *PGRepo) UpdateOrderStatus(id int, status string) error {
	_, err := repo.pool.Exec(context.Background(), `UPDATE orders SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		return err
	}
	return nil
}
