// cmd/main.go
package main

import (
    "database/sql"
    "log"
    "os"
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
    "github.com/saisubin/auth-service/internal/handler"
    "github.com/saisubin/auth-service/internal/repository"
    "github.com/saisubin/auth-service/internal/usecase"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    if err := db.Ping(); err != nil {
        log.Fatal("Database ping failed:", err)
    }
    repo := repository.NewUserRepository(db)
    uc := usecase.NewAuthUseCase(repo, os.Getenv("JWT_SECRET"))
    h := handler.NewAuthHandler(uc, os.Getenv("JWT_SECRET"))
    r := gin.Default()
    r.POST("/signup", h.Signup)
    r.POST("/login", h.Login)
    protected := r.Group("/api")
    protected.Use(h.AuthMiddleware)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server starting on :%s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}