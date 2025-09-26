import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError, BehaviorSubject } from 'rxjs';
import { catchError } from 'rxjs/operators'; // tap foi removido!
import { LoginRequest, LoginResponse, RegisterRequest } from '../model/login.model';
import { environment } from '../../../../environments/environment';

const TOKEN_KEY = 'auth_token';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private http = inject(HttpClient); 
  private apiUrl = `${environment.apiUrl}/auth`;
  
  private isAuthenticatedSubject = new BehaviorSubject<boolean>(this.hasToken());
  public isAuthenticated$ = this.isAuthenticatedSubject.asObservable();

  private hasToken(): boolean {
    return !!localStorage.getItem(TOKEN_KEY);
  }

  public setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token);
    this.isAuthenticatedSubject.next(true);
  }

  /**
   * GETTER: Obtém o token JWT armazenado.
   */
  public getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY);
  }

  login(credentials: LoginRequest): Observable<LoginResponse> {
    const url = `${this.apiUrl}/login`;
    
    return this.http.post<LoginResponse>(url, credentials).pipe(
      catchError((error) => this.handleError(error))
    );
  }

  register(data: RegisterRequest): Observable<void> {
    const url = `${this.apiUrl}/register`;
    
    return this.http.post<void>(url, data).pipe(
      catchError((error) => this.handleError(error))
    );
  }

 
  logout(): void {
    localStorage.removeItem(TOKEN_KEY);
    this.isAuthenticatedSubject.next(false);
  }


  private handleError(error: HttpErrorResponse) {
    let errorMessage = 'Ocorreu um erro desconhecido.';
    
    if (error.status === 401) {
      errorMessage = 'Credenciais inválidas. Verifique seu email e senha.';
    } else if (error.status === 400) {
      errorMessage = 'Dados de entrada inválidos.';
    } else if (error.error && error.error.message) {
      errorMessage = error.error.message;
    } else {
      errorMessage = `Erro ${error.status}: Tente novamente mais tarde.`;
    }
    
    return throwError(() => new Error(errorMessage));
  }
}