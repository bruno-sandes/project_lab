import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { LoginRequest,LoginResponse,RegisterRequest } from '../model/auth.model';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  // A URL base da API, conforme o YAML (você pode usar um ambiente para isso)
  private apiUrl = 'http://localhost:8080/auth';

  constructor(private http: HttpClient) { }

  /**
   * Realiza a requisição de login para a API.
   * @param credentials Os dados de email e senha do usuário.
   * @returns Um Observable com a resposta do login (contendo o token).
   */
  login(credentials: LoginRequest): Observable<LoginResponse> {
    const url = `${this.apiUrl}/login`;
    return this.http.post<LoginResponse>(url, credentials).pipe(
      // Você pode usar 'tap' para salvar o token no localStorage
      tap(response => {
        if (response.token) {
          localStorage.setItem('authToken', response.token);
          console.log('Login bem-sucedido. Token salvo.');
        }
      }),
      // 'catchError' para tratar erros, como credenciais inválidas (401)
      catchError(this.handleError)
    );
  }

  /**
   * Realiza a requisição de cadastro de novo usuário para a API.
   * @param userData Os dados de nome, email e senha do novo usuário.
   * @returns Um Observable sem um corpo de resposta específico (201 Created).
   */
  register(userData: RegisterRequest): Observable<any> {
    const url = `${this.apiUrl}/register`;
    return this.http.post(url, userData).pipe(
      // 'catchError' para tratar erros, como email já cadastrado (400)
      catchError(this.handleError)
    );
  }

  /**
   * Remove o token de autenticação e realiza o logout local.
   */
  logout(): void {
    localStorage.removeItem('authToken');
    console.log('Logout realizado. Token removido.');
  }

  /**
   * Checa se o usuário está autenticado verificando a existência do token.
   * @returns true se o token existir, false caso contrário.
   */
  isLoggedIn(): boolean {
    return !!localStorage.getItem('authToken');
  }

  /**
   * Método privado para lidar com erros de requisição HTTP.
   * @param error O erro de resposta HTTP.
   * @returns Um Observable com o erro.
   */
  private handleError(error: HttpErrorResponse): Observable<never> {
    let errorMessage = 'Ocorreu um erro desconhecido.';
    if (error.error instanceof ErrorEvent) {
      // Erro do lado do cliente ou de rede
      errorMessage = `Erro: ${error.error.message}`;
    } else {
      // Erro retornado pelo backend
      console.error(
        `Código do erro: ${error.status}, ` +
        `Corpo do erro: ${JSON.stringify(error.error)}`);

      if (error.status === 401) {
        errorMessage = 'Credenciais inválidas. Por favor, verifique seu email e senha.';
      } else if (error.status === 400) {
        errorMessage = 'Erro de validação. Verifique os dados enviados.';
      } else {
        errorMessage = `Erro no servidor: ${error.status}`;
      }
    }
    // Retorna um Observable com uma mensagem de erro amigável para o componente
    return throwError(() => new Error(errorMessage));
  }
}