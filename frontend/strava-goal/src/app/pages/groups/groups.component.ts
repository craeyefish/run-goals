import { CommonModule } from '@angular/common';
import { Component, signal, WritableSignal } from '@angular/core';
import { AppButtonComponent } from 'src/app/components/app-button/app-button.component';
import { GoalsCreateFormComponent } from 'src/app/components/goals-create-form/goals-create-form.component';
import { GroupsCreateFormComponent } from 'src/app/components/groups-create-form/groups-create-form.component';
import { GroupsGoalsTableComponent } from 'src/app/components/groups-goals-table/groups-goals-table.component';
import { GroupsMembersTableComponent } from 'src/app/components/groups-members-table/groups-members-table.component';
import { GroupsProgressBarComponent } from 'src/app/components/groups-progress-bar/groups-progress-bar.component';
import { GroupsTableComponent } from 'src/app/components/groups-table/groups-table.component';
import { CreateGoalRequest, CreateGroupRequest, GroupService } from 'src/app/services/groups.service';


@Component({
  selector: 'app-groups',
  standalone: true,
  imports: [
    CommonModule,
    AppButtonComponent,
    GroupsTableComponent,
    GroupsMembersTableComponent,
    GroupsProgressBarComponent,
    GroupsCreateFormComponent,
    GoalsCreateFormComponent,
    GroupsGoalsTableComponent
  ],
  templateUrl: './groups.component.html',
  styleUrls: ['./groups.component.scss'],
})
export class GroupsComponent {

  constructor(private groupService: GroupService) { }

  showCreateGroupForm: WritableSignal<boolean> = signal(false);
  showCreateGoalForm: WritableSignal<boolean> = signal(false);

  openCreateGroupForm = () => this.showCreateGroupForm.set(true)
  openCreateGoalForm = () => this.showCreateGoalForm.set(true)

  onCreateGroupFormSubmit = (data: { name: string }) => {
    const requestPayload: CreateGroupRequest = {
      name: data.name,
      created_by: 1,
    };

    this.groupService.createGroup(requestPayload).subscribe({
      next: (response) => {
        console.log('Group Created:', data);
        this.showCreateGroupForm.set(false);
      },
      error: (err) => {
        console.error('Error creating group:', err)
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
        this.groupService.notifyGoalCreated();
      },
      error: (err) => {
        console.error('Error creating group goal:', err)
      }
    });
  }
}
