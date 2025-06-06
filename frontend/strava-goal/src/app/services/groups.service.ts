import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, signal } from '@angular/core';
import { Observable } from 'rxjs';
import { DatePipe } from '@angular/common';

@Injectable({
  providedIn: 'root',
})
export class GroupService {

  groups = signal<Group[]>([]);
  groupUpdate = signal<number | null>(null);
  selectedGroup = signal<Group | null>(null);
  goals = signal<Goal[]>([]);
  goalCreated = signal<number | null>(null);
  selectedGoal = signal<Goal | null>(null);
  selectedGoalChange = signal<boolean>(false);
  membersContribution = signal<MemberContribution[]>([]);
  memberAddedOrRemoved = signal<boolean>(false);

  constructor(private http: HttpClient, private datePipe: DatePipe) { }

  createGroup(request: CreateGroupRequest): Observable<CreateGroupResponse> {
    console.log('create group request: ', request);
    return this.http.post<CreateGroupResponse>('/api/groups', request);
  }

  updateGroup(request: UpdateGroupRequest): Observable<any> {
    return this.http.put<any>('/api/groups', request);
  }

  getGroups(userID: number): Observable<GetGroupsResponse> {
    const params = new HttpParams().set('userID', userID)
    return this.http.get<GetGroupsResponse>('/api/groups', { params })
  }

  createGoal(request: CreateGoalRequest): Observable<CreateGoalResponse> {
    return this.http.post<CreateGoalResponse>('/api/group-goal', request);
  }

  updateGoal(request: UpdateGoalRequest): Observable<any> {
    return this.http.put<any>('/api/group-goal', request);
  }

  getGroupGoals(groupID: number): Observable<GetGroupGoalsResponse> {
    const params = new HttpParams().set('groupID', groupID);
    return this.http.get<GetGroupGoalsResponse>('/api/group-goals', { params })
  }

  getGroupMembers(groupID: number): Observable<GetGroupMembersResponse> {
    const params = new HttpParams().set('groupID', groupID);
    return this.http.get<GetGroupMembersResponse>('/api/group-members', { params })
  }

  createGroupMember(request: CreateGroupMemberRequest): Observable<any> {
    return this.http.post<any>('/api/group-member', request);
  }

  leaveGroup(request: LeaveGroupRequest): Observable<any> {
    const params = new HttpParams()
      .set('groupID', request.groupID)
      .set('userID', request.userID);
    return this.http.delete<any>('/api/group-member', { params });
  }

  getGroupMembersGoalContribution(groupID: number, startDate: string, endDate: string): Observable<GetGroupMembersGoalContributionResponse> {
    const formattedStartDate = this.datePipe.transform(startDate, 'yyyy-MM-dd');
    const formattedEndDate = this.datePipe.transform(endDate, 'yyyy-MM-dd');
    if (formattedStartDate && formattedEndDate) {
      const params = new HttpParams()
        .set('groupID', groupID)
        .set('startDate', formattedStartDate)
        .set('endDate', formattedEndDate);
      return this.http.get<GetGroupMembersGoalContributionResponse>('/api/group-members-contribution', { params })
    } else {
      console.log("Failed to format date");
      throw new Error('Failed to format date');
    }
  }

  loadGroups() {
    const userID = 1;
    this.getGroups(userID).subscribe({
      next: (response) => {
        this.groups.set(response.groups);
      },
      error: (err) => {
        console.error('Failed to load groups', err)
      }
    })
  }

  loadGoals(groupID: number) {
    this.getGroupGoals(groupID).subscribe({
      next: (response) => {
        const goals = response.goals.map((goal) => ({
          ...goal,
          start_date: this.datePipe.transform(goal.start_date, 'yyyy-MM-dd'),
          end_date: this.datePipe.transform(goal.end_date, 'yyyy-MM-dd'),
          created_at: this.datePipe.transform(goal.created_at, 'yyyy-MM-dd'),
        }))
        this.goals.set(goals);
        const createdGoal = this.goals().find(goal => goal.id === this.goalCreated());
        if (createdGoal) {
          this.selectedGoal.set(createdGoal);
          this.resetGoalCreated();
        } else if (response.goals.length > 0) {
          this.selectedGoal.set(response.goals[0]);
        } else {
          this.selectedGoal.set(null);
        }
        this.notifySelectedGoalChange();
      },
      error: (err) => {
        console.error('Failed to load group goals', err)
      }
    })
  }

  loadMembersGoalContribution(groupID: number, startDate: string, endDate: string) {
    this.getGroupMembersGoalContribution(groupID, startDate, endDate).subscribe({
      next: (response) => {
        this.membersContribution.set(response.members);
      },
      error: (err) => {
        console.error('Failed to load members goal contribution', err)
      }
    })
  }

  notifyGoalCreated(goalID: number) {
    this.goalCreated.set(goalID);
  }
  resetGoalCreated() {
    this.goalCreated.set(null);
  }

  notifySelectedGoalChange() {
    this.selectedGoalChange.set(true);
  }
  resetSelectedGoalChange() {
    this.selectedGoalChange.set(false);
  }

  notifyGroupUpdate() {
    this.loadGroups();
  }
  resetGroupUpdate() {
    this.groupUpdate.set(null);
  }

  notifyMemberAddedOrRemoved(groupMemberID: number) {
    this.memberAddedOrRemoved.set(true)
  }
  resetMemberAddedOrRemoved() {
    this.memberAddedOrRemoved.set(false)
  }
}

export interface Group {
  id: number;
  name: string;
  code: string;
  created_by: number;
  created_at: string;
}

export interface Goal {
  id: number;
  group_id: number;
  name: string;
  target_value: string;
  start_date: string | null;
  end_date: string | null;
  created_at: string | null;
}

export interface Member {
  id: number;
  group_id: number;
  user_id: number;
  role: string;
  joined_at: string;
}

export interface MemberContribution {
  group_member_id: number;
  group_id: number;
  user_id: number;
  role: string;
  joined_at: string;
  total_activities: number | null;
  total_distance: number | null;
  total_unique_summits: number | null;
  total_summits: number | null;
}

export interface CreateGroupRequest {
  name: string;
  created_by: number;
}

export interface CreateGroupResponse {
  group_id: number;
}

export interface UpdateGroupRequest {
  id: number;
  name: string;
  created_by: number;
  created_at: string;
}

export interface GetGroupsResponse {
  groups: Group[];
}

export interface CreateGoalRequest {
  group_id: number;
  name: string;
  target_value: string;
  start_date: Date;
  end_date: Date;
}

export interface CreateGoalResponse {
  goal_id: number;
}

export interface UpdateGoalRequest {
  id: number;
  group_id: number;
  name: string;
  target_value: string;
  start_date: Date;
  end_date: Date;
  created_at: Date;
}

export interface GetGroupGoalsResponse {
  goals: Goal[];
}

export interface GetGroupMembersGoalContributionResponse {
  members: MemberContribution[];
}

export interface GetGroupMembersResponse {
  members: Member[];
}

export interface CreateGroupMemberRequest {
  group_code: string;
  user_id: number;
  role: string;
}

export interface LeaveGroupRequest {
  userID: number;
  groupID: number;
}
