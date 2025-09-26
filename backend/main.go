package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"project_lab/internal/handlers"
	"project_lab/internal/repositories"
	"project_lab/internal/services"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Printf("FRONTEND_URL não definido. Usando: %s", frontendURL)
	}

	db := connectDB()
	defer db.Close()
	createTables(db)

	// Inicializa as camadas da aplicação, injetando as dependências.
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", authHandler.RegisterUserHandler)
	mux.HandleFunc("/auth/login", authHandler.LoginUserHandler)

	// Configuração do middleware CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handlerWithCORS := c.Handler(mux)

	fmt.Println("🚀 Servidor rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", handlerWithCORS); err != nil {
		log.Fatal(err)
	}
}
