package repository

import (
	"User_Service/internal/models"
	"context"
)

func (repo *PGRepo) CreateUser(user models.User) (int, error) {
	err := repo.pool.QueryRow(context.Background(), `INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id`, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (repo *PGRepo) GetUserByEmail(email string) (user models.User, err error) {
	err = repo.pool.QueryRow(context.Background(), `SELECT id, username, email, password, role FROM users WHERE email=$1`, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	return user, err
}

func (repo *PGRepo) GetUserByID(id int) (user models.User, err error) {
	err = repo.pool.QueryRow(context.Background(), `SELECT id, username, email, password, role FROM users WHERE id=$1`, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	return user, err
}
