import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class GroupService {

  constructor(private http: HttpClient) { }

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
