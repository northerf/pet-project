package repository

import (
	"database/sql"
	"pet-project/pkg/model"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}

func (r *PostgresUserRepository) Create(user *model.User) error {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) Update(user *model.User) error {
	query := `UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4`
	_, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) Delete(user *model.User) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.DB.Exec(query, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (*model.User, error) {
	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	user := &model.User{}
	row := r.DB.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
