package models

// VoteRequest é o payload para registrar um voto
type VoteRequest struct {
	SelectedOption string `json:"selectedOption"`
}

// Vote Model (para registro interno)
type Vote struct {
	VotingID       int
	UserID         int
	SelectedOption string
}
