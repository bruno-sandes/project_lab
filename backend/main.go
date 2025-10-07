package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"project_lab/internal/handlers"
	"project_lab/internal/middleware"
	"project_lab/internal/repositories"
	"project_lab/internal/services"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func groupsRouter(h *handlers.TravelGroupHandler) http.HandlerFunc {
	// Retorna uma função anônima que é o nosso Handler unificado
	return func(w http.ResponseWriter, r *http.Request) {
		// Usa o r.Method para decidir qual Handler real chamar
		switch r.Method {
		case "GET":
			h.ListGroups(w, r) // Vai listar os grupos
		case "POST":
			h.CreateGroupHandler(w, r) // Vai criar um novo grupo
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}
}

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

	travelGroupsRepo := repositories.NewTravelGroupRepository(db)
	travelGroupsHandler := handlers.NewTravelGroupHandler(travelGroupsRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", authHandler.RegisterUserHandler)
	mux.HandleFunc("/auth/login", authHandler.LoginUserHandler)
	mux.Handle("/groups", middleware.AuthMiddleware(groupsRouter(travelGroupsHandler)))

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
