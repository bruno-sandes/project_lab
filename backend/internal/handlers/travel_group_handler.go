package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project_lab/internal/middleware"
	"project_lab/internal/models"
	"project_lab/internal/repositories"
	"strconv"
	"time"
)

type TravelGroupHandler struct {
	repo repositories.TravelGroupRepository
}

func NewTravelGroupHandler(repo repositories.TravelGroupRepository) *TravelGroupHandler {
	return &TravelGroupHandler{repo: repo}
}

func (h *TravelGroupHandler) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {

	var req models.TravelGroupCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida ou formato JSON incorreto.", http.StatusBadRequest)
		return
	}

	const layout = "2006-01-02"

	startDate, err := time.Parse(layout, req.StartDate)
	if err != nil {
		http.Error(w, "Formato de data de início inválido. Use YYYY-MM-DD.", http.StatusUnprocessableEntity)
		return
	}

	endDate, err := time.Parse(layout, req.EndDate)
	if err != nil {
		http.Error(w, "Formato de data de término inválido. Use YYYY-MM-DD.", http.StatusUnprocessableEntity)
		return
	}

	if req.Name == "" || req.StartDate == "" || req.EndDate == "" {
		http.Error(w, "Nome, data de início e data de término são obrigatórios.", http.StatusUnprocessableEntity)
		return
	}

	if startDate.After(endDate) {
		http.Error(w, "A data de início deve ser anterior ou igual à data de término.", http.StatusUnprocessableEntity)
		return
	}

	userIDValue := r.Context().Value(middleware.UserIDKey)
	creatorID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Não autorizado. ID do usuário não encontrado.", http.StatusUnauthorized)
		return
	}

	group := models.TravelGroup{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatorID:   creatorID,
	}

	if err := h.repo.CreateTravelGroup(&group); err != nil {
		fmt.Printf("Erro ao criar grupo no BD: %v\n", err)
		http.Error(w, "Erro interno ao salvar grupo de viagem.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

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
func (h *TravelGroupHandler) GetGroupDetailsWithID(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido. Deve ser um número.", http.StatusBadRequest)
		return
	}

	// Obter o ID do usuário autenticado (lógica de autorização)
	userIDValue := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Não autorizado. ID do usuário não encontrado.", http.StatusUnauthorized)
		return
	}

	details, err := h.repo.GetGroupDetails(groupID, userID)
	if err != nil {
		if err.Error() == "grupo não encontrado ou usuário não autorizado a visualizá-lo" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		fmt.Printf("Erro ao buscar detalhes do grupo %d: %v\n", groupID, err)
		http.Error(w, "Erro interno ao buscar detalhes do grupo.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}
