// internal/usecase/auth.go
package usecase

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v4"
    "github.com/yourusername/auth-service/internal/entity"
)

type AuthRepository interface {
    CreateUser(user *entity.User) error
    GetUserByEmail(email string) (*entity.User, error)
}

type AuthUseCase struct {
    repo   AuthRepository
    jwtKey string
}

func NewAuthUseCase(repo AuthRepository, jwtKey string) *AuthUseCase {
    return &AuthUseCase{repo: repo, jwtKey: jwtKey}
}

func (u *AuthUseCase) Signup(email, password string) error {
    user := &entity.User{Email: email, Password: password}
    if email == "" || password == "" {
        return errors.New("email and password are required")
    }
    if err := user.HashPassword(); err != nil {
        return err
    }
    return u.repo.CreateUser(user)
}

func (u *AuthUseCase) Login(email, password string) (string, error) {
    user, err := u.repo.GetUserByEmail(email)
    if err != nil || user == nil || !user.CheckPassword(password) {
        return "", errors.New("invalid credentials")
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })
    return token.SignedString([]byte(u.jwtKey))
}