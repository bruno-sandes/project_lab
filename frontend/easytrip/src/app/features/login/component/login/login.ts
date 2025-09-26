import { Component, inject, signal } from '@angular/core';
import { AuthService } from '../../service/login-service';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { LoginRequest } from '../../model/login.model';
import { NavigateService } from '../../../../shared/services/navigate-service';

@Component({
  selector: 'app-login',
  imports: [ReactiveFormsModule],
  standalone: true,
  templateUrl: './login.html',
  styleUrl: './login.css'
})
export class Login {
 private authService = inject(AuthService);
 private navigateService = inject(NavigateService)

  public isLoading = signal(false);
  public errorMessage = signal<string | null>(null);

  loginForm = new FormGroup({
    email: new FormControl('', [Validators.required, Validators.email]),
    password: new FormControl('', [Validators.required, Validators.minLength(6)]),
  });


  get email() {
    return this.loginForm.get('email');
  }
  
  get password() {
    return this.loginForm.get('password');
  }
  
  navigateToRegister(){
    this.navigateService.toRegister();
  }

  onSubmit(): void {
    this.errorMessage.set(null); 

    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      return;
    }

    this.isLoading.set(true); 

    const credentials = this.loginForm.value as LoginRequest;

    this.authService.login(credentials).subscribe({
      next: (response) => {
        this.authService.setToken(response.token);
        this.navigateService.toGroupsDashboard();         
        this.isLoading.set(false);
      },
      error: (err: Error) => {
        this.errorMessage.set(err.message); 
        this.isLoading.set(false); 
      },
    });
  }

}