import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, signal } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class GroupService {

  selectedGroup = signal<Group | null>(null);
  selectedGoal = signal<Goal | null>(null);
  goalCreated = signal<boolean>(false);
  goals = signal<Goal[]>([]);

  constructor(private http: HttpClient) { }

  notifyGoalCreated() {
    this.goalCreated.set(true);
  }
  resetGoalCreated() {
    this.goalCreated.set(false);
  }

  createGroup(request: CreateGroupRequest): Observable<any> {
    return this.http.post('/api/groups', request);
  }

  getGroups(userID: number): Observable<GetGroupsResponse> {
    const params = new HttpParams().set('userID', userID)
    return this.http.get<GetGroupsResponse>('/api/groups', { params })
  }

  createGoal(request: CreateGoalRequest): Observable<any> {
    return this.http.post('/api/group-goal', request);
  }

  getGroupGoals(groupID: number): Observable<GetGroupGoalsResponse> {
    const params = new HttpParams().set('groupID', groupID)
    return this.http.get<GetGroupGoalsResponse>('/api/group-goals', { params })
  }

  getGoalProgress(groupID: number, goalID: number): Observable<GetGroupGoalProgressResponse> {
    const params = new HttpParams()
      .set('goalID', goalID)
      .set('groupID', groupID)
    // this.http.get<GetGroupGoalProgressResponse>('/api/group-goal-progress', { params })
    return this.http.get<GetGroupGoalProgressResponse>('/api/group-goals', { params })
  }

  loadGoals(groupID: number) {
    this.getGroupGoals(groupID).subscribe({
      next: (response) => {
        this.goals.set(response.goals);
        if (!this.selectedGoal() && response.goals.length > 0) {
          this.selectedGoal.set(response.goals[0]);
        }
      },
      error: (err) => {
        console.error('Failed to load group goals', err)
      }
    })
  }

}

export interface Group {
  id: number;
  name: string;
  created_by: number;
  created_at: string;
}

export interface Goal {
  id: number;
  group_id: number;
  name: string;
  target_value: string;
  start_date: string;
  end_date: string;
  create_at: string;
}

export interface Member {
  id: number;
  name: string;
  activities: number;
  summits: number;
  distance: number;
  contribution: number;
}

export interface CreateGroupRequest {
  name: string;
  created_by: number;
}

export interface GetGroupsResponse {
  groups: Group[];
}

export interface CreateGoalRequest {
  group_id: number;
  name: string;
  target_value: string;
  start_date: string;
  end_date: string;
}

export interface GetGroupGoalsResponse {
  goals: Goal[];
}

export interface GetGroupGoalProgressResponse {
  members: Member[];
}
