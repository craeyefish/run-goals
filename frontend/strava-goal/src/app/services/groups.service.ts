import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, signal } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class GroupService {

  groups = signal<Group[]>([]);
  groupCreated = signal<number | null>(null);
  selectedGroup = signal<Group | null>(null);
  goals = signal<Goal[]>([]);
  goalCreated = signal<number | null>(null);
  selectedGoal = signal<Goal | null>(null);

  constructor(private http: HttpClient) { }

  createGroup(request: CreateGroupRequest): Observable<CreateGroupResponse> {
    return this.http.post<CreateGroupResponse>('/api/groups', request);
  }

  getGroups(userID: number): Observable<GetGroupsResponse> {
    const params = new HttpParams().set('userID', userID)
    return this.http.get<GetGroupsResponse>('/api/groups', { params })
  }

  createGoal(request: CreateGoalRequest): Observable<CreateGoalResponse> {
    return this.http.post<CreateGoalResponse>('/api/group-goal', request);
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

  loadGroups() {
    const userID = 1;
    this.getGroups(userID).subscribe({
      next: (response) => {
        this.groups.set(response.groups);
        const createdGroup = this.groups().find(group => group.id === this.groupCreated());
        if (createdGroup) {
          this.selectedGroup.set(createdGroup);
          this.resetGroupCreated();
        } else if (response.groups.length > 0) {
          this.selectedGroup.set(response.groups[0]);
        }
      },
      error: (err) => {
        console.error('Failed to load groups', err)
      }
    })
  }

  loadGoals(groupID: number) {
    this.getGroupGoals(groupID).subscribe({
      next: (response) => {
        this.goals.set(response.goals);
        const createdGoal = this.goals().find(goal => goal.id === this.goalCreated());
        if (createdGoal) {
          this.selectedGoal.set(createdGoal);
          this.resetGoalCreated();
        } else if (response.goals.length > 0) {
          this.selectedGoal.set(response.goals[0]);
        }
      },
      error: (err) => {
        console.error('Failed to load group goals', err)
      }
    })
  }

  notifyGoalCreated(goalID: number) {
    this.goalCreated.set(goalID);
  }
  resetGoalCreated() {
    this.goalCreated.set(null);
  }

  notifyGroupCreated(groupID: number) {
    this.groupCreated.set(groupID);
  }
  resetGroupCreated() {
    this.groupCreated.set(null);
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

export interface CreateGroupResponse {
  group_id: number;
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

export interface CreateGoalResponse {
  goal_id: number;
}

export interface GetGroupGoalsResponse {
  goals: Goal[];
}

export interface GetGroupGoalProgressResponse {
  members: Member[];
}
