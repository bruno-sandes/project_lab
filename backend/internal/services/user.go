package services

import (
	"errors"
	"project_lab/internal/models"
	"project_lab/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

// AuthService é a interface que define a lógica de negócio de autenticação.
type AuthService interface {
	RegisterUser(user *models.User) error
	Authenticate(email, password string) (string, error)
}

// authService implementa a interface AuthService.
type authService struct {
	userRepo repositories.UserRepository
}

// NewAuthService cria uma nova instância de AuthService.
func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// RegisterUser lida com a lógica de negócio do cadastro.
func (s *authService) RegisterUser(user *models.User) error {
	// 1. Gera o hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("erro ao gerar hash da senha")
	}
	user.PasswordHash = string(hashedPassword)

	// 2. Salva o usuário no banco de dados
	if err := s.userRepo.CreateUser(user); err != nil {
		return err // Passa o erro (incluindo o de conflito) para o handler
	}

	return nil
}

// Authenticate autentica um usuário.
func (s *authService) Authenticate(email, password string) (string, error) {
	//Busca o usuário pelo e-mail
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("usuário ou senha incorretos")
	}

	//Compara a senha com o hash no banco
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("usuário ou senha incorretos")
	}

	// Lógica de sucesso: Gerar e retornar o JWT
	// Por enquanto, retorna uma string de teste
	return "token_de_teste", nil
}
