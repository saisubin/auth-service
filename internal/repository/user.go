// internal/repository/user.go
package repository

import (
    "database/sql"
    "github.com/saisubin/auth-service/internal/entity"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *entity.User) error {
    query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
    return r.db.QueryRow(query, user.Email, user.Password).Scan(&user.ID)
}

func (r *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
    user := &entity.User{}
    query := "SELECT id, email, password FROM users WHERE email = $1"
    err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return user, err
}