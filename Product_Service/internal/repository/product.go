package repository

import (
	"Product_Service/internal/models"
	"context"
)

func (repo *PGRepo) CreateProduct(product models.Product) (int, error) {
	err := repo.pool.QueryRow(context.Background(), `INSERT INTO products (name, description, price, user_id) VALUES ($1, $2, $3, $4) RETURNING id`, product.Name, product.Description, product.Price, product.UserID).Scan(&product.ID)
	if err != nil {
		return 0, err
	}
	return product.ID, nil
}

func (repo *PGRepo) GetProductByID(id int) (models.Product, error) {
	var product models.Product
	err := repo.pool.QueryRow(context.Background(), `SELECT id, name, description, price FROM products WHERE id=$1`, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	return product, err
}

func (repo *PGRepo) GetAllProductsForClient() ([]models.Product, error) {
	var products []models.Product
	rows, err := repo.pool.Query(context.Background(), `SELECT id, name, description, price FROM products`)
	if err != nil {
		return products, err
	}
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
		if err != nil {
			return products, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (repo *PGRepo) DeleteProductByID(id int) error {
	_, err := repo.pool.Exec(context.Background(), `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PGRepo) GetAllProductsForSupplier(userID int) ([]models.Product, error) {
	var products []models.Product
	rows, err := repo.pool.Query(context.Background(), `SELECT id, name, description, price FROM products WHERE user_id=$1`, userID)
	if err != nil {
		return products, err
	}
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
		if err != nil {
			return products, err
		}
		products = append(products, product)
	}
	return products, nil
}
