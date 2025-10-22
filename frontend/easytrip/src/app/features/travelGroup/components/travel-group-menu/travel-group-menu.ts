import { Component, inject, signal } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

interface MemberListItem {
  id: number;
  name: string;
  email: string;
  role: 'Organizador' | 'Participante';
  avatarUrl: string;
}

interface Destination {
  id: number;
  name: string;
  location: string;
  description: string;
}

interface Voting {
  id: number;
  question: string;
  options: string[];
  status: 'Aberto' | 'Fechado';
  votedOption?: string; 
}

interface Expense {
  id: number;
  description: string;
  amount: number;
  payerName: string;
  participantsCount: number; 
}

type TabType = 'members' | 'destinations' | 'votings' | 'expenses';


@Component({
  selector: 'app-travel-group-menu',
  imports: [],
  templateUrl: './travel-group-menu.html',
  styleUrl: './travel-group-menu.css'
})
export class TravelGroupMenu {
// Signals para o estado da página
  groupName = signal<string>('Viagem de Verão para a Europa');
  organizerName = signal<string>('Sofia Mendes');
  activeTab = signal<TabType>('members');

  // Signals para os dados (usando arrays vazios como placeholder)
  members = signal<MemberListItem[]>([]);
  destinations = signal<Destination[]>([]);
  votings = signal<Voting[]>([]);
  expenses = signal<Expense[]>([]);

  // Array de abas para a navegação
  tabs: { key: TabType, label: string }[] = [
    { key: 'members', label: 'Membros' },
    { key: 'destinations', label: 'Destinos' },
    { key: 'votings', label: 'Votações' },
    { key: 'expenses', label: 'Despesas' },
  ];
  
  activatedRouteService = inject(ActivatedRoute)


  ngOnInit(): void {
    const groupId = this.activatedRouteService.snapshot.paramMap.get('travelGroupId');
    console.log(`Carregando dados para o Grupo ID: ${groupId}`);
    
    this.loadGroupData(groupId);
  }

  loadGroupData(groupId: string | null) {
    //TO DO: Implementar chamadas a API aqui (GET /groups/{id}/...)

    // MOCK DATA 
    this.members.set([
      { id: 1, name: 'Sofia Mendes', email: 'sofia@email.com', role: 'Organizador', avatarUrl: '...url1' },
      { id: 2, name: 'Carlos Pereira', email: 'carlos@email.com', role: 'Participante', avatarUrl: '...url2' },
      { id: 3, name: 'Ana Silva', email: 'ana@email.com', role: 'Participante', avatarUrl: '...url3' },
      { id: 4, name: 'Ricardo Almeida', email: 'ricardo@email.com', role: 'Participante', avatarUrl: '...url4' },
    ]);

    // Vazio para testar o @empty
    this.destinations.set([]); 
    this.votings.set([]);
    this.expenses.set([]);
  }

  // Método para trocar de aba
  setActiveTab(tab: TabType): void {
    this.activeTab.set(tab);
  }

}
