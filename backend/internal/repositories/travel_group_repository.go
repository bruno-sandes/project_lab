package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"project_lab/internal/models"
	"strconv"
	"strings"
)

type TravelGroupRepository interface {
	ListGroupsByUserId(userID int) ([]models.TravelGroupListItem, error)
	CreateTravelGroup(group *models.TravelGroup) error
	GetGroupDetails(groupID int, userID int) (*models.TravelGroupDetails, error)
	ListGroupMembers(groupID int) ([]models.GroupMemberDTO, error)
	ListGroupDestinations(groupID int) ([]models.DestinationDTO, error)
	ListGroupVotings(groupID int, userID int) ([]models.VotingDTO, error)
	ListGroupExpenses(groupID int) ([]models.ExpenseDTO, error)
	CreateDestination(destination *models.Destination) error
	CreateVoting(groupID int, question string, optionsJSON string) (int, error)
	CreateExpense(expense *models.Expense) error
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
        (name,description , creator_id, start_date, end_date, created_at) 
        VALUES 
        ($1, $2, $3, $4, $5, NOW())
        RETURNING id
    `
	err = tx.QueryRow(query,
		group.Name,
		group.Description,
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
			tg.description,
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
			&g.Description,
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

func (r *postgresTravelGroupRepository) GetGroupDetails(groupID int, userID int) (*models.TravelGroupDetails, error) {

	// Consulta SQL para obter detalhes básicos, nome do criador e contagem de membros.
	// Também valida se o usuário (userID) é membro/criador.
	query := `
        SELECT 
            tg.id,
            tg.name,
            tg.description,
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
            tg.id = $1
            -- Verifica se o usuário é membro ou criador antes de retornar os detalhes
            AND (tg.creator_id = $2 OR tg.id IN (SELECT travel_group_id FROM group_members WHERE user_id = $2));
    `

	var details models.TravelGroupDetails
	var memberCount sql.NullInt32

	err := r.db.QueryRow(query, groupID, userID).Scan(
		&details.ID,
		&details.Name,
		&details.Description,
		&details.StartDate,
		&details.EndDate,
		&details.CreatorID,
		&details.CreatorName,
		&memberCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("grupo não encontrado ou usuário não autorizado a visualizá-lo")
		}
		return nil, fmt.Errorf("erro ao buscar detalhes do grupo: %w", err)
	}

	// Garante que a contagem é um inteiro
	if memberCount.Valid {
		details.MemberCount = int(memberCount.Int32)
	} else {
		details.MemberCount = 0
	}

	return &details, nil
}

func (r *postgresTravelGroupRepository) ListGroupMembers(groupID int) ([]models.GroupMemberDTO, error) {
	query := `
        SELECT 
            u.id AS user_id,
            u.name,
            u.email,
            CASE 
                WHEN tg.creator_id = u.id THEN 'Organizador'
                ELSE 'Participante'
            END AS role
        FROM 
            group_members gm
        JOIN 
            users u ON gm.user_id = u.id
        JOIN 
            travel_groups tg ON gm.travel_group_id = tg.id
        WHERE 
            gm.travel_group_id = $1
        ORDER BY role DESC, u.name ASC; -- Organizador sempre primeiro
    `

	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar membros do grupo: %w", err)
	}
	defer rows.Close()

	members := []models.GroupMemberDTO{}
	for rows.Next() {
		var member models.GroupMemberDTO
		err := rows.Scan(
			&member.UserID,
			&member.Name,
			&member.Email,
			&member.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear membro: %w", err)
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro na iteração dos membros: %w", err)
	}

	return members, nil
}

func (r *postgresTravelGroupRepository) ListGroupDestinations(groupID int) ([]models.DestinationDTO, error) {
	query := `
        SELECT 
            id,
            name,
            location,
            description
        FROM 
            destinations
        WHERE 
            travel_group_id = $1
        ORDER BY name ASC;
    `

	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar destinos do grupo: %w", err)
	}
	defer rows.Close()

	destinations := []models.DestinationDTO{}
	for rows.Next() {
		var d models.DestinationDTO
		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Location,
			&d.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear destino: %w", err)
		}
		destinations = append(destinations, d)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro na iteração dos destinos: %w", err)
	}

	return destinations, nil
}

func (r *postgresTravelGroupRepository) ListGroupVotings(groupID int, userID int) ([]models.VotingDTO, error) {
	// Usamos left join para trazer a votação do usuário (se existir)
	query := `
        SELECT 
            v.id,
            v.question,
            v.options,
            COUNT(vt.id) AS total_votes,
            uv.selected_option AS user_vote_option,
            v.created_at
        FROM 
            votings v
        LEFT JOIN 
            votes vt ON v.id = vt.voting_id
        LEFT JOIN 
            votes uv ON v.id = uv.voting_id AND uv.user_id = $2 -- Voto do Usuário Logado
        WHERE 
            v.travel_group_id = $1
        GROUP BY v.id, v.question, v.options, v.created_at, uv.selected_option
        ORDER BY v.created_at DESC;
    `

	rows, err := r.db.Query(query, groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar votações: %w", err)
	}
	defer rows.Close()

	votings := []models.VotingDTO{}
	for rows.Next() {
		var v models.VotingDTO
		var optionsJSON string
		var totalVotes sql.NullInt64
		var userVote sql.NullString

		err := rows.Scan(
			&v.ID,
			&v.Question,
			&optionsJSON, // String JSON
			&totalVotes,
			&userVote,
			&v.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear votação: %w", err)
		}

		// 1. Deserializar as opções
		if err := json.Unmarshal([]byte(optionsJSON), &v.Options); err != nil {
			// Logar ou tratar erro, mas continuar se possível
			fmt.Printf("Aviso: Falha ao deserializar opções JSON para votação %d: %v\n", v.ID, err)
		}

		// 2. Setar total de votos
		v.TotalVotes = int(totalVotes.Int64)

		// 3. Setar voto do usuário (se houver)
		if userVote.Valid {
			v.UserVote = &userVote.String
		}

		votings = append(votings, v)
	}

	return votings, nil
}

func (r *postgresTravelGroupRepository) ListGroupExpenses(groupID int) ([]models.ExpenseDTO, error) {
	// A consulta usa a agregação STRING_AGG para obter a lista de IDs de participantes
	query := `
        SELECT 
            e.id,
            e.description,
            e.amount,
            e.payer_id,
            u.name AS payer_name,
            e.created_at,
            COUNT(ep.user_id) AS participants_count,
            COALESCE(STRING_AGG(ep.user_id::text, ',' ORDER BY ep.user_id), '') AS participants_ids
        FROM 
            expenses e
        JOIN 
            users u ON e.payer_id = u.id
        LEFT JOIN 
            expense_participants ep ON e.id = ep.expense_id
        WHERE 
            e.travel_group_id = $1
        GROUP BY e.id, u.name
        ORDER BY e.created_at DESC;
    `

	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar despesas: %w", err)
	}
	defer rows.Close()

	expenses := []models.ExpenseDTO{}
	for rows.Next() {
		var e models.ExpenseDTO
		var participantsIDsStr sql.NullString
		var participantsCount sql.NullInt64

		err := rows.Scan(
			&e.ID,
			&e.Description,
			&e.Amount, // Mapeado diretamente para float64 (no Struct)
			&e.PayerID,
			&e.PayerName,
			&e.CreatedAt,
			&participantsCount,
			&participantsIDsStr, // IDs separados por vírgula
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear despesa: %w", err)
		}

		// 1. Setar contagem de participantes
		e.ParticipantsCount = int(participantsCount.Int64)

		// 2. Processar a lista de IDs de participantes
		e.ParticipantsIDs = []int{}
		if participantsIDsStr.Valid && participantsIDsStr.String != "" {
			parts := strings.Split(participantsIDsStr.String, ",")
			for _, p := range parts {
				if id, err := strconv.Atoi(p); err == nil {
					e.ParticipantsIDs = append(e.ParticipantsIDs, id)
				}
			}
		}

		expenses = append(expenses, e)
	}

	return expenses, nil
}
func (r *postgresTravelGroupRepository) CreateDestination(destination *models.Destination) error {
	query := `
        INSERT INTO destinations 
        (travel_group_id, name, location, description, created_at) 
        VALUES 
        ($1, $2, $3, $4, NOW())
        RETURNING id;
    `
	// O ID retornado é setado de volta no struct 'destination'
	err := r.db.QueryRow(query,
		destination.TravelGroupID,
		destination.Name,
		destination.Location,
		destination.Description,
	).Scan(&destination.ID)

	if err != nil {
		return fmt.Errorf("erro ao inserir destino: %w", err)
	}
	return nil
}

func (r *postgresTravelGroupRepository) CreateVoting(groupID int, question string, optionsJSON string) (int, error) {
	var newID int
	query := `
        INSERT INTO votings 
        (travel_group_id, question, options, created_at) 
        VALUES 
        ($1, $2, $3, NOW())
        RETURNING id;
    `
	err := r.db.QueryRow(query,
		groupID,
		question,
		optionsJSON,
	).Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("erro ao inserir votação: %w", err)
	}
	return newID, nil
}

func (r *postgresTravelGroupRepository) CreateExpense(expense *models.Expense) error {

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação para despesa: %w", err)
	}
	defer tx.Rollback()

	expenseQuery := `
        INSERT INTO expenses 
        (travel_group_id, description, amount, payer_id, created_at) 
        VALUES 
        ($1, $2, $3, $4, NOW())
        RETURNING id;
    `
	err = tx.QueryRow(expenseQuery,
		expense.TravelGroupID,
		expense.Description,
		expense.Amount,
		expense.PayerID,
	).Scan(&expense.ID)

	if err != nil {
		return fmt.Errorf("erro ao inserir despesa: %w", err)
	}

	if len(expense.ParticipantIDs) > 0 {
		participantQuery := `INSERT INTO expense_participants (expense_id, user_id) VALUES ($1, $2)`

		for _, userID := range expense.ParticipantIDs {
			_, err := tx.Exec(participantQuery, expense.ID, userID)
			if err != nil {
				return fmt.Errorf("erro ao inserir participante %d para despesa %d: %w", userID, expense.ID, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("falha ao commitar transação da despesa: %w", err)
	}

	return nil
}
