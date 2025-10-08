// src/app/core/interceptors/auth.interceptor.ts

import { HttpInterceptorFn, HttpRequest, HttpHandlerFn, HttpEvent} from '@angular/common/http';
import { Observable } from 'rxjs';
import { Login } from '../features/login/component/login/login';
import { inject } from '@angular/core';
import { AuthService } from '../features/login/service/login-service';

export const authInterceptor: HttpInterceptorFn = (req: HttpRequest<unknown>, next: HttpHandlerFn): Observable<HttpEvent<unknown>> => {
  
  const tokenService = inject(AuthService); 
  const authToken = tokenService.getToken();
   //TO DO 
  //Define a URL do seu backend para evitar enviar token para APIs externas, mudar depois para uma variavel de ambiente
  const isApiUrl = req.url.includes('localhost:8080'); 

  if (authToken && isApiUrl) {
    //Clona a requisição e adiciona o cabeçalho Authorization
    const authRequest = req.clone({
      setHeaders: {
        Authorization: `Bearer ${authToken}`
      }
    });

    return next(authRequest);
  }
  
  return next(req);
};