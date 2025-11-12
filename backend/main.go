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
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func profileRouter(h *handlers.ProfileHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/profile" {
			switch r.Method {
			case "GET":
				h.GetProfileHandler(w, r)
			case "PATCH":
				h.UpdateProfileHandler(w, r)
			default:
				http.Error(w, "Método não permitido para /profile.", http.StatusMethodNotAllowed)
			}
			return
		}

		http.NotFound(w, r)
	}
}

func groupsRouter(h *handlers.TravelGroupHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		if len(pathSegments) == 1 && pathSegments[0] == "groups" {
			switch r.Method {
			case "GET":
				h.ListGroups(w, r)
			case "POST":
				h.CreateGroupHandler(w, r)
			default:
				http.Error(w, "Método não permitido para /groups", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(pathSegments) >= 2 && pathSegments[0] == "groups" {
			groupIDStr := pathSegments[1]

			if len(pathSegments) == 3 {
				resource := pathSegments[2]

				switch resource {
				case "members":
					if r.Method == "GET" {
						h.ListGroupMembersHandler(w, r, groupIDStr)
						return
					}
				case "destinations":
					switch r.Method {
					case "GET":
						h.ListGroupDestinationsHandler(w, r, groupIDStr)
					case "POST":
						h.CreateDestinationHandler(w, r, groupIDStr)
					default:
						http.Error(w, "Método não permitido para /destinations", http.StatusMethodNotAllowed)
					}
					return
				case "votings":
					switch r.Method {
					case "GET":
						h.ListGroupVotingsHandler(w, r, groupIDStr)
					case "POST":
						h.CreateVotingHandler(w, r, groupIDStr)
					default:
						http.Error(w, "Método não permitido para /votings", http.StatusMethodNotAllowed)
					}
					return
				case "expenses":
					if r.Method == "GET" {
						h.ListGroupExpensesHandler(w, r, groupIDStr)
						return
					}
					// A lógica de POST para expenses (futuro) entraria aqui
				}
				http.Error(w, "Recurso ou Método não permitido.", http.StatusMethodNotAllowed)
				return
			}

			if len(pathSegments) == 2 {
				if r.Method == "GET" {
					h.GetGroupDetailsWithID(w, r, groupIDStr)
					return
				}
				http.Error(w, "Método não permitido para detalhes do grupo", http.StatusMethodNotAllowed)
				return
			}
		}

		http.NotFound(w, r)
	}
}

func votingsRouter(h *handlers.VoteHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		// Esperamos a rota /votings/{id}/vote
		if len(pathSegments) == 3 && pathSegments[0] == "votings" && pathSegments[2] == "vote" {
			votingIDStr := pathSegments[1]

			if r.Method == "POST" {
				h.VoteHandler(w, r, votingIDStr)
				return
			}
		}

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
	profileHandler := handlers.NewProfileHandler(userRepo)

	travelGroupsRepo := repositories.NewTravelGroupRepository(db)
	travelGroupsHandler := handlers.NewTravelGroupHandler(travelGroupsRepo)

	voteRepo := repositories.NewVoteRepository(db)
	voteHandler := handlers.NewVoteHandler(voteRepo, travelGroupsRepo) // Passa travelGroupsRepo para validações

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", authHandler.RegisterUserHandler)
	mux.Handle("/auth/login", middleware.RateLimitMiddleware(http.HandlerFunc(authHandler.LoginUserHandler)))
	mux.Handle("/profile", middleware.AuthMiddleware(profileRouter(profileHandler)))
	mux.Handle("/groups/", middleware.AuthMiddleware(groupsRouter(travelGroupsHandler)))
	mux.Handle("/groups", middleware.AuthMiddleware(groupsRouter(travelGroupsHandler)))
	mux.Handle("/votings/", middleware.AuthMiddleware(votingsRouter(voteHandler)))

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
