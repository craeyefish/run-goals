import { CommonModule } from '@angular/common';
import { Component, signal, WritableSignal, OnInit, computed } from '@angular/core';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
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
import { ChallengeService } from 'src/app/services/challenge.service';
import { Challenge, ChallengeWithProgress } from 'src/app/models/challenge.model';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'group-details-page',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
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
  private destroy$ = new Subject<void>();

  constructor(
    private groupService: GroupService,
    private challengeService: ChallengeService,
    private breadcrumbService: BreadcrumbService,
    private router: Router
  ) { }

  createGoalFormSignal: WritableSignal<{ show: boolean; data: Goal | null }> =
    signal({ show: false, data: null });
  editGoalFormSignal: WritableSignal<{ show: boolean; data: Goal | null }> =
    signal({ show: false, data: null });
  deleteConfirmationSignal: WritableSignal<{
    show: boolean;
    data: Goal | null;
  }> = signal({ show: false, data: null });

  // Challenge-related signals
  groupChallenges = signal<Challenge[]>([]);
  availableChallenges = signal<Challenge[]>([]);
  filteredAvailableChallenges = signal<Challenge[]>([]);
  showAdoptChallengeModal = signal(false);
  loadingChallenges = signal(false);
  challengeSearchQuery = '';

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

    // Load group challenges
    this.loadGroupChallenges(selectedGroup.id);
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  // ==================== Challenge Methods ====================

  loadGroupChallenges(groupId: number): void {
    this.challengeService.getGroupChallenges(groupId).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (response) => {
        this.groupChallenges.set(response.challenges || []);
      },
      error: (err) => console.error('Error loading group challenges:', err)
    });
  }

  openAdoptChallengeModal(): void {
    this.showAdoptChallengeModal.set(true);
    this.loadAvailableChallenges();
  }

  closeAdoptChallengeModal(): void {
    this.showAdoptChallengeModal.set(false);
    this.challengeSearchQuery = '';
  }

  loadAvailableChallenges(): void {
    this.loadingChallenges.set(true);
    this.challengeService.getPublicChallenges().pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (response) => {
        this.availableChallenges.set(response.challenges || []);
        this.filteredAvailableChallenges.set(response.challenges || []);
        this.loadingChallenges.set(false);
      },
      error: (err) => {
        console.error('Error loading available challenges:', err);
        this.loadingChallenges.set(false);
      }
    });
  }

  filterAvailableChallenges(): void {
    const query = this.challengeSearchQuery.toLowerCase().trim();
    if (!query) {
      this.filteredAvailableChallenges.set(this.availableChallenges());
      return;
    }

    const filtered = this.availableChallenges().filter(c =>
      c.name.toLowerCase().includes(query) ||
      c.region?.toLowerCase().includes(query) ||
      c.description?.toLowerCase().includes(query)
    );
    this.filteredAvailableChallenges.set(filtered);
  }

  isChallengeAdopted(challengeId: number): boolean {
    return this.groupChallenges().some(c => c.id === challengeId);
  }

  adoptChallenge(challenge: Challenge): void {
    if (this.isChallengeAdopted(challenge.id)) return;

    const selectedGroup = this.selectedGroup();
    if (!selectedGroup) return;

    this.challengeService.addGroupToChallenge(challenge.id, { groupId: selectedGroup.id }).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: () => {
        console.log('Challenge adopted:', challenge.name);
        this.loadGroupChallenges(selectedGroup.id);
        this.closeAdoptChallengeModal();
      },
      error: (err) => {
        console.error('Error adopting challenge:', err);
        alert('Failed to adopt challenge. Please try again.');
      }
    });
  }

  removeGroupChallenge(challengeId: number, event: Event): void {
    event.stopPropagation();

    const selectedGroup = this.selectedGroup();
    if (!selectedGroup) return;

    if (!confirm('Remove this challenge from the group?')) return;

    this.challengeService.removeGroupFromChallenge(challengeId, selectedGroup.id).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: () => {
        console.log('Challenge removed from group');
        this.loadGroupChallenges(selectedGroup.id);
      },
      error: (err) => {
        console.error('Error removing challenge:', err);
        alert('Failed to remove challenge. Please try again.');
      }
    });
  }

  navigateToChallenge(challengeId: number): void {
    this.router.navigate(['/challenges', challengeId]);
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
