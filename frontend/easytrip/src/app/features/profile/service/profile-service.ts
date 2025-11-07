import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { 
  UserProfileResponse, 
  UserProfileUpdateRequest 
} from '../models/profile.model'; 

@Injectable({
  providedIn: 'root'
})
export class ProfileService {
  private apiUrl = 'http://localhost:8080/profile'; 
  private http = inject(HttpClient);

  
  getProfile(): Observable<UserProfileResponse> {
    return this.http.get<UserProfileResponse>(this.apiUrl);
  }

  updateProfile(updateData: UserProfileUpdateRequest): Observable<UserProfileResponse> {
    return this.http.patch<UserProfileResponse>(this.apiUrl, updateData);
  }
}