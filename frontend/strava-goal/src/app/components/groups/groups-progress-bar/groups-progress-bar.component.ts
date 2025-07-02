import { Component, computed, effect, Input, signal } from '@angular/core';
import {
  Goal,
  GroupService,
  MemberContribution,
} from 'src/app/services/groups.service';

@Component({
  selector: 'groups-progress-bar',
  imports: [],
  templateUrl: './groups-progress-bar.component.html',
  styleUrl: './groups-progress-bar.component.scss',
})
export class GroupsProgressBarComponent {
  @Input({ required: true }) goal!: Goal;

  constructor(private groupService: GroupService) {}

  // Each progress bar maintains its own member contribution data
  private goalMembersContribution = signal<MemberContribution[]>([]);

  // Load member contribution data specific to this goal's date range
  private loadGoalContribution = effect(() => {
    const selectedGroup = this.groupService.selectedGroup();
    if (selectedGroup && this.goal) {
      this.groupService
        .getGroupMembersGoalContribution(
          selectedGroup.id,
          this.goal.start_date,
          this.goal.end_date
        )
        .subscribe({
          next: (response) =>
            this.goalMembersContribution.set(response.members),
          error: (err) =>
            console.error(
              'Failed to load goal-specific members contribution',
              err
            ),
        });
    }
  });

  // Calculate current value based on goal type
  currentValue = computed(() => {
    const members = this.goalMembersContribution();
    if (!members || !this.goal) return 0;

    switch (this.goal.goal_type) {
      case 'distance':
        // Distance is stored in meters, convert to kilometers for display
        return (
          members.reduce(
            (sum, member) => sum + (member.total_distance ?? 0),
            0
          ) / 1000
        );
      case 'elevation':
        // Note: This would need elevation data from members - currently not available
        return 0; // TODO: Add elevation tracking
      case 'summit_count':
        return members.reduce(
          (sum, member) => sum + (member.total_unique_summits ?? 0),
          0
        );
      case 'specific_summits':
        return members.reduce(
          (sum, member) => sum + (member.total_unique_summits ?? 0),
          0
        );
      default:
        return 0;
    }
  });

  // Calculate progress percentage based on goal type
  progressPercentage = computed(() => {
    const current = this.currentValue();
    if (!this.goal || this.goal.target_value <= 0) return 0;

    return Math.min(
      100,
      Math.round((current / Number(this.goal.target_value)) * 100)
    );
  });

  // Display text for progress (current value)
  displayValue = computed(() => {
    const current = this.currentValue();
    const target = this.goal?.target_value || 0;

    switch (this.goal?.goal_type) {
      case 'distance':
        return `${current.toFixed(1)} / ${target} km`;
      case 'elevation':
        return `${current} / ${target} m`;
      case 'summit_count':
        return `${current} / ${target}`;
      case 'specific_summits':
        return `${current} / ${target}`;
      default:
        return `${current} / ${target}`;
    }
  });
}
