import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface UserContribution {
  name: string;
  totalDistance: number;
}

export interface GoalProgress {
  goal: number;
  currentProgress: number;
  contributions: UserContribution[];
}

@Injectable({
  providedIn: 'root',
})
export class ProgressService {
  constructor(private http: HttpClient) {}

  getProgress(): Observable<GoalProgress> {
    return this.http.get<GoalProgress>('/api/progress');
  }
}
