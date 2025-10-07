package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project_lab/internal/middleware"
	"project_lab/internal/models"
	"project_lab/internal/repositories"
)

type TravelGroupHandler struct {
	repo repositories.TravelGroupRepository
}

func NewTravelGroupHandler(repo repositories.TravelGroupRepository) *TravelGroupHandler {
	return &TravelGroupHandler{repo: repo}
}

func (h *TravelGroupHandler) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {

	userIDValue := r.Context().Value(middleware.UserIDKey)
	creatorID, ok := userIDValue.(int)
	if !ok {
		// Se o middleware de Auth falhar, ele já deveria ter retornado 401.
		http.Error(w, "Não autorizado. ID do usuário não encontrado.", http.StatusUnauthorized)
		return
	}

	var group models.TravelGroup
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Requisição inválida ou formato JSON incorreto.", http.StatusBadRequest)
		return
	}

	// VALIDAÇÃO SIMPLES (Retorna 422 Unprocessable Entity)
	if group.Name == "" || group.StartDate.IsZero() || group.EndDate.IsZero() {
		http.Error(w, "Nome, data de início e data de término são obrigatórios.", http.StatusUnprocessableEntity)
		return
	}
	//  Data de Início deve ser antes da Data de Fim
	if group.StartDate.After(group.EndDate) {
		http.Error(w, "A data de início deve ser anterior ou igual à data de término.", http.StatusUnprocessableEntity)
		return
	}

	// O ID é pego do contexto, garantindo que o usuário só pode criar grupos para si mesmo.
	group.CreatorID = creatorID

	if err := h.repo.CreateTravelGroup(&group); err != nil {
		fmt.Printf("Erro ao criar grupo no BD: %v\n", err)
		http.Error(w, "Erro interno ao salvar grupo de viagem.", http.StatusInternalServerError)
		return
	}

	// Define cabeçalhos e status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Retorna a struct completa do grupo, incluindo o novo ID gerado pelo banco.
	if err := json.NewEncoder(w).Encode(group); err != nil {
		http.Error(w, "Erro ao serializar resposta.", http.StatusInternalServerError)
		return
	}
}

func (h *TravelGroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {

	userIdValue := r.Context().Value(middleware.UserIDKey)

	userID, ok := userIdValue.(int)
	if !ok {
		http.Error(w, "Falha na autenticação. ID de usuário não disponível.", http.StatusInternalServerError)
		return
	}

	groups, err := h.repo.ListGroupsByUserId(userID)
	if err != nil {
		// Loga o erro detalhado no backend
		fmt.Printf("Erro ao buscar grupos para userID %d: %v\n", userID, err)

		// Retorna um erro genérico para o frontend
		http.Error(w, "Erro interno ao buscar grupos de viagem.", http.StatusInternalServerError)
		return
	}

	// Define o cabeçalho para JSON.
	w.Header().Set("Content-Type", "application/json")

	// Resposta 200 OK (implícito pelo Encoder ou explícito com w.WriteHeader(http.StatusOK))

	// Codifica o resultado (slice de TravelGroupListItem) para JSON e envia.
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, "Erro ao serializar resposta JSON.", http.StatusInternalServerError)
		return
	}
}
