package handlers

import (
	"encoding/json"
	"net/http"
	"project_lab/internal/models"
	"project_lab/internal/services"
)

// AuthHandler gerencia as requisições HTTP para a autenticação.
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler cria uma nova instância de AuthHandler.
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterUserHandler lida com a requisição de cadastro de usuário.
func (h *AuthHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Dados de requisição inválidos", http.StatusBadRequest)
		return
	}

	if err := h.authService.RegisterUser(&user); err != nil {
		if err.Error() == "e-mail já está em uso" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Erro ao registrar usuário", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Usuário registrado com sucesso!"))
}

// LoginUserHandler lida com a requisição de login de usuário.
func (h *AuthHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Dados de requisição inválidos", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Authenticate(loginRequest.Email, loginRequest.Password)
	if err != nil {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	response := map[string]string{"token": token}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Erro ao serializar resposta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
