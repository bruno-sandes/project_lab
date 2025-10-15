import { inject, Injectable } from '@angular/core';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class NavigateService {
  private router = inject(Router);

  public toRegister(): void {
    this.router.navigate(['register']);
  }

  public toLogin(): void {
    this.router.navigate(['']);
  }

  public toGroupsDashboard(): void {
    this.router.navigate(['/inicio/groups-dashboard']);
  }

  public toCreateTravelGroup(): void {
    this.router.navigate(['/inicio/create-group']);
  }

  public toProfile(): void {
    this.router.navigate(['/inicio/profile']);

  }

  //metodo generico que constroi um path
  public toRoute(path: string, params?: any): void {
    this.router.navigate([path], { queryParams: params });
  }
}
