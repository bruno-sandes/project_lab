import { Routes } from '@angular/router';
import { LoginComponent } from './features/auth/login/components/login/login';

export const routes: Routes = [
  // Define a rota raiz ('/') para carregar o LoginComponent
  { path: '', component: LoginComponent },
  
  // Rota de fallback para qualquer caminho não encontrado
  { path: '**', redirectTo: '' }
];
