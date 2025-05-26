import { CommonModule } from "@angular/common";
import { Component, signal, WritableSignal } from "@angular/core";
import { GoalProgressComponent } from "src/app/components/goal-progress/goal-progress.component";
import { GoalsCreateFormComponent } from "src/app/components/groups/goals-create-form/goals-create-form.component";
import { GoalsEditFormComponent } from "src/app/components/groups/goals-edit-form/goals-edit-form.component";
import { GroupsGoalsTableComponent } from "src/app/components/groups/groups-goals-table/groups-goals-table.component";
import { GroupsMembersTableComponent } from "src/app/components/groups/groups-members-table/groups-members-table.component";
import { CreateGoalRequest, Goal, GroupService, UpdateGoalRequest } from "src/app/services/groups.service";

@Component({
  selector: 'group-details-page',
  standalone: true,
  imports: [
    CommonModule,
    GoalsCreateFormComponent,
    GoalsEditFormComponent,
    GroupsGoalsTableComponent,
    GroupsMembersTableComponent
  ],
  templateUrl: './group-details.component.html',
  styleUrls: ['./group-details.component.scss'],
})
export class GroupsDetailsPageComponent {

  constructor(
    private groupService: GroupService
  ) { }

  createGoalFormSignal: WritableSignal<{ show: boolean, data: Goal | null }> = signal({ show: false, data: null });
  editGoalFormSignal: WritableSignal<{ show: boolean, data: Goal | null }> = signal({ show: false, data: null });

  openCreateGoalForm = () => this.createGoalFormSignal.set({ show: true, data: null });
  openEditGoalForm = (goal: Goal) => {
    this.editGoalFormSignal.set({ show: true, data: goal });
  }

  selectedGroup = this.groupService.selectedGroup;
  selectedGoal = this.groupService.selectedGoal;

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
      start_date: new Date(data.startDate),
      end_date: new Date(data.endDate),
    };

    this.groupService.createGoal(requestPayload).subscribe({
      next: (response) => {
        console.log('Goal Created:', data);
        this.createGoalFormSignal.set({ show: false, data: null });
        this.groupService.notifyGoalCreated(response.goal_id);
      },
      error: (err) => {
        console.error('Error creating group goal:', err)
      }
    });
  }

  onEditGoalFormSubmit = (data: Goal) => {
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
      target_value: data.target_value,
      start_date: new Date(data.start_date!),
      end_date: new Date(data.end_date!),
      created_at: new Date(selectedGoal.created_at!),
    };

    this.groupService.updateGoal(requestPayload).subscribe({
      next: () => {
        console.log('Goal Updated: ', selectedGoal.id);
        this.editGoalFormSignal.set({ show: false, data: null });
        this.groupService.notifyGoalCreated(selectedGoal.id);
      },
      error: (err) => {
        console.error('Error updating goal:', err)
      }
    })
  }

}
