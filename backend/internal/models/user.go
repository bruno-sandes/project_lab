package models

// User representa o modelo de dados de um usuário.
type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"-"` // O hash não deve ser exposto no JSON
}

// UserLogin representa a requisição de login.
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserProfileResponse corresponde ao DTO retornado por GET /profile
type UserProfileResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserProfileUpdateRequest corresponde ao payload de PATCH /profile
type UserProfileUpdateRequest struct {
	Name string `json:"name"`
}
