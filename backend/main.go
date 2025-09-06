package main

import (
	"fmt"
	"log"
	"net/http"
	"project_lab/internal/handlers"
	"project_lab/internal/repositories"
	"project_lab/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	// Carrega as variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	//Conecta ao banco de dados e garante que as tabelas existam.
	db := connectDB()
	defer db.Close()
	createTables(db)

	//Inicializa as camadas da aplicação, injetando as dependências.
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	//Registra as rotas da API.
	http.HandleFunc("/auth/register", authHandler.RegisterUserHandler)
	http.HandleFunc("/auth/login", authHandler.LoginUserHandler)

	fmt.Println("🚀 Servidor rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
