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

  constructor(private http: HttpClient, private datePipe: DatePipe) {}

  createGroup(request: CreateGroupRequest): Observable<CreateGroupResponse> {
    console.log('create group request: ', request);
    return this.http.post<CreateGroupResponse>('/api/groups', request);
  }

  updateGroup(request: UpdateGroupRequest): Observable<any> {
    return this.http.put<any>('/api/groups', request);
  }

  getGroups(): Observable<GetGroupsResponse> {
    return this.http.get<GetGroupsResponse>('/api/groups');
  }

  createGoal(request: CreateGoalRequest): Observable<CreateGoalResponse> {
    return this.http.post<CreateGoalResponse>('/api/group-goal', request);
  }

  updateGoal(request: UpdateGoalRequest): Observable<any> {
    return this.http.put<any>('/api/group-goal', request);
  }

  getGroupGoals(groupID: number): Observable<GetGroupGoalsResponse> {
    const params = new HttpParams().set('groupID', groupID);
    return this.http.get<GetGroupGoalsResponse>('/api/group-goals', { params });
  }

  getGroupMembers(groupID: number): Observable<GetGroupMembersResponse> {
    const params = new HttpParams().set('groupID', groupID);
    return this.http.get<GetGroupMembersResponse>('/api/group-members', {
      params,
    });
  }

  createGroupMember(request: CreateGroupMemberRequest): Observable<any> {
    return this.http.post<any>('/api/group-member', request);
  }

  leaveGroup(request: LeaveGroupRequest): Observable<any> {
    const params = new HttpParams().set('groupID', request.groupID);
    return this.http.delete<any>('/api/group-member', { params });
  }

  getGroupMembersGoalContribution(
    groupID: number,
    startDate: string,
    endDate: string
  ): Observable<GetGroupMembersGoalContributionResponse> {
    const formattedStartDate = this.datePipe.transform(startDate, 'yyyy-MM-dd');
    const formattedEndDate = this.datePipe.transform(endDate, 'yyyy-MM-dd');
    if (formattedStartDate && formattedEndDate) {
      const params = new HttpParams()
        .set('groupID', groupID)
        .set('startDate', formattedStartDate)
        .set('endDate', formattedEndDate);
      return this.http.get<GetGroupMembersGoalContributionResponse>(
        '/api/group-members-contribution',
        { params }
      );
    } else {
      console.log('Failed to format date');
      throw new Error('Failed to format date');
    }
  }

  loadGroups() {
    this.getGroups().subscribe({
      next: (response) => {
        this.groups.set(response.groups);
      },
      error: (err) => {
        console.error('Failed to load groups', err);
      },
    });
  }

  loadGoals(groupID: number) {
    this.getGroupGoals(groupID).subscribe({
      next: (response) => {
        // Don't convert dates to Date objects - keep them as strings for consistency
        const goals = response.goals.map((goal) => ({
          ...goal,
          // Keep dates as strings since the rest of the app expects strings
          start_date: goal.start_date,
          end_date: goal.end_date,
          created_at: goal.created_at,
        }));

        this.goals.set(goals);

        // Find and maintain the currently selected goal after update
        const currentSelectedGoal = this.selectedGoal();
        if (currentSelectedGoal) {
          const updatedSelectedGoal = goals.find(
            (goal) => goal.id === currentSelectedGoal.id
          );
          if (updatedSelectedGoal) {
            this.selectedGoal.set(updatedSelectedGoal);
          }
        }

        // Handle goal creation selection
        const createdGoalId = this.goalCreated();
        if (createdGoalId) {
          const createdGoal = goals.find((goal) => goal.id === createdGoalId);
          if (createdGoal) {
            this.selectedGoal.set(createdGoal);
            this.resetGoalCreated();
          }
        } else if (response.goals.length > 0 && !this.selectedGoal()) {
          // Only set first goal as selected if no goal is currently selected
          this.selectedGoal.set(goals[0]);
        }

        this.notifySelectedGoalChange();
      },
      error: (err) => {
        console.error('Failed to load group goals', err);
      },
    });
  }

  loadMembersGoalContribution(
    groupID: number,
    startDate: string,
    endDate: string
  ) {
    this.getGroupMembersGoalContribution(groupID, startDate, endDate).subscribe(
      {
        next: (response) => {
          this.membersContribution.set(response.members);
        },
        error: (err) => {
          console.error('Failed to load members goal contribution', err);
        },
      }
    );
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
    this.memberAddedOrRemoved.set(true);
  }
  resetMemberAddedOrRemoved() {
    this.memberAddedOrRemoved.set(false);
  }

  setSelectedGoal(goal: Goal | null): void {
    this.selectedGoal.set(goal);
    this.notifySelectedGoalChange();
  }

  refreshGoals(groupID: number) {
    this.loadGoals(groupID);
  }

  notifyGoalUpdated(goalID: number) {
    const selectedGroup = this.selectedGroup();
    if (selectedGroup) {
      this.loadGoals(selectedGroup.id);
    }
  }

  deleteGoal(goalId: number): Observable<any> {
    const params = new HttpParams().set('goalID', goalId);
    return this.http.delete<any>('/api/group-goal', { params });
  }

  notifyGoalDeleted(goalId: number) {
    const selectedGroup = this.selectedGroup();
    if (selectedGroup) {
      // If the deleted goal was selected, clear the selection
      if (this.selectedGoal()?.id === goalId) {
        this.selectedGoal.set(null);
      }
      this.loadGoals(selectedGroup.id);
    }
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
  description?: string;
  goal_type: 'distance' | 'elevation' | 'summit_count' | 'specific_summits';
  target_value: number; // For distance/elevation/summit count
  target_summits?: number[]; // Array of peak IDs for specific summits goal
  start_date: string;
  end_date: string;
  created_at: string;
  updated_at?: string;
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
}

export interface CreateGroupResponse {
  group_id: number;
}

export interface UpdateGroupRequest {
  id: number;
  name: string;
}

export interface GetGroupsResponse {
  groups: Group[];
}

export interface CreateGoalRequest {
  group_id: number;
  name: string;
  description?: string;
  goal_type: 'distance' | 'elevation' | 'summit_count' | 'specific_summits';
  target_value: number;
  target_summits?: number[];
  start_date: string;
  end_date: string;
}

export interface CreateGoalResponse {
  goal_id: number;
}

export interface UpdateGoalRequest {
  id: number;
  group_id: number;
  name?: string;
  description?: string;
  goal_type?: 'distance' | 'elevation' | 'summit_count' | 'specific_summits';
  target_value?: number;
  target_summits?: number[];
  start_date?: string;
  end_date?: string;
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
  role: string;
}

export interface LeaveGroupRequest {
  groupID: number;
}
