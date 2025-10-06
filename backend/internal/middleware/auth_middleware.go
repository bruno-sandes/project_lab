package middleware

import (
	"context"
	"net/http"
	"project_lab/internal/services"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware é um middleware que protege rotas para verificar o token JWT .
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token de autenticação é necessário.", http.StatusUnauthorized)
			return
		}

		// O formato é esperado: "Bearer [TOKEN]"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato do token inválido. Use 'Bearer <token>'.", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		//Valida e faz o parse do token
		claims := &services.UserClaims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				// Retorna a chave secreta usada para assinar
				return []byte(services.JwtSecret), nil
			},
		)

		if err != nil || !token.Valid {
			http.Error(w, "Token inválido ou expirado.", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
