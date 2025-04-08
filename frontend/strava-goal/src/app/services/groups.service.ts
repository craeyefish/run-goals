import { HttpClient, HttpParams } from '@angular/common/http';
import { StringToken } from '@angular/compiler';
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

}

export interface Group {
  id: number;
  name: string;
  created_by: number;
  created_at: string;
}

export interface CreateGroupRequest {
  name: string;
  created_by: number;
}

export interface GetGroupsResponse {
  groups: Group[];
}
