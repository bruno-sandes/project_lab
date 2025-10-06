package services

import (
	"errors"
	"project_lab/internal/models"
	"project_lab/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("erro ao gerar hash da senha")
	}
	user.PasswordHash = string(hashedPassword)

	if err := s.userRepo.CreateUser(user); err != nil {
		return err
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

	//  Gerar e retornar o JWT
	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", errors.New("falha ao gerar token de autenticação")
	}
	return token, nil
}

func (s *authService) generateToken(userID int) (string, error) {
	// Define o tempo de expiração do token (Nesse casos 24 horas)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Cria os dados do token
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Cria o token usando o algoritmo de assinatura (HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assina o token com a chave secreta
	tokenString, err := token.SignedString([]byte(JwtSecret))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
