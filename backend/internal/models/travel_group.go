package models

import (
	"time"
)

// TravelGroup é a estrutura principal que mapeia a tabela travel_groups.
type TravelGroup struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	CreatorID   int       `json:"creator_id"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
}

type TravelGroupCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

// TravelGroupListItem é a estrutura simplificada para a tela de listagem.
// Inclui o nome e id do criador
type TravelGroupListItem struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	MemberCount int       `json:"member_count"`
	CreatorId   int       `json:"creator_id"`
	CreatorName string    `json:"creator_name"`
}
