import { CommonModule } from '@angular/common';
import { Component, signal, WritableSignal, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { GoalProgressComponent } from 'src/app/components/goal-progress/goal-progress.component';
import { GoalDeleteConfirmationComponent } from 'src/app/components/groups/goal-delete-confirmation/goal-delete-confirmation.component';
import { GoalsCreateFormComponent } from 'src/app/components/groups/goals-create-form/goals-create-form.component';
import { GoalsEditFormComponent } from 'src/app/components/groups/goals-edit-form/goals-edit-form.component';
import { GroupsGoalsTableComponent } from 'src/app/components/groups/groups-goals-table/groups-goals-table.component';
import { GroupsMembersTableComponent } from 'src/app/components/groups/groups-members-table/groups-members-table.component';
import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
import {
  CreateGoalRequest,
  Goal,
  GroupService,
  UpdateGoalRequest,
} from 'src/app/services/groups.service';

@Component({
  selector: 'group-details-page',
  standalone: true,
  imports: [
    CommonModule,
    GoalsCreateFormComponent,
    GoalsEditFormComponent,
    GoalDeleteConfirmationComponent,
    GroupsGoalsTableComponent,
    GroupsMembersTableComponent,
  ],
  templateUrl: './group-details.component.html',
  styleUrls: ['./group-details.component.scss'],
})
export class GroupsDetailsPageComponent implements OnInit {
  constructor(
    private groupService: GroupService,
    private breadcrumbService: BreadcrumbService,
    private router: Router
  ) {}

  createGoalFormSignal: WritableSignal<{ show: boolean; data: Goal | null }> =
    signal({ show: false, data: null });
  editGoalFormSignal: WritableSignal<{ show: boolean; data: Goal | null }> =
    signal({ show: false, data: null });
  deleteConfirmationSignal: WritableSignal<{
    show: boolean;
    data: Goal | null;
  }> = signal({ show: false, data: null });

  openCreateGoalForm = () =>
    this.createGoalFormSignal.set({ show: true, data: null });
  openEditGoalForm = (goal: Goal) => {
    this.editGoalFormSignal.set({ show: true, data: goal });
  };
  openDeleteConfirmation = (goal: Goal) => {
    this.deleteConfirmationSignal.set({ show: true, data: goal });
  };

  selectedGroup = this.groupService.selectedGroup;
  selectedGoal = this.groupService.selectedGoal;

  ngOnInit() {
    const selectedGroup = this.groupService.selectedGroup();

    // Check if we have a selected group
    if (!selectedGroup) {
      console.warn('No group selected, redirecting to groups list');
      // Redirect to groups list if no group is selected
      this.router.navigate(['/groups']);
      return;
    }

    // Set breadcrumb with the group name
    this.breadcrumbService.addItem({
      label: selectedGroup.name,
    });

    // Load group data if needed
    this.groupService.loadGoals(selectedGroup.id);
    this.groupService.getGroupMembers(selectedGroup.id);
  }

  onCreateGoalFormSubmit = (data: CreateGoalRequest) => {
    const selectedGroup = this.groupService.selectedGroup();
    if (!selectedGroup) {
      console.error('No selectedGroup available');
      return;
    }

    const requestPayload: CreateGoalRequest = {
      group_id: selectedGroup.id,
      name: data.name,
      description: data.description || '',
      goal_type: data.goal_type,
      target_summits: data.target_summits || [],
      target_value: Number(data.target_value),
      start_date: data.start_date,
      end_date: data.end_date,
    };

    this.groupService.createGoal(requestPayload).subscribe({
      next: (response) => {
        console.log('Goal Created:', response);
        this.createGoalFormSignal.set({ show: false, data: null });

        // Store the created goal ID and refresh
        this.groupService.notifyGoalCreated(response.goal_id);
        this.groupService.refreshGoals(selectedGroup.id);
      },
      error: (err) => {
        console.error('Error creating group goal:', err);
      },
    });
  };

  onEditGoalFormSubmit = (updateData: UpdateGoalRequest) => {
    const selectedGroup = this.selectedGroup();

    if (!selectedGroup) {
      console.error('No group selected');
      return;
    }

    this.groupService.updateGoal(updateData).subscribe({
      next: (response) => {
        console.log('Goal Updated:', updateData.id);
        this.editGoalFormSignal.set({ show: false, data: null });

        // Use the new refresh method
        this.groupService.refreshGoals(selectedGroup.id);
      },
      error: (err) => {
        console.error('Error updating goal:', err);
      },
    });
  };

  onDeleteGoalConfirm = () => {
    const goalToDelete = this.deleteConfirmationSignal().data;
    const selectedGroup = this.selectedGroup();

    if (!goalToDelete || !selectedGroup) {
      console.error('No goal or group selected for deletion');
      return;
    }

    this.groupService.deleteGoal(goalToDelete.id).subscribe({
      next: () => {
        console.log('Goal deleted:', goalToDelete.id);
        this.deleteConfirmationSignal.set({ show: false, data: null });
        this.groupService.notifyGoalDeleted(goalToDelete.id);
      },
      error: (err) => {
        console.error('Error deleting goal:', err);
        // You could show an error message here
        alert('Failed to delete goal. Please try again.');
      },
    });
  };

  onDeleteGoalCancel = () => {
    this.deleteConfirmationSignal.set({ show: false, data: null });
  };
}
