import { CommonModule } from "@angular/common";
import { Component, signal, WritableSignal } from "@angular/core";
import { GoalProgressComponent } from "src/app/components/goal-progress/goal-progress.component";
import { GoalsFormComponent } from "src/app/components/groups/goals-form/goals-form.component";
import { GroupsGoalsTableComponent } from "src/app/components/groups/groups-goals-table/groups-goals-table.component";
import { CreateGoalRequest, GroupService, UpdateGoalRequest } from "src/app/services/groups.service";

@Component({
  selector: 'group-details-page',
  standalone: true,
  imports: [
    CommonModule,
    GoalsFormComponent,
    GoalProgressComponent,
    GroupsGoalsTableComponent
  ],
  templateUrl: './group-details.component.html',
  styleUrls: ['./group-details.component.scss'],
})
export class GroupsDetailsPageComponent {

  constructor(
    private groupService: GroupService
  ) { }

  showCreateGoalForm: WritableSignal<boolean> = signal(false);
  showEditGoalForm: WritableSignal<boolean> = signal(false);

  openCreateGoalForm = () => this.showCreateGoalForm.set(true);
  openEditGoalForm = () => {
    this.showEditGoalForm.set(true);
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
        this.showEditGoalForm.set(false);
        this.groupService.notifyGoalCreated(selectedGoal.id);
      },
      error: (err) => {
        console.error('Error updating goal:', err)
      }
    })
  }

}
