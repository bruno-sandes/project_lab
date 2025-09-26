import { Component, inject, signal } from '@angular/core'; // Importamos 'signal'
import { CommonModule } from '@angular/common';
import { FormGroup, FormControl, Validators, ReactiveFormsModule } from '@angular/forms';
import { AuthService } from '../../../login/service/login-service'; 
import { RegisterRequest } from '../../model/ register.model';
import { NavigateService } from '../../../../shared/services/navigate-service'; 


@Component({
  selector: 'app-register',
  standalone: true,
  imports: [
    CommonModule, 
    ReactiveFormsModule, 
  ],
  templateUrl: './register.html',
  styleUrls: ['./register.css'], 
})
export class RegisterComponent {
  private authService = inject(AuthService);
  private navigateService = inject(NavigateService);

  public isLoading = signal(false);
  public errorMessage = signal<string | null>(null);
  public successMessage = signal<string | null>(null);

  registerForm = new FormGroup({
    name: new FormControl('', [Validators.required]),
    email: new FormControl('', [Validators.required, Validators.email]),
    password: new FormControl('', [Validators.required, Validators.minLength(6)]),
  });
  
  get name() { return this.registerForm.get('name'); }
  get email() { return this.registerForm.get('email'); }
  get password() { return this.registerForm.get('password'); }

  navigateToLogin(){
    this.navigateService.toLogin();
  }

  onSubmit(): void {
    this.errorMessage.set(null);
    this.successMessage.set(null);

    if (this.registerForm.invalid) {
      this.registerForm.markAllAsTouched();
      return;
    }

    this.isLoading.set(true); 
    const data = this.registerForm.value as RegisterRequest;

    this.authService.register(data).subscribe({
      next: () => {
        this.successMessage.set('Conta criada com sucesso! Redirecionando para o login...');
        this.registerForm.reset();
        
        setTimeout(() => {
            this.navigateService.toLogin()
        }, 2000);
      },
      error: (err: Error) => {
        this.errorMessage.set(err.message); 
        this.isLoading.set(false);
      },
      complete: () => {
        this.isLoading.set(false); 
      }
    });
  }
}