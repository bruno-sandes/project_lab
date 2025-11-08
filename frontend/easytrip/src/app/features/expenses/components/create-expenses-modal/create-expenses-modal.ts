import { Component, EventEmitter, inject, Input, Output, signal, OnInit } from '@angular/core';
import { FormGroup, FormControl, Validators, FormArray, ReactiveFormsModule } from '@angular/forms';
import { GroupMemberDTO, ExpenseCreateRequest } from '../../../groups-dashboard/models/travel_groups';
import { TravelGroupsService } from '../../../groups-dashboard/services/travel-groups-service';
import { CommonModule } from '@angular/common'; // Necessário para @if/@for

@Component({
  selector: 'app-create-expenses-modal',
  standalone: true,
  imports: [ReactiveFormsModule, CommonModule],
  templateUrl: './create-expenses-modal.html',
  styleUrl: './create-expenses-modal.css'
})
export class CreateExpensesModal implements OnInit { // IMPLEMENTADO ONINIT
  private groupsService = inject(TravelGroupsService);

  @Input({ required: true }) groupId!: number;
  @Input({ required: true }) members: GroupMemberDTO[] = [];
  @Input({ required: true }) loggedUserId!: number; 
  
  @Output() close = new EventEmitter<void>();
  @Output() expenseCreated = new EventEmitter<void>();

  isLoading = signal(false);
  errorMessage = signal<string | null>(null);

  expenseForm = new FormGroup({
    description: new FormControl('', { 
      nonNullable: true, 
      validators: [Validators.required, Validators.minLength(3)] 
    }),
    amount: new FormControl(0, { 
      nonNullable: true, 
      validators: [Validators.required, Validators.min(0.01)]
    }),
    payerId: new FormControl(0, {
      nonNullable: true, 
      validators: [Validators.required]
    }),
    participantIds: new FormArray<FormControl<boolean>>([], []) 
  });

  ngOnInit(): void {
    if (this.loggedUserId <= 0 || this.members.length === 0) {
        return;
    }

    // Inicializa o FormArray com um checkbox para cada membro, marcado como true por padrão
    this.members.forEach(() => {
        this.participantIdsArray.push(new FormControl(true, { nonNullable: true })); 
    });
    
    // Define o pagador padrão como o usuário logado
    this.expenseForm.controls.payerId.setValue(this.loggedUserId);
    this.expenseForm.updateValueAndValidity(); // Força a revalidação após o preenchimento inicial
  }

  get participantIdsArray() {
    return this.expenseForm.get('participantIds') as FormArray<FormControl<boolean>>;
  }

  get hasSelectedParticipants(): boolean {
    return this.participantIdsArray.controls.some(control => control.value);
  }
  
  // Getter unificado para desabilitar o botão: checa campos principais E participantes
  get isFormInvalid(): boolean {
      return this.expenseForm.invalid 
  }

  onSubmit() {
    // Primeira verificação de validade e estado de carregamento
    if (this.isFormInvalid || this.isLoading()) { 
      this.errorMessage.set('Verifique os campos. Descrição, Valor, Pagador e pelo menos um participante são obrigatórios.');
      this.expenseForm.markAllAsTouched(); 
      return;
    }
    
    this.isLoading.set(true);
    this.errorMessage.set(null);

    // Mapeia apenas os IDs dos membros que estão com o checkbox marcado
    const selectedParticipantIds = this.members
        .filter((_, index) => this.participantIdsArray.controls[index].value)
        .map(member => member.userId);

    const data: ExpenseCreateRequest = {
        description: this.expenseForm.controls.description.value,
        amount: this.expenseForm.controls.amount.value,
        payerId: this.expenseForm.controls.payerId.value,
        participantIds: selectedParticipantIds 
    };

    // Chamada do serviço para enviar a despesa
    this.groupsService.createExpense(this.groupId, data)
      .subscribe({
        next: () => {
          this.expenseCreated.emit();
          this.close.emit();
        },
        error: (err) => {
          this.isLoading.set(false);
          this.errorMessage.set('Falha ao registrar despesa. Verifique o console para detalhes.');
          console.error('Erro ao criar despesa:', err);
        }
      });
  }

  onClose() {
    this.close.emit();
  }
}