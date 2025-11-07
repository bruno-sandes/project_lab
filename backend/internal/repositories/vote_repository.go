package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"project_lab/internal/models"
)

type VoteRepository interface {
	CastVote(vote *models.Vote) error
	CheckUserVote(votingID int, userID int) (bool, error)
	GetVotingOptions(votingID int) ([]string, error)
}

type postgresVoteRepository struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) VoteRepository {
	return &postgresVoteRepository{db: db}
}

// CheckUserVote verifica se o usuário já votou nesta votação.
func (r *postgresVoteRepository) CheckUserVote(votingID int, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM votes WHERE voting_id = $1 AND user_id = $2;`
	var count int
	err := r.db.QueryRow(query, votingID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao checar voto existente: %w", err)
	}
	return count > 0, nil
}

// GetVotingOptions busca as opções válidas para uma votação (requer o parse JSON)
func (r *postgresVoteRepository) GetVotingOptions(votingID int) ([]string, error) {
	query := `SELECT options FROM votings WHERE id = $1;`
	var optionsJSON string
	err := r.db.QueryRow(query, votingID).Scan(&optionsJSON)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("votação não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar opções da votação: %w", err)
	}

	var options []string
	if err := json.Unmarshal([]byte(optionsJSON), &options); err != nil {
		return nil, fmt.Errorf("erro ao deserializar opções JSON: %w", err)
	}

	return options, nil
}

// CastVote insere ou atualiza o voto do usuário (dependendo da sua regra de negócio, aqui faremos INSERIR)
func (r *postgresVoteRepository) CastVote(vote *models.Vote) error {
	query := `
		INSERT INTO votes 
		(voting_id, user_id, selected_option, created_at) 
		VALUES 
		($1, $2, $3, NOW());
	`
	_, err := r.db.Exec(query, vote.VotingID, vote.UserID, vote.SelectedOption)
	if err != nil {
		return fmt.Errorf("erro ao registrar voto: %w", err)
	}
	return nil
}
