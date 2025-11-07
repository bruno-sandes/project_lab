// src/app/core/models/travel_groups.model.ts

// ====================================================================
// --- TIPOS DE GRUPOS E DETALHES ---
// ====================================================================

/**
 * Interface para os dados de criação de um novo grupo.
 * Corresponde ao TravelGroupCreateRequest no YAML.
 */
export interface TravelGroupCreateRequest {
  name: string;
  description?: string;
  start_date: string; // 'YYYY-MM-DD'
  end_date: string;   // 'YYYY-MM-DD'
}

/**
 * Interface para a resposta de listagem de grupos.
 * Corresponde ao TravelGroupListItem no YAML.
 */
export interface TravelGroupListItem {
  id: number;
  name: string;
  description?: string;
  start_date: string; 
  end_date: string;   
  memberCount: number;
  creatorId: number;
  creatorName: string;
}

/**
 * Interface para a resposta de detalhes de um grupo.
 * Corresponde ao TravelGroupDetails no YAML.
 */
export interface TravelGroupDetails {
  id: number;
  name: string;
  description?: string;
  creatorId: number;
  start_date: string;
  end_date: string;
  createdAt: string;
}


// ====================================================================
// --- TIPOS DE MEMBROS E AÇÕES (/members) ---
// ====================================================================

/**
 * DTO para listar os membros de um grupo.
 * Corresponde ao GroupMemberDTO no YAML.
 */
export interface GroupMemberDTO {
    userId: number;
    name: string;
    email: string;
    role: 'Organizador' | 'Participante';
}

/**
 * Payload para adicionar um membro a um grupo.
 * Corresponde ao MemberAddRequest no YAML.
 */
export interface MemberAddRequest {
    userId: number;
}


// ====================================================================
// --- TIPOS DE DESTINOS (/destinations) ---
// ====================================================================

/**
 * DTO para listar ou retornar um destino criado.
 * Corresponde ao DestinationDTO no YAML.
 */
export interface DestinationDTO {
    id: number;
    name: string;
    location: string;
    description: string;
}

/**
 * Payload para criar um novo destino.
 * Corresponde ao DestinationCreateRequest no YAML.
 */
export interface DestinationCreateRequest {
    name: string;
    location: string;
    description: string;
}


// ====================================================================
// --- TIPOS DE VOTAÇÕES E VOTOS (/votings, /vote) ---
// ====================================================================

/**
 * DTO para listar as votações de um grupo.
 * Corresponde ao VotingDTO no YAML.
 */
export interface VotingDTO {
    id: number;
    question: string;
    options: string[];
    totalVotes: number;
    userVote: string | null;
    createdAt: string;
}

/**
 * Payload para criar uma nova votação.
 * Corresponde ao VotingCreateRequest no YAML.
 */
export interface VotingCreateRequest {
    question: string;
    options: string[];
}

/**
 * Payload para registrar um voto.
 * Corresponde ao VoteRequest no YAML.
 */
export interface VoteRequest {
    selectedOption: string;
}


// ====================================================================
// --- TIPOS DE DESPESAS (/expenses) ---
// ====================================================================

/**
 * DTO para listar ou retornar uma despesa.
 * Corresponde ao ExpenseDTO no YAML.
 */
export interface ExpenseDTO {
    id: number;
    description: string;
    amount: number;
    payerId: number;
    payerName: string;
    participantsIds: number[];
    participantsCount: number;
    createdAt: string;
}

/**
 * Payload para criar uma nova despesa.
 * Corresponde ao ExpenseCreateRequest no YAML.
 */
export interface ExpenseCreateRequest {
    description: string;
    amount: number;
    payerId: number;
    participantIds: number[];
}

/**
 * DTO para a resposta da criação de uma despesa.
 * Corresponde ao ExpenseResponse no YAML (que usa allOf com ExpenseDTO).
 */
export interface ExpenseResponse extends ExpenseDTO {}