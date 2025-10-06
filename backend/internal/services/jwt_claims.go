package services

import (
	"github.com/golang-jwt/jwt/v5"
)

const JwtSecret = "sua_chave_secreta_muito_forte_aqui"

// UserClaims define a estrutura dos dados que serão armazenados no token JWT.
type UserClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
