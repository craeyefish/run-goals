import { HttpClient } from '@angular/common/http';
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

}

export interface CreateGroupRequest {
  name: string;
  created_by: number;
}
