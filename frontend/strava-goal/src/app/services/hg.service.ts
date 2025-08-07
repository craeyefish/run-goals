import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject } from 'rxjs';

export interface Activity {
  id: number;
  strava_activity_id: number;
  user_id: number;
  name: string;
  description: string;
  distance: number;
  start_date: string; // or Date if you parse it
  map_polyline: string;
  has_summit: boolean; // true if at least one peak was bagged on this activity
  photo_url?: string;
}

@Injectable({ providedIn: 'root' })
export class HgService {
  private activitiesSubject = new BehaviorSubject<Activity[] | null>(null);
  activities$ = this.activitiesSubject.asObservable();

  private loading = false;

  constructor(private http: HttpClient) {}

  loadActivities(forceRefresh: boolean = false): void {
    if (this.activitiesSubject.value && !forceRefresh) {
      return;
    }

    if (this.loading) {
      return;
    }

    this.loading = true;

    this.http.get<Activity[]>('/hikegang/activities').subscribe({
      next: (acts) => {
        this.activitiesSubject.next(acts);
        this.loading = false;
      },
      error: (err) => {
        console.error('Failed to load activities', err);
        this.loading = false;
      },
    });
  }

  refreshActivities(): void {
    this.loadActivities(true);
  }

  triggerSync() {
    return this.http.post<any>('/hikegang/sync', {});
  }
}
