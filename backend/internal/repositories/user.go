package repositories

import (
	"database/sql"
	"errors"
	"project_lab/internal/models"
)

// UserRepository é a interface que define os métodos de acesso a dados para usuários.
type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
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

// FindByEmail busca um usuário no banco de dados por e-mail.
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	// A query SQL para buscar o usuário.
	query := `SELECT id, email, password FROM users WHERE email = $1`

	user := &models.User{}

	// Executa a query e escaneia o resultado para a struct User.
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// Se não houver linhas, o usuário não foi encontrado.
			return nil, errors.New("usuário não encontrado")
		}
		// Se for outro erro, retorne a mensagem de erro.
		return nil, err
	}

	return user, nil
}
