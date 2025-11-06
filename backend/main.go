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
	return func(w http.ResponseWriter, r *http.Request) {

		// Rota de Listagem/Criação (Exatamente /groups)
		if r.URL.Path == "/groups" || r.URL.Path == "/groups/" {
			switch r.Method {
			case "GET":
				h.ListGroups(w, r) // GET /groups -> Listar
			case "POST":
				h.CreateGroupHandler(w, r) // POST /groups -> Criar
			default:
				http.Error(w, "Método não permitido para /groups", http.StatusMethodNotAllowed)
			}
			return
		}

		// Rota de Detalhe (Exemplo: /groups/123)
		// Verifica se o caminho começa com "/groups/" e tem mais caracteres (o ID)
		if len(r.URL.Path) > len("/groups/") && r.URL.Path[:len("/groups/")] == "/groups/" {

			// Tenta extrair o ID do final da rota
			groupIDStr := r.URL.Path[len("/groups/"):]

			// O ID deve ser um número inteiro, e não deve conter sub-caminhos (ex: /groups/123/members)
			switch r.Method {
			case "GET":
				h.GetGroupDetailsWithID(w, r, groupIDStr)
			default:
				http.Error(w, "Método não permitido para detalhes do grupo", http.StatusMethodNotAllowed)
			}
			return
		}

		// Se a URL não for /groups e nem /groups/{id}, retorna 404
		http.NotFound(w, r)
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
