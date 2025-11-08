import { Component, EventEmitter, inject, Input, Output, signal } from '@angular/core';
import { FormGroup, FormControl, Validators, ReactiveFormsModule } from '@angular/forms';
import { DestinationCreateRequest } from '../../../groups-dashboard/models/travel_groups';
import { TravelGroupsService } from '../../../groups-dashboard/services/travel-groups-service';

@Component({
  selector: 'app-creat-destination-modal',
  imports: [ReactiveFormsModule],
  templateUrl: './creat-destination-modal.html',
  styleUrl: './creat-destination-modal.css'
})
export class CreatDestinationModal {
private groupsService = inject(TravelGroupsService);

  // ID do grupo é obrigatório para a criação
  @Input({ required: true }) groupId!: number;
  
  // Eventos para fechar a modal ou notificar sucesso
  @Output() close = new EventEmitter<void>();
  @Output() destinationCreated = new EventEmitter<void>();

  isLoading = signal(false);
  errorMessage = signal<string | null>(null);

  destinationForm = new FormGroup({
    name: new FormControl('', { 
      nonNullable: true, 
      validators: [Validators.required, Validators.minLength(3)] 
    }),
    location: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    description: new FormControl('', { nonNullable: true }),
  });

  // Getter de conveniência
  get nameControl() {
    return this.destinationForm.get('name') as FormControl;
  }
  
  get locationControl() {
    return this.destinationForm.get('location') as FormControl;
  }

  onSubmit() {
    if (this.destinationForm.invalid || this.isLoading()) {
      return;
    }
    
    this.isLoading.set(true);
    this.errorMessage.set(null);

    const data: DestinationCreateRequest = this.destinationForm.getRawValue();

    this.groupsService.createDestination(this.groupId, data)
      .subscribe({
        next: () => {
          this.destinationCreated.emit(); // Notifica a criação
          this.close.emit(); // Fecha a modal
        },
        error: (err) => {
          this.isLoading.set(false);
          this.errorMessage.set('Falha ao sugerir destino. Verifique os dados.');
          console.error('Erro ao criar destino:', err);
        }
      });
  }

  onClose() {
    this.close.emit();
  }
}
