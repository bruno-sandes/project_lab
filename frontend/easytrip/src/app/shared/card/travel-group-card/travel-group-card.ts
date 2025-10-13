import { Component, computed, input } from '@angular/core';
import { TravelGroupListItem } from '../../../features/groups-dashboard/models/travel_groups';
import { DatePipe } from '@angular/common';

export type CardLayout = 'upcoming' | 'past'

@Component({
  selector: 'app-travel-group-card',
  imports: [DatePipe],
  templateUrl: './travel-group-card.html',
  styleUrl: './travel-group-card.css'
})
export class TravelGroupCard {
  group = input.required<TravelGroupListItem>();
  layout = input.required<CardLayout>();

  daysDifference = computed(() => {
    const groupData = this.group();
    if (!groupData) return 0;

    const today = new Date();

    const dateToCheck = new Date(groupData.start_date); 
    const diffTime = Math.abs(dateToCheck.getTime() - today.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24)); 
    
    return diffDays;
  });

  displayedMembers = computed(() => {
      const memberCount = this.group().memberCount;
      return Array(Math.min(4, memberCount)).fill(0);
  });

  remainingMembersCount = computed(() => {
      const memberCount = this.group().memberCount;
      return memberCount > 4 ? memberCount - 4 : 0;
  });
}
