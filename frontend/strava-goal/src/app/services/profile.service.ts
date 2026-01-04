import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface UserProfile {
  id: number;
  strava_athelete_id: number;
  username?: string;
  is_admin: boolean;
  last_distance: number;
  last_updated: string;
  created_at: string;
  updated_at: string;
}

@Injectable({
  providedIn: 'root',
})
export class ProfileService {
  constructor(private http: HttpClient) {}

  getUserProfile(): Observable<UserProfile> {
    return this.http.get<UserProfile>('/api/profile');
  }

  updateUsername(username: string): Observable<UserProfile> {
    return this.http.put<UserProfile>('/api/profile', { username });
  }
}
