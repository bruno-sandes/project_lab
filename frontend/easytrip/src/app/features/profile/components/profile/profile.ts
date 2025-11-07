import { Component, OnInit, signal, inject } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { catchError, finalize, of } from 'rxjs'; 
import { ProfileService } from '../../service/profile-service';
import { UserProfileResponse, UserProfileUpdateRequest } from '../../models/profile.model';

@Component({
  selector: 'app-profile',
  imports: [ReactiveFormsModule, CommonModule], 
  templateUrl: './profile.html',
  styleUrl: './profile.css'
})
export class Profile implements OnInit {
  
  // Dependências
  private profileService = inject(ProfileService);
  
  // Estado
  userEmail = signal<string>('');
  isLoading = signal<boolean>(false);
  successMessage = signal<string | null>(null);
  errorMessage = signal<string | null>(null);
  
  profileForm = new FormGroup({
    name: new FormControl('', { 
      nonNullable: true, 
      validators: [Validators.required, Validators.minLength(3)] 
    }),
  });

  // Getter de conveniência
  get nameControl(): FormControl {
    return this.profileForm.get('name') as FormControl;
  }

  ngOnInit(): void {
    this.loadProfile();
  }

  /**
   * Carrega os dados do perfil da API.
   */
  loadProfile() {
    this.isLoading.set(true);
    this.errorMessage.set(null);
    
    this.profileService.getProfile().pipe(
      finalize(() => this.isLoading.set(false)),
      catchError(err => {
        this.errorMessage.set('Falha ao carregar o perfil. Tente novamente.');
        console.error('Erro ao buscar perfil:', err);
        return of(null);
      })
    ).subscribe((profile: UserProfileResponse | null) => {
      if (profile) {
        this.userEmail.set(profile.email);
        
        this.profileForm.controls.name.setValue(profile.name);
        
        this.profileForm.markAsPristine();
      }
    });
  }

  /**
   * Envia a requisição de atualização de nome.
   */
  onSubmit() {
    if (this.profileForm.invalid || this.profileForm.pristine) {
      return;
    }
    
    this.isLoading.set(true);
    this.successMessage.set(null);
    this.errorMessage.set(null);

    const updateData: UserProfileUpdateRequest = {
      name: this.profileForm.controls.name.value
    };

    this.profileService.updateProfile(updateData).pipe(
      finalize(() => this.isLoading.set(false)),
      catchError(err => {
        this.errorMessage.set('Falha ao salvar. Verifique se o nome é válido.');
        console.error('Erro ao atualizar perfil:', err);
        return of(null);
      })
    ).subscribe((updatedProfile: UserProfileResponse | null) => {
      if (updatedProfile) {
        this.profileForm.controls.name.setValue(updatedProfile.name); 
        this.profileForm.markAsPristine(); 
        this.successMessage.set('Perfil atualizado com sucesso!');
        setTimeout(() => this.successMessage.set(null), 3000); 
      }
    });
  }
}