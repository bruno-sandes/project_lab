import { Component, EventEmitter, inject, Input, Output, signal, OnInit } from '@angular/core';
import { FormControl, ReactiveFormsModule, Validators } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { TravelGroupsService } from '../../../groups-dashboard/services/travel-groups-service'; // Ajuste o caminho conforme sua estrutura

export interface VotingDTO {
    id: number;
    question: string;
    options: string[]; 
    userVote: string | null;
    totalVotes: number;
}

export interface VoteRequest {
    votingId: number;
    selectedOption: string;
}

@Component({
  selector: 'app-vote-modal',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  styleUrl:'./vote-modal.css',
  templateUrl: './vote-modal.html',
})
export class VoteModal implements OnInit {
  
  private groupsService = inject(TravelGroupsService);

  @Input({ required: true }) groupId!: number;
  @Input({ required: true }) voting!: VotingDTO; 
  
  @Output() close = new EventEmitter<void>();
  @Output() voteSubmitted = new EventEmitter<void>();

  isLoading = signal(false);
  errorMessage = signal<string | null>(null);

  voteControl = new FormControl<string>('', {
    nonNullable: true,
    validators: [Validators.required]
  });

  ngOnInit(): void {
      if (this.voting.userVote) {
          this.voteControl.setValue(this.voting.userVote);
      }
  }

  onSubmit() {
    if (this.voteControl.invalid || this.isLoading()) {
      this.errorMessage.set('Selecione uma opção válida para registrar seu voto.');
      this.voteControl.markAsTouched();
      return;
    }
    
    this.isLoading.set(true);
    this.errorMessage.set(null);

    const data: VoteRequest = {
        votingId: this.voting.id,
        selectedOption: this.voteControl.value
    };

    this.groupsService.castVote(this.voting.id, data) 
      .subscribe({
        next: () => {
          this.voteSubmitted.emit();
          this.close.emit();
        },
        error: (err) => {
          this.isLoading.set(false);
          this.errorMessage.set('Falha ao registrar voto. Tente novamente.');
          console.error('Erro ao votar:', err);
        }
      });
  }

  onClose() {
    this.close.emit();
  }
}