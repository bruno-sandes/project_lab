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

// TravelGroupDetails representa os dados básicos de um grupo para a página de visualização
type TravelGroupDetails struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	CreatorID   int       `json:"creatorId"`
	CreatorName string    `json:"organizerName"`
	MemberCount int       `json:"memberCount"`
}

// GroupMemberDTO representa um item na lista de membros (para a aba Membros)
type GroupMemberDTO struct {
	UserID int    `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// DestinationDTO representa um destino sugerido para um grupo
type DestinationDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type VotingDTO struct {
	ID         int       `json:"id"`
	Question   string    `json:"question"`
	Options    []string  `json:"options"`
	TotalVotes int       `json:"totalVotes"`
	UserVote   *string   `json:"userVote"`
	CreatedAt  time.Time `json:"createdAt"`
	// Status (Aberto/Fechado) pode ser inferido pelo backend ou adicionado aqui.
}

// ExpenseDTO representa uma despesa do grupo
type ExpenseDTO struct {
	ID                int       `json:"id"`
	Description       string    `json:"description"`
	Amount            float64   `json:"amount"`
	PayerID           int       `json:"payerId"`
	PayerName         string    `json:"payerName"`
	ParticipantsIDs   []int     `json:"participantsIds"`
	ParticipantsCount int       `json:"participantsCount"`
	CreatedAt         time.Time `json:"createdAt"`
}

// DestinationCreateRequest é o payload para criar um novo destino
type DestinationCreateRequest struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

// VotingCreateRequest é o payload para criar uma nova votação
type VotingCreateRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

// Destination Model (para passar para o repository se necessário)
type Destination struct {
	ID            int
	TravelGroupID int
	Name          string
	Location      string
	Description   string
}

// ExpenseCreateRequest é o payload para criar uma nova despesa
type ExpenseCreateRequest struct {
	Description    string  `json:"description"`
	Amount         float64 `json:"amount"`
	PayerID        int     `json:"payerId"`
	ParticipantIDs []int   `json:"participantIds"`
}

// Expense Model (para uso interno no Repository)
type Expense struct {
	ID             int
	TravelGroupID  int
	Description    string
	Amount         float64
	PayerID        int
	ParticipantIDs []int
}
