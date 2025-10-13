import { Component, computed, inject } from '@angular/core';
import { TravelGroupsService } from '../../services/travel-groups-service';
import { Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { TravelGroupListItem } from '../../models/travel_groups';
import { TravelGroupCard } from "../../../../shared/card/travel-group-card/travel-group-card";
import { NavigateService } from '../../../../shared/services/navigate-service';

@Component({
  selector: 'app-groups-dashboard',
  imports: [TravelGroupCard],
  templateUrl: './groups-dashboard.html',
  styleUrl: './groups-dashboard.css'
})
export class GroupsDashboard {
 
  private groupService = inject(TravelGroupsService);
  private navigateService = inject(NavigateService);

  allGroups = toSignal(this.groupService.listGroups(), { 
    initialValue: [] as TravelGroupListItem[] 
  });

  upcomingGroups = computed(() => {
    const groups = this.allGroups(); 
    const today = new Date().setHours(0, 0, 0, 0); 
    
    return groups
      .filter(group => new Date(group.end_date).getTime() >= today)
      .sort((a, b) => new Date(a.start_date).getTime() - new Date(b.start_date).getTime()); 
  });

  pastGroups = computed(() => {
    const groups = this.allGroups();
    const today = new Date().setHours(0, 0, 0, 0); 
    
    return groups
      .filter(group => new Date(group.end_date).getTime() < today)
      .sort((a, b) => new Date(b.start_date).getTime() - new Date(a.start_date).getTime()); 
  });

  onNavigateToCreateGroup() {
    this.navigateService.toCreateTravelGroup();
  }
//TODO LOGICA DE ENTRAR COM CODIGO ALEATORIO GERADO E VCALIDP POR APENAS 10 MINUTOS
  onJoinGroup() {
    console.log('Abrir modal para entrar em grupo');
  }
}
