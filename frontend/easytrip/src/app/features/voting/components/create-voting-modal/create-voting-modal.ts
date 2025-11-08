import { Component, EventEmitter, inject, Input, Output, signal } from '@angular/core';
import { FormGroup, FormControl, Validators, FormArray, ReactiveFormsModule } from '@angular/forms';
import { VotingCreateRequest } from '../../../groups-dashboard/models/travel_groups';
import { TravelGroupsService } from '../../../groups-dashboard/services/travel-groups-service';

@Component({
  selector: 'app-create-voting-modal',
  imports: [ReactiveFormsModule],
  templateUrl: './create-voting-modal.html',
  styleUrl: './create-voting-modal.css'
})
export class CreateVotingModal {
private groupsService = inject(TravelGroupsService);

  @Input({ required: true }) groupId!: number;
  
  @Output() close = new EventEmitter<void>();
  @Output() votingCreated = new EventEmitter<void>();

  isLoading = signal(false);
  errorMessage = signal<string | null>(null);

  votingForm = new FormGroup({
    question: new FormControl('', { 
      nonNullable: true, 
      validators: [Validators.required, Validators.minLength(5)] 
    }),
    options: new FormArray([
      new FormControl('', { nonNullable: true, validators: [Validators.required] }),
      new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    ], Validators.minLength(2))
  });

  get questionControl() {
    return this.votingForm.get('question') as FormControl;
  }
  
  get optionsArray() {
    return this.votingForm.get('options') as FormArray<FormControl<string>>;
  }

  addOption() {
    this.optionsArray.push(new FormControl('', { 
        nonNullable: true, 
        validators: [Validators.required] 
    }));
  }

  removeOption(index: number) {
    if (this.optionsArray.length > 2) {
      this.optionsArray.removeAt(index);
    }
  }

  onSubmit() {
    if (this.votingForm.invalid || this.isLoading()) {
      return;
    }
    
    this.isLoading.set(true);
    this.errorMessage.set(null);

    const options = this.optionsArray.controls
        .map(control => control.value)
        .filter(value => value.trim() !== '');

    const data: VotingCreateRequest = {
        question: this.votingForm.controls.question.value,
        options: options
    };

    this.groupsService.createVoting(this.groupId, data)
      .subscribe({
        next: () => {
          this.votingCreated.emit();
          this.close.emit();
        },
        error: (err) => {
          this.isLoading.set(false);
          this.errorMessage.set('Falha ao criar votação. Verifique a pergunta e opções.');
          console.error('Erro ao criar votação:', err);
        }
      });
  }

  onClose() {
    this.close.emit();
  }
}
