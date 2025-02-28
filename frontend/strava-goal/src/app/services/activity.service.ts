import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Activity {
  id: number;
  strava_activity_id: number;
  user_id: number;
  name: string;
  distance: number;
  start_date: string; // or Date if you parse it
  map_polyline: string;
  has_summit: boolean; // true if at least one peak was bagged on this activity
}

@Injectable({
  providedIn: 'root',
})
export class ActivityService {
  constructor(private http: HttpClient) {}

  // Adjust userId param as needed
  getActivitiesForUser(userId: number): Observable<Activity[]> {
    return this.http.get<Activity[]>(`/api/activities?userId=${userId}`);
  }
}
