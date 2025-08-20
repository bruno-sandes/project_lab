package models

// User representa o modelo de dados de um usuário.
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"` // A senha sempre deve ser armazenada com hash!
}
