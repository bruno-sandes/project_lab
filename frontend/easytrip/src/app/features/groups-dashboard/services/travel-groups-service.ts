import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { DestinationCreateRequest, DestinationDTO, ExpenseCreateRequest, ExpenseDTO, GroupMemberDTO, MemberAddRequest, TravelGroupCreateRequest, TravelGroupDetails, TravelGroupListItem, VoteRequest, VotingCreateRequest, VotingDTO } from '../models/travel_groups';

@Injectable({
  providedIn: 'root'
})
export class TravelGroupsService {
  private baseUrl = 'http://localhost:8080'; // Base URL para rotas fora de /groups
  private groupsApiUrl = `${this.baseUrl}/groups`; 

  http = inject(HttpClient);


  // --- CRUD BÁSICO DE GRUPOS (/groups) ---

  /**
   * GET /groups
   * Lista todos os grupos aos quais o usuário logado pertence ou criou.
   */
  listGroups(): Observable<TravelGroupListItem[]> {
    return this.http.get<TravelGroupListItem[]>(this.groupsApiUrl);
  }

  /**
   * POST /groups
   * Cria um novo grupo de viagem.
   */
  createGroup(groupData: TravelGroupCreateRequest): Observable<TravelGroupDetails> {
    return this.http.post<TravelGroupDetails>(this.groupsApiUrl, groupData);
  }
  
  /**
   * GET /groups/{id}
   * Obtém os detalhes de um grupo específico.
   */
  getGroupDetails(groupId: number): Observable<TravelGroupDetails> {
    return this.http.get<TravelGroupDetails>(`${this.groupsApiUrl}/${groupId}`);
  }


  // --- SUB-RECURSOS DE GRUPO (/groups/{id}/...) ---

  // 1. MEMBROS
  /**
   * GET /groups/{id}/members
   * Lista os membros de um grupo.
   */
  listGroupMembers(groupId: number): Observable<GroupMemberDTO[]> {
    return this.http.get<GroupMemberDTO[]>(`${this.groupsApiUrl}/${groupId}/members`);
  }

  /**
   * POST /groups/{id}/members
   * Adiciona um membro ao grupo.
   */
  addMember(groupId: number, userData: MemberAddRequest): Observable<void> {
    // Usamos `void` porque o backend retorna 204 No Content
    return this.http.post<void>(`${this.groupsApiUrl}/${groupId}/members`, userData);
  }


  // 2. DESTINOS
  /**
   * GET /groups/{id}/destinations
   * Lista os destinos sugeridos para um grupo.
   */
  listDestinations(groupId: number): Observable<DestinationDTO[]> {
    return this.http.get<DestinationDTO[]>(`${this.groupsApiUrl}/${groupId}/destinations`);
  }

  /**
   * POST /groups/{id}/destinations
   * Cria um novo destino.
   */
  createDestination(groupId: number, destinationData: DestinationCreateRequest): Observable<DestinationDTO> {
    return this.http.post<DestinationDTO>(`${this.groupsApiUrl}/${groupId}/destinations`, destinationData);
  }


  // 3. VOTAÇÕES
  /**
   * GET /groups/{id}/votings
   * Lista as votações de um grupo.
   */
  listVotings(groupId: number): Observable<VotingDTO[]> {
    return this.http.get<VotingDTO[]>(`${this.groupsApiUrl}/${groupId}/votings`);
  }

  /**
   * POST /groups/{id}/votings
   * Cria uma nova votação.
   */
  createVoting(groupId: number, votingData: VotingCreateRequest): Observable<{ id: number }> {
    return this.http.post<{ id: number }>(`${this.groupsApiUrl}/${groupId}/votings`, votingData);
  }


  // 4. DESPESAS
  /**
   * GET /groups/{id}/expenses
   * Lista as despesas de um grupo.
   */
  listExpenses(groupId: number): Observable<ExpenseDTO[]> {
    return this.http.get<ExpenseDTO[]>(`${this.groupsApiUrl}/${groupId}/expenses`);
  }

  /**
   * POST /groups/{id}/expenses
   * Cria uma nova despesa.
   */
  createExpense(groupId: number, expenseData: ExpenseCreateRequest): Observable<ExpenseDTO> {
    // O backend retorna o DTO da despesa criada (com ID e timestamps)
    return this.http.post<ExpenseDTO>(`${this.groupsApiUrl}/${groupId}/expenses`, expenseData);
  }

  
  // --- AÇÃO DE VOTO (/votings/{id}/vote) ---

  /**
   * POST /votings/{id}/vote
   * Registra o voto do usuário em uma votação.
   */
  castVote(votingId: number, voteData: VoteRequest): Observable<void> {
    // Usamos `void` pois o backend retorna 201 com um corpo simples de mensagem.
    return this.http.post<void>(`${this.baseUrl}/votings/${votingId}/vote`, voteData);
  }
}