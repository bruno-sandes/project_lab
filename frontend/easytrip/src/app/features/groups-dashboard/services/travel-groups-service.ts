import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { TravelGroupCreateRequest, TravelGroupDetails, TravelGroupListItem } from '../models/travel_groups';

@Injectable({
  providedIn: 'root'
})
export class TravelGroupsService {
  private apiUrl = 'http://localhost:8080/groups'; 

  http = inject(HttpClient)


  /**
   * GET /groups
   * Lista todos os grupos aos quais o usuário logado pertence ou criou.
   */
  listGroups(): Observable<TravelGroupListItem[]> {
    return this.http.get<TravelGroupListItem[]>(this.apiUrl);
  }

  /**
   * POST /groups
   * Cria um novo grupo de viagem.
   */
  createGroup(groupData: TravelGroupCreateRequest): Observable<TravelGroupDetails> {
    return this.http.post<TravelGroupDetails>(this.apiUrl, groupData);
  }
}
