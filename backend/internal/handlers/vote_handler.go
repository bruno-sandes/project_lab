package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project_lab/internal/middleware"
	"project_lab/internal/models"
	"project_lab/internal/repositories"
	"strconv"
)

type VoteHandler struct {
	voteRepo  repositories.VoteRepository
	groupRepo repositories.TravelGroupRepository
}

func NewVoteHandler(voteRepo repositories.VoteRepository, groupRepo repositories.TravelGroupRepository) *VoteHandler {
	return &VoteHandler{voteRepo: voteRepo, groupRepo: groupRepo}
}

// VoteHandler lida com o registro de um voto
func (h *VoteHandler) VoteHandler(w http.ResponseWriter, r *http.Request, votingIDStr string) {

	votingID, err := strconv.Atoi(votingIDStr)
	if err != nil {
		http.Error(w, "ID da votação inválido.", http.StatusBadRequest)
		return
	}

	userIDValue := r.Context().Value(middleware.UserIDKey)
	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Não autorizado. ID do usuário não encontrado.", http.StatusUnauthorized)
		return
	}

	var req models.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida (JSON).", http.StatusBadRequest)
		return
	}

	if req.SelectedOption == "" {
		http.Error(w, "Opção de voto é obrigatória.", http.StatusUnprocessableEntity)
		return
	}

	//  Lógica de Validação e Registro
	validOptions, err := h.voteRepo.GetVotingOptions(votingID)
	if err != nil {
		if err.Error() == "votação não encontrada" {
			http.Error(w, "Votação não encontrada.", http.StatusNotFound)
			return
		}
		fmt.Printf("Erro ao buscar opções da votação %d: %v\n", votingID, err)
		http.Error(w, "Erro interno de validação.", http.StatusInternalServerError)
		return
	}

	isValidOption := false
	for _, opt := range validOptions {
		if opt == req.SelectedOption {
			isValidOption = true
			break
		}
	}
	if !isValidOption {
		http.Error(w, "Opção de voto inválida para esta votação.", http.StatusUnprocessableEntity)
		return
	}

	//  Verificar se o usuário já votou
	alreadyVoted, err := h.voteRepo.CheckUserVote(votingID, userID)
	if err != nil {
		fmt.Printf("Erro ao checar voto: %v\n", err)
		http.Error(w, "Erro interno de checagem de voto.", http.StatusInternalServerError)
		return
	}

	if alreadyVoted {
		// Se a regra é "só pode votar uma vez", retorne um erro.
		http.Error(w, "Você já votou nesta enquete.", http.StatusConflict)
		return
	}

	// Nota: Idealmente, você também deveria checar se o usuário (userID) é membro do grupo
	// ao qual a votação pertence. Isso requer uma função extra no TravelGroupRepository.

	vote := models.Vote{
		VotingID:       votingID,
		UserID:         userID,
		SelectedOption: req.SelectedOption,
	}

	if err := h.voteRepo.CastVote(&vote); err != nil {
		fmt.Printf("Erro ao registrar voto: %v\n", err)
		http.Error(w, "Erro interno ao registrar voto.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Voto registrado com sucesso."}`))
}
