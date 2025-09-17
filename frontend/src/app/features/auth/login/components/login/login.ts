import { Component } from '@angular/core';
import { Router } from '@angular/router'; // Importar Router para navegação
import { AuthService } from '../../services/auth';
import { LoginRequest } from '../../model/auth.model';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-login',
  templateUrl: './login.html',
  imports:[CommonModule]
})
export class  LoginComponent {
  // Propriedades para vincular aos campos do formulário
  email!: string;
  password!: string;
  
  // Propriedade para controlar a exibição de mensagens de erro
  errorMessage: string | null = null;
  loading = false;

  constructor(
    private authService: AuthService, // Injeta o serviço AuthService
    private router: Router // Injeta o serviço Router para navegação
  ) { }

  /**
   * Este método é chamado quando o formulário de login é enviado.
   * Ele utiliza o AuthService para se comunicar com a API.
   */
  onSubmit(): void {
    // Reseta a mensagem de erro a cada nova tentativa
    this.errorMessage = null;
    this.loading = true;

    // Cria um objeto com os dados de login no formato esperado pela API
    const credentials: LoginRequest = {
      email: this.email,
      password: this.password
    };

    // Chama o método login() do serviço e se inscreve no Observable
    this.authService.login(credentials).subscribe({
      next: (response) => {
        // Se a requisição for bem-sucedida (status 200)
        console.log('Login realizado com sucesso!', response);
        this.loading = false;
        
        // Redireciona o usuário para a página de grupos, por exemplo
        this.router.navigate(['/grupos']);
      },
      error: (err) => {
        // Se a requisição falhar (ex: status 401, 400, etc.)
        this.loading = false;
        // O `AuthService` já tratou o erro e retornou uma mensagem amigável
        this.errorMessage = err.message || 'Ocorreu um erro no login.';
        console.error('Falha no login:', err);
      }
    });
  }
}