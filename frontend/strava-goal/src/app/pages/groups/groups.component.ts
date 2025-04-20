import { CommonModule } from '@angular/common';
import { Component, signal, WritableSignal } from '@angular/core';
import { AppButtonComponent } from 'src/app/components/app/app-button/app-button.component';
import { GoalsFormComponent } from 'src/app/components/groups/goals-form/goals-form.component';
import { GroupsFormComponent } from 'src/app/components/groups/groups-form/groups-form.component';
import { GroupsGoalsTableComponent } from 'src/app/components/groups/groups-goals-table/groups-goals-table.component';
import { GroupsMembersTableComponent } from 'src/app/components/groups/groups-members-table/groups-members-table.component';
import { GroupsProgressBarComponent } from 'src/app/components/groups/groups-progress-bar/groups-progress-bar.component';
import { GroupsTableComponent } from 'src/app/components/groups/groups-table/groups-table.component';
import { CreateGoalRequest, CreateGroupRequest, CreateGroupResponse, Group, GroupService, UpdateGoalRequest, UpdateGroupRequest } from 'src/app/services/groups.service';


@Component({
  selector: 'app-groups',
  standalone: true,
  imports: [
    CommonModule,
    AppButtonComponent,
    GroupsTableComponent,
    GroupsMembersTableComponent,
    GroupsProgressBarComponent,
    GroupsFormComponent,
    GoalsFormComponent,
    GroupsGoalsTableComponent
  ],
  templateUrl: './groups.component.html',
  styleUrls: ['./groups.component.scss'],
})
export class GroupsComponent {

  constructor(private groupService: GroupService) { }

  showCreateGroupForm: WritableSignal<boolean> = signal(false);
  showEditGroupForm: WritableSignal<boolean> = signal(false);
  showCreateGoalForm: WritableSignal<boolean> = signal(false);
  showEditGoalForm: WritableSignal<boolean> = signal(false);

  openCreateGroupForm = () => this.showCreateGroupForm.set(true);
  openEditGroupForm = () => {
    this.showEditGroupForm.set(true);
  }
  openCreateGoalForm = () => this.showCreateGoalForm.set(true);
  openEditGoalForm = () => {
    this.showEditGoalForm.set(true);
  }

  selectedGroup = this.groupService.selectedGroup;
  selectedGoal = this.groupService.selectedGoal;

  onCreateGroupFormSubmit = (data: { name: string }) => {
    const requestPayload: CreateGroupRequest = {
      name: data.name,
      created_by: 1,
    };

    this.groupService.createGroup(requestPayload).subscribe({
      next: (response) => {
        console.log('Group Created:', response);
        this.showCreateGroupForm.set(false);
        this.groupService.notifyGroupCreated(response.group_id);
      },
      error: (err) => {
        console.error('Error creating group:', err)
      }
    })
  }

  onEditGroupFormSubmit = (data: { name: string }) => {
    const selectedGroup = this.selectedGroup();
    if (!selectedGroup) {
      console.log('No selectedGroup');
      return;
    }
    const requestPayload: UpdateGroupRequest = {
      id: selectedGroup.id,
      name: data.name,
      created_by: selectedGroup.created_by,
      created_at: selectedGroup.created_at,
    };

    this.groupService.updateGroup(requestPayload).subscribe({
      next: () => {
        console.log('Group Updated: ', selectedGroup.id);
        this.groupService.loadGroups();
        this.showEditGroupForm.set(false);
      },
      error: (err) => {
        console.error('Error updating group:', err)
      }
    })
  }

  onCreateGoalFormSubmit = (
    data: {
      name: string,
      targetValue: string,
      startDate: string,
      endDate: string
    }
  ) => {
    const selectedGroup = this.groupService.selectedGroup();
    if (!selectedGroup) {
      console.log('No selectedGroup');
      return;
    }
    const requestPayload: CreateGoalRequest = {
      group_id: selectedGroup.id,
      name: data.name,
      target_value: data.targetValue,
      start_date: data.startDate,
      end_date: data.endDate,
    };

    this.groupService.createGoal(requestPayload).subscribe({
      next: (response) => {
        console.log('Goal Created:', data);
        this.showCreateGoalForm.set(false);
        this.groupService.notifyGoalCreated(response.goal_id);
      },
      error: (err) => {
        console.error('Error creating group goal:', err)
      }
    });
  }

  onEditGoalFormSubmit = (
    data: {
      name: string,
      targetValue: string,
      startDate: string,
      endDate: string
    }
  ) => {
    const selectedGroup = this.selectedGroup();
    const selectedGoal = this.selectedGoal();
    if (!selectedGroup || !selectedGoal) {
      console.log('Goal or group not selected');
      return;
    }
    const requestPayload: UpdateGoalRequest = {
      id: selectedGoal.id,
      group_id: selectedGroup.id,
      name: data.name,
      target_value: data.targetValue,
      start_date: data.startDate,
      end_date: data.endDate,
      created_at: selectedGoal.created_at,
    };

    this.groupService.updateGoal(requestPayload).subscribe({
      next: () => {
        console.log('Goal Updated: ', selectedGoal.id);
        this.groupService.loadGoals(selectedGroup.id);
        this.showEditGoalForm.set(false);
      },
      error: (err) => {
        console.error('Error updating goal:', err)
      }
    })
  }
}
