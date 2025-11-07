package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project_lab/internal/middleware"
	"project_lab/internal/models"
	"project_lab/internal/repositories"
)

type ProfileHandler struct {
	userRepo repositories.UserRepository
}

func NewProfileHandler(userRepo repositories.UserRepository) *ProfileHandler {
	return &ProfileHandler{userRepo: userRepo}
}

// GetProfileHandler lida com GET /profile
func (h *ProfileHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Obtém o ID do usuário do contexto JWT (AuthMiddleware)
	userIDValue := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Não autorizado. ID do usuário não encontrado.", http.StatusUnauthorized)
		return
	}

	profile, err := h.userRepo.GetUserProfile(userID)
	if err != nil {
		if err.Error() == "usuário não encontrado" {
			http.Error(w, "Perfil não encontrado.", http.StatusNotFound)
			return
		}
		fmt.Printf("Erro ao buscar perfil do BD: %v\n", err)
		http.Error(w, "Erro interno ao buscar perfil.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfileHandler lida com PATCH /profile
func (h *ProfileHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Obtém o ID do usuário do contexto JWT (AuthMiddleware)
	userIDValue := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Não autorizado. ID do usuário não encontrado.", http.StatusUnauthorized)
		return
	}

	var req models.UserProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida (JSON).", http.StatusBadRequest)
		return
	}

	// Validação de negócio (conforme o que foi feito no frontend)
	if req.Name == "" || len(req.Name) < 3 {
		http.Error(w, "Nome é obrigatório e deve ter no mínimo 3 caracteres.", http.StatusUnprocessableEntity)
		return
	}

	if err := h.userRepo.UpdateUserName(userID, req.Name); err != nil {
		fmt.Printf("Erro ao atualizar nome do usuário %d: %v\n", userID, err)
		http.Error(w, "Erro interno ao atualizar perfil.", http.StatusInternalServerError)
		return
	}

	// Retorna o perfil atualizado, conforme o YAML (200 OK)
	updatedProfile, err := h.userRepo.GetUserProfile(userID)
	if err != nil {
		fmt.Printf("Erro ao buscar perfil atualizado do BD: %v\n", err)
		// A atualização foi feita, mas falhamos ao ler.
		http.Error(w, "Perfil atualizado, mas falha ao retornar os dados.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProfile)
}
