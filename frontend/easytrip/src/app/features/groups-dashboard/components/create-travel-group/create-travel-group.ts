import { Component, inject, signal, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormControl, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { Subscription } from 'rxjs'; 
import { TravelGroupsService } from '../../services/travel-groups-service'; 
import { TravelGroupCreateRequest } from '../../models/travel_groups'; 
import { NavigateService } from '../../../../shared/services/navigate-service';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'app-create-travel-group',
  imports: [CommonModule, ReactiveFormsModule], 
  standalone: true,
  templateUrl: './create-travel-group.html',
  styleUrl: './create-travel-group.css'
})
export class CreateTravelGroup {
 
  private travelGroupService = inject(TravelGroupsService);
  private navigationService = inject(NavigateService)
  public isLoading = signal<boolean>(false);
  public errorMessage = signal<string | null>(null);
  public successMessage = signal<string | null>(null);

  
  tripForm: FormGroup = new FormGroup({
    tripName: new FormControl('', [Validators.required, Validators.minLength(3)]),
    description: new FormControl(''),
    startDate: new FormControl('', [Validators.required]),
    endDate: new FormControl('', [Validators.required])
  });

  
  onSubmit(): void {
    this.errorMessage.set(null); 

    if (this.tripForm.invalid) {
      this.tripForm.markAllAsTouched();
      return;
    }

    this.isLoading.set(true); 

    const formData = this.tripForm.value;

    const payload: TravelGroupCreateRequest = {
      name: formData.tripName,
      description: formData.description,
      start_date: formData.startDate,
      end_date: formData.endDate,
    };

    this.travelGroupService.createGroup(payload).subscribe({
      next: () => {
        this.successMessage.set('Viagem criada com sucesso!');
        this.tripForm.reset();
        
       
      },
      error: (err: HttpErrorResponse) => {
        const message = err.error?.message || err.message || 'Erro desconhecido ao criar viagem.';
        this.errorMessage.set(message);
        this.isLoading.set(false); 
      },
      complete: () => {
        this.isLoading.set(false);
      }
    });    
  }

  onCancel(): void {
    this.navigationService.toGroupsDashboard(); 
  }
}