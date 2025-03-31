package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authUseCase := usecase.NewAuthUseCase(userRepo)

	handler.SetAuthUseCase(authUseCase) // Set the use case for handlers

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/signup", handler.SignupHandler).Methods("POST")
	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
