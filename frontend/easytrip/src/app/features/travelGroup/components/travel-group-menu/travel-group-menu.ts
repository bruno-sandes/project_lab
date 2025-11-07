import { Component, inject, signal, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { CommonModule } from '@angular/common';
import { forkJoin, catchError, of, switchMap, tap } from 'rxjs';
import { DestinationDTO, ExpenseDTO, GroupMemberDTO, TravelGroupDetails, VotingDTO } from '../../../groups-dashboard/models/travel_groups';
import { TravelGroupsService } from '../../../groups-dashboard/services/travel-groups-service';

type TabType = 'members' | 'destinations' | 'votings' | 'expenses';

@Component({
  selector: 'app-travel-group-menu',
  imports: [CommonModule], 
  templateUrl: './travel-group-menu.html',
  styleUrl: './travel-group-menu.css'
})
export class TravelGroupMenu implements OnInit {
  
  // Dependências
  private activatedRouteService = inject(ActivatedRoute);
  private groupsService = inject(TravelGroupsService);

  // Estado
  groupId = signal<number | null>(null);
  isLoading = signal<boolean>(true);
  errorMessage = signal<string | null>(null);
  activeTab = signal<TabType>('members');

  // Dados do Grupo Principal
  groupName = signal<string>('');
  organizerName = signal<string>('');
  
  // Dados dos Sub-recursos (usando modelos reais)
  members = signal<GroupMemberDTO[]>([]);
  destinations = signal<DestinationDTO[]>([]);
  votings = signal<VotingDTO[]>([]);
  expenses = signal<ExpenseDTO[]>([]);
  

  // Array de abas para a navegação
  tabs: { key: TabType, label: string }[] = [
    { key: 'members', label: 'Membros' },
    { key: 'destinations', label: 'Destinos' },
    { key: 'votings', label: 'Votações' },
    { key: 'expenses', label: 'Despesas' },
  ];
  

  ngOnInit(): void {
    const idStr = this.activatedRouteService.snapshot.paramMap.get('travelGroupId');
    const id = idStr ? parseInt(idStr, 10) : null;

    if (id === null || isNaN(id)) {
      this.errorMessage.set('ID do grupo inválido na URL.');
      this.isLoading.set(false);
      return;
    }
    
    this.groupId.set(id);
    this.loadGroupData(id);
  }

  /**
   * Carrega os dados principais do grupo e todos os sub-recursos em paralelo.
   */
  loadGroupData(groupId: number): void {
    this.isLoading.set(true);
    this.errorMessage.set(null);

    this.groupsService.getGroupDetails(groupId).pipe(
      switchMap((details: TravelGroupDetails) => {
        this.updateGroupHeader(details);

        return forkJoin({
          members: this.groupsService.listGroupMembers(groupId),
          destinations: this.groupsService.listDestinations(groupId),
          votings: this.groupsService.listVotings(groupId),
          expenses: this.groupsService.listExpenses(groupId),
        });
      }),
      catchError(err => {
        this.errorMessage.set('Falha ao carregar todos os dados do grupo. Verifique sua conexão e permissões.');
        console.error('Erro ao carregar dados do grupo:', err);
        return of(null);
      }),
      tap(() => this.isLoading.set(false))
    ).subscribe((data) => {
      if (data) {
        this.members.set(data.members);
        this.destinations.set(data.destinations);
        this.votings.set(data.votings);
        this.expenses.set(data.expenses);
      }
    });
  }

  /**
   * Atualiza os signals de nome e organizador.
   */
  private updateGroupHeader(details: TravelGroupDetails): void {
    this.groupName.set(details.name);
    this.organizerName.set(`Usuário ID ${details.creatorId}`);
  }

  // Método para trocar de aba
  setActiveTab(tab: TabType): void {
    this.activeTab.set(tab);
  }
}