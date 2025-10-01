import { Component, inject, signal } from '@angular/core';
import { NavigateService } from '../services/navigate-service';
import { AuthService } from '../../features/login/service/login-service';

@Component({
  selector: 'app-header',
  imports: [],
  templateUrl: './header.html',
  styleUrl: './header.css'
})
export class Header {
 private navigateService = inject(NavigateService);
 private authService = inject(AuthService);

  isDropdownOpen = signal(false);
  userProfileImageUrl = signal('https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_960_720.png'); // Imagem de perfil padrão


  toggleDropdown() {
    this.isDropdownOpen.update(val => !val);
  }

  logout() {
    this.authService.logout();
    this.navigateService.toLogin(); 
  }


  navigateToExplorar() {
    console.log('Navegar para Explorar');
    // this.navigateService.toExplore();
  }

  navigateToHome() {
    console.log('Navegar para inicio');
    this.navigateService.toGroupsDashboard();
  }
}

