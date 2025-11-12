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

// checkGroupMembership é uma função auxiliar interna para verificar a autorização (Mitigação A01).
// Ela reutiliza o GetGroupDetails para garantir que o usuário é membro.
func (h *TravelGroupHandler) checkGroupMembership(w http.ResponseWriter, r *http.Request, groupID int) (userID int, ok bool) {
	userIDValue := r.Context().Value(middleware.UserIDKey)
	userID, ok = userIDValue.(int)
	if !ok {
		http.Error(w, "Não autorizado. ID do usuário não encontrado.", http.StatusUnauthorized)
		return 0, false
	}

	// MITIGAÇÃO A01: Verifica se o usuário tem permissão para acessar este groupID
	_, err := h.repo.GetGroupDetails(groupID, userID)
	if err != nil {
		// Se GetGroupDetails falhar, o usuário não é membro ou o grupo não existe.
		http.Error(w, "Grupo não encontrado ou não autorizado", http.StatusNotFound)
		return userID, false
	}

	return userID, true
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
		fmt.Printf("Erro ao buscar grupos para userID %d: %v\n", userID, err)
		http.Error(w, "Erro interno ao buscar grupos de viagem.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, "Erro ao serializar resposta JSON.", http.StatusInternalServerError)
		return
	}
}

// GetGroupDetailsWithID (MITIGADO)
func (h *TravelGroupHandler) GetGroupDetailsWithID(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido. Deve ser um número.", http.StatusBadRequest)
		return
	}

	// MITIGAÇÃO A01: Verifica se o usuário é membro ANTES de buscar detalhes.
	userID, ok := h.checkGroupMembership(w, r, groupID)
	if !ok {
		return // O erro já foi enviado pela função auxiliar
	}

	details, err := h.repo.GetGroupDetails(groupID, userID)
	if err != nil {
		// Este erro não deve ocorrer se o checkGroupMembership passou, mas é uma boa defesa.
		fmt.Printf("Erro ao buscar detalhes do grupo %d: %v\n", groupID, err)
		http.Error(w, "Erro interno ao buscar detalhes do grupo.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

// ListGroupMembersHandler (MITIGADO)
func (h *TravelGroupHandler) ListGroupMembersHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// MITIGAÇÃO A01 (IDOR): Verifica se o usuário é membro ANTES de listar.
	if _, ok := h.checkGroupMembership(w, r, groupID); !ok {
		return // Bloqueia se não for membro
	}

	members, err := h.repo.ListGroupMembers(groupID)
	if err != nil {
		fmt.Printf("Erro ao buscar lista de membros do grupo %d: %v\n", groupID, err)
		http.Error(w, "Erro interno ao buscar membros do grupo.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// ListGroupDestinationsHandler (MITIGADO)
func (h *TravelGroupHandler) ListGroupDestinationsHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// MITIGAÇÃO A01 (IDOR): Verifica se o usuário é membro ANTES de listar.
	if _, ok := h.checkGroupMembership(w, r, groupID); !ok {
		return // Bloqueia se não for membro
	}

	destinations, err := h.repo.ListGroupDestinations(groupID)
	if err != nil {
		fmt.Printf("Erro ao buscar lista de destinos do grupo %d: %v\n", groupID, err)
		http.Error(w, "Erro interno ao buscar destinos do grupo.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(destinations)
}

// ListGroupVotingsHandler (MITIGADO)
func (h *TravelGroupHandler) ListGroupVotingsHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// MITIGAÇÃO A01 (IDOR): Verifica se o usuário é membro ANTES de listar.
	userID, ok := h.checkGroupMembership(w, r, groupID)
	if !ok {
		return // Bloqueia se não for membro
	}

	votings, err := h.repo.ListGroupVotings(groupID, userID)
	if err != nil {
		fmt.Printf("Erro ao buscar votações do grupo %d: %v\n", groupID, err)
		http.Error(w, "Erro interno ao buscar votações do grupo.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votings)
}

// ListGroupExpensesHandler (MITIGADO)
func (h *TravelGroupHandler) ListGroupExpensesHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// MITIGAÇÃO A01 (IDOR): Verifica se o usuário é membro ANTES de listar.
	if _, ok := h.checkGroupMembership(w, r, groupID); !ok {
		return // Bloqueia se não for membro
	}

	expenses, err := h.repo.ListGroupExpenses(groupID)
	if err != nil {
		fmt.Printf("Erro ao buscar despesas do grupo %d: %v\n", groupID, err)
		http.Error(w, "Erro interno ao buscar despesas do grupo.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// CreateDestinationHandler (MITIGADO)
func (h *TravelGroupHandler) CreateDestinationHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// MITIGAÇÃO A01 (Criação Arbitrária): Verifica se o usuário é membro ANTES de criar.
	if _, ok := h.checkGroupMembership(w, r, groupID); !ok {
		return // Bloqueia se não for membro
	}

	var req models.DestinationCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida (JSON).", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "O nome do destino é obrigatório.", http.StatusUnprocessableEntity)
		return
	}

	destination := models.Destination{
		TravelGroupID: groupID,
		Name:          req.Name,
		Location:      req.Location,
		Description:   req.Description,
	}

	if err := h.repo.CreateDestination(&destination); err != nil {
		fmt.Printf("Erro ao criar destino no BD: %v\n", err)
		http.Error(w, "Erro interno ao salvar destino.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(destination)
}

// CreateVotingHandler (MITIGADO)
func (h *TravelGroupHandler) CreateVotingHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// MITIGAÇÃO A01 (Criação Arbitrária): Verifica se o usuário é membro ANTES de criar.
	if _, ok := h.checkGroupMembership(w, r, groupID); !ok {
		return // Bloqueia se não for membro
	}

	var req models.VotingCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida (JSON).", http.StatusBadRequest)
		return
	}

	if req.Question == "" || len(req.Options) < 2 {
		http.Error(w, "A pergunta e pelo menos 2 opções são obrigatórias.", http.StatusUnprocessableEntity)
		return
	}

	optionsJSON, err := json.Marshal(req.Options)
	if err != nil {
		http.Error(w, "Erro ao processar opções da votação.", http.StatusInternalServerError)
		return
	}

	newVotingID, err := h.repo.CreateVoting(groupID, req.Question, string(optionsJSON))
	if err != nil {
		fmt.Printf("Erro ao criar votação no BD: %v\n", err)
		http.Error(w, "Erro interno ao salvar votação.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": newVotingID})
}

func (h *TravelGroupHandler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// 1. MITIGAÇÃO A01 (Check de Membro): Verifica se o usuário é membro do grupo.
	// O userID é o ID do usuário AUTENTICADO (do token).
	userID, ok := h.checkGroupMembership(w, r, groupID)
	if !ok {
		return // Bloqueia se não for membro (IDOR)
	}

	var req models.ExpenseCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida (JSON).", http.StatusBadRequest)
		return
	}

	// Validações básicas
	if req.Description == "" || req.Amount <= 0 || len(req.ParticipantIDs) == 0 {
		// Removemos a checagem de req.PayerID <= 0, pois não usaremos o PayerID do JSON.
		http.Error(w, "Descrição, valor (positivo) e lista de participantes são obrigatórios.", http.StatusUnprocessableEntity)
		return
	}

	// 2. MITIGAÇÃO A01 (AÇÃO FORJADA):
	// O PayerID da despesa AGORA usa o ID do usuário AUTENTICADO (userID)
	// em vez de confiar no valor enviado no corpo da requisição (req.PayerID).
	expense := models.Expense{
		TravelGroupID:  groupID,
		Description:    req.Description,
		Amount:         req.Amount,
		PayerID:        userID, // <-- CORRIGIDO! Usa o ID do token.
		ParticipantIDs: req.ParticipantIDs,
	}

	if err := h.repo.CreateExpense(&expense); err != nil {
		fmt.Printf("Erro ao criar despesa no BD: %v\n", err)
		http.Error(w, "Erro interno ao salvar despesa e participantes.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expense)
}
