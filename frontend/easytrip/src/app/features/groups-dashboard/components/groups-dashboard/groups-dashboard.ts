import { Component, computed, inject } from '@angular/core';
import { TravelGroupsService } from '../../services/travel-groups-service';
import { Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { TravelGroupListItem } from '../../models/travel_groups';
import { TravelGroupCard } from "../../../../shared/card/travel-group-card/travel-group-card";

@Component({
  selector: 'app-groups-dashboard',
  imports: [TravelGroupCard],
  templateUrl: './groups-dashboard.html',
  styleUrl: './groups-dashboard.css'
})
export class GroupsDashboard {
 
  private groupService = inject(TravelGroupsService);
  private router = inject(Router);
  
  allGroups = toSignal(this.groupService.listGroups(), { 
    initialValue: [] as TravelGroupListItem[] 
  });

  upcomingGroups = computed(() => {
    const groups = this.allGroups(); 
    const today = new Date().setHours(0, 0, 0, 0); 
    
    return groups
      .filter(group => new Date(group.endDate).getTime() >= today)
      .sort((a, b) => new Date(a.startDate).getTime() - new Date(b.startDate).getTime()); 
  });

  pastGroups = computed(() => {
    const groups = this.allGroups();
    const today = new Date().setHours(0, 0, 0, 0); 
    
    return groups
      .filter(group => new Date(group.endDate).getTime() < today)
      .sort((a, b) => new Date(b.startDate).getTime() - new Date(a.startDate).getTime()); 
  });

  onNavigateToCreateGroup() {
    this.router.navigate(['/groups/create']); 
  }

  onJoinGroup() {
    console.log('Abrir modal para entrar em grupo');
  }
}
