package repositories

import (
	"database/sql"
	"errors"
	"project_lab/internal/models"

	"github.com/lib/pq"
)

// ErrEmailAlreadyExists é um erro customizado para duplicidade de e-mail.
var ErrEmailAlreadyExists = errors.New("e-mail já está em uso")

// UserRepository é a interface que define os métodos de acesso a dados para usuários.
type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	GetUserProfile(userID int) (*models.UserProfileResponse, error)
	UpdateUserName(userID int, newName string) error
}

// userRepository representa a implementação do repositório com o banco de dados.
type userRepository struct {
	db *sql.DB
}

// NewUserRepository cria uma nova instância de UserRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// CreateUser insere um novo usuário no banco de dados.
func (r *userRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrEmailAlreadyExists
		}
		return errors.New("erro ao criar usuário: " + err.Error())
	}
	return nil
}

// FindByEmail busca um usuário no banco de dados por e-mail.
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, password_hash FROM users WHERE email = $1`
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}
	return user, nil
}

// GetUserProfile busca o nome e email do usuário pelo ID.
func (r *userRepository) GetUserProfile(userID int) (*models.UserProfileResponse, error) {
	var profile models.UserProfileResponse
	query := `SELECT name, email FROM users WHERE id = $1`

	err := r.db.QueryRow(query, userID).Scan(&profile.Name, &profile.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			// Se o usuário não existir
			return nil, errors.New("usuário não encontrado")
		}
		// Outros erros de banco de dados
		return nil, errors.New("erro ao buscar perfil do usuário: " + err.Error())
	}
	return &profile, nil
}

// UpdateUserName atualiza o nome do usuário.
func (r *userRepository) UpdateUserName(userID int, newName string) error {
	query := `
		UPDATE users 
		SET name = $2, updated_at = NOW()
		WHERE id = $1
	`
	// Note: Eu corrigi o caractere inválido ' ' que estava no seu código original.
	result, err := r.db.Exec(query, userID, newName)
	if err != nil {
		return errors.New("erro ao atualizar nome do usuário: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("erro ao verificar linhas afetadas: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("usuário não encontrado para atualização")
	}

	return nil
}
