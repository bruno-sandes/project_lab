package services

import (
	"errors"
	"project_lab/internal/models"
	"project_lab/internal/repositories"
)

// AuthService é a interface que define a lógica de negócio de autenticação.
type AuthService interface {
	Authenticate(email, password string) (*models.User, error)
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

// Authenticate autentica um usuário.
func (s *authService) Authenticate(email, password string) (*models.User, error) {
	// 1. Busca o usuário pelo e-mail na camada de Repositories
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("usuário ou senha incorretos")
	}

	// 2. Compara a senha (aqui você usaria uma biblioteca de hash como bcrypt)
	// Por enquanto, vamos fazer uma comparação simples para focar na lógica
	if user.Password != password {
		return nil, errors.New("usuário ou senha incorretos")
	}

	// Lógica de sucesso: aqui você geraria um JWT
	// Por agora, apenas retornamos o usuário
	return user, nil
}
