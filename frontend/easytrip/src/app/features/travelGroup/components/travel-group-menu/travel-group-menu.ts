import { Component, inject, signal, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { CommonModule } from '@angular/common';
import { forkJoin, catchError, of, switchMap, tap } from 'rxjs';
import { DestinationDTO, ExpenseDTO, GroupMemberDTO, TravelGroupDetails, VotingDTO } from '../../../groups-dashboard/models/travel_groups';
import { TravelGroupsService } from '../../../groups-dashboard/services/travel-groups-service';
import { AuthService } from '../../../login/service/login-service';
import { CreatDestinationModal } from "../../../destinations/components/creat-destination-modal/creat-destination-modal";
import { CreateVotingModal } from "../../../voting/components/create-voting-modal/create-voting-modal";
import { CreateExpensesModal } from "../../../expenses/components/create-expenses-modal/create-expenses-modal";
import { VoteModal } from '../../../voting/components/vote-modal/vote-modal';

type TabType = 'members' | 'destinations' | 'votings' | 'expenses';

@Component({
  selector: 'app-travel-group-menu',
  imports: [CommonModule, CreatDestinationModal, CreateVotingModal, CreateExpensesModal, VoteModal], 
  templateUrl: './travel-group-menu.html',
  styleUrl: './travel-group-menu.css'
})
export class TravelGroupMenu implements OnInit {
  
  private activatedRouteService = inject(ActivatedRoute);
  private groupsService = inject(TravelGroupsService);
  private authService = inject(AuthService); 

  // Estado
  groupId = signal<number | null>(null);
  isLoading = signal<boolean>(true);
  errorMessage = signal<string | null>(null);
  activeTab = signal<TabType>('members');
  loggedUserId = signal<number>(0);

  // Estado das Modais
  showCreateDestination = signal(false);
  showCreateVoting = signal(false);
  showCreateExpense = signal(false);
  showVoteModal = signal(false); 
  selectedVoting = signal<VotingDTO | null>(null); 
  
  // Dados do Grupo Principal
  groupName = signal<string>('');
  organizerName = signal<string>('');
  
  // Dados dos Sub-recursos
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

  private updateGroupHeader(details: TravelGroupDetails): void {
    this.groupName.set(details.name);
 }

  setActiveTab(tab: TabType): void {
    this.activeTab.set(tab);
  }

  openModal(modalType: 'destination' | 'voting' | 'expense'): void {
    if (modalType === 'destination') {
      this.showCreateDestination.set(true);
    } else if (modalType === 'voting') {
      this.showCreateVoting.set(true);
    } else if (modalType === 'expense') {
      this.showCreateExpense.set(true);
    }
  }

  openVoteModal(voting: VotingDTO): void {
    this.selectedVoting.set(voting);
    this.showVoteModal.set(true);
  }

  closeAllModals(): void {
    this.showCreateDestination.set(false);
    this.showCreateVoting.set(false);
    this.showCreateExpense.set(false);
    this.showVoteModal.set(false);
    this.selectedVoting.set(null); // Limpa o estado
  }

  handleCreationSuccess(): void {
    this.closeAllModals();
    if (this.groupId()) {
      this.loadGroupData(this.groupId()!); 
    }
  }
}