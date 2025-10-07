package repositories

import (
	"database/sql"
	"fmt"
	"project_lab/internal/models"
)

type TravelGroupRepository interface {
	ListGroupsByUserId(userID int) ([]models.TravelGroupListItem, error)
	CreateTravelGroup(group *models.TravelGroup) error
}

type postgresTravelGroupRepository struct {
	db *sql.DB
}

func NewTravelGroupRepository(db *sql.DB) TravelGroupRepository {
	return &postgresTravelGroupRepository{db: db}
}

// this method creates a travel group and sets its creator
func (r *postgresTravelGroupRepository) CreateTravelGroup(group *models.TravelGroup) error {

	// start of the transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação: %w", err)
	}
	// DEFER: rollback.
	defer tx.Rollback()

	query := `
        INSERT INTO travel_groups 
        (name, creator_id, start_date, end_date, created_at) 
        VALUES 
        ($1, $2, $3, $4, NOW())
        RETURNING id
    `
	err = tx.QueryRow(query,
		group.Name,
		group.CreatorID,
		group.StartDate,
		group.EndDate,
	).Scan(&group.ID)

	if err != nil {
		return fmt.Errorf("erro ao inserir grupo de viagem: %w", err)
	}

	memberQuery := `INSERT INTO group_members (travel_group_id, user_id, created_at) VALUES ($1, $2, NOW())`

	_, err = tx.Exec(memberQuery, group.ID, group.CreatorID)

	if err != nil {
		return fmt.Errorf("erro ao adicionar criador como membro: %w", err)
	}

	// end of transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("falha ao commitar transação: %w", err)
	}

	return nil
}

func (r *postgresTravelGroupRepository) ListGroupsByUserId(userId int) ([]models.TravelGroupListItem, error) {
	query := `
		SELECT 
			tg.id,
			tg.name,
			tg.start_date,
			tg.end_date,
			tg.creator_id,
			u.name AS creator_name,
			(SELECT COUNT(*) FROM group_members gm WHERE gm.travel_group_id = tg.id) AS member_count
		FROM 
			travel_groups tg
		JOIN 
			users u ON tg.creator_id = u.id
		WHERE
			-- O usuário é o criador OU é um membro
			tg.creator_id = $1 OR tg.id IN (SELECT travel_group_id FROM group_members WHERE user_id = $1)
		ORDER BY tg.start_date DESC;
	`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	groups := []models.TravelGroupListItem{}
	for rows.Next() {
		var g models.TravelGroupListItem
		var memberCount sql.NullInt32 // Usar sql.NullInt32 para garantir compatibilidade com COUNT

		err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.StartDate,
			&g.EndDate,
			&g.CreatorId,
			&g.CreatorName,
			&memberCount,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler linha: %w", err)
		}

		// Garante que a contagem é um inteiro
		if memberCount.Valid {
			g.MemberCount = int(memberCount.Int32)
		} else {
			g.MemberCount = 0
		}

		groups = append(groups, g)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iteração: %w", err)
	}

	return groups, nil
}
