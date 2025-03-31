// internal/handler/auth.go
package handler

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "github.com/yourusername/auth-service/internal/usecase"
)

type AuthHandler struct {
    uc     *usecase.AuthUseCase
    jwtKey string
}

func NewAuthHandler(uc *usecase.AuthUseCase, jwtKey string) *AuthHandler {
    return &AuthHandler{uc: uc, jwtKey: jwtKey}
}

func (h *AuthHandler) Signup(c *gin.Context) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.BindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }
    if err := h.uc.Signup(input.Email, input.Password); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (h *AuthHandler) Login(c *gin.Context) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.BindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }
    token, err := h.uc.Login(input.Email, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) AuthMiddleware(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
        c.Abort()
        return
    }
    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(h.jwtKey), nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
        c.Abort()
        return
    }
    c.Next()
}