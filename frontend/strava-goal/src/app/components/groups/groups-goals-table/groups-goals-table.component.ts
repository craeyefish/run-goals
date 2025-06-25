import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { GroupsProgressBarComponent } from '../groups-progress-bar/groups-progress-bar.component';
import { GroupService, Goal } from '../../../services/groups.service';
import { PeakService, Peak } from '../../../services/peak.service';
import { GroupsSummitDetailsComponent } from '../groups-summit-details/groups-summit-details.component';

@Component({
  selector: 'groups-goals-table',
  standalone: true,
  imports: [
    CommonModule,
    GroupsProgressBarComponent,
    GroupsSummitDetailsComponent,
  ],
  templateUrl: './groups-goals-table.component.html',
  styleUrls: ['./groups-goals-table.component.scss'],
})
export class GroupsGoalsTableComponent {
  @Input() onEditGoalClick?: (goal: Goal) => void;
  @Input() onDeleteGoalClick?: (goal: Goal) => void;

  constructor(
    public groupService: GroupService,
    private peakService: PeakService
  ) {}

  // Updated selectGoal method to handle summit details toggle
  selectGoal(goal: Goal) {
    // If this is a specific summits goal, toggle the selection
    if (
      goal.goal_type === 'specific_summits' &&
      this.hasSpecificSummits(goal)
    ) {
      // If this goal is already selected, deselect it to close details
      if (this.groupService.selectedGoal() === goal) {
        this.groupService.setSelectedGoal(null);
      } else {
        // Otherwise, select this goal to show details
        this.groupService.setSelectedGoal(goal);
      }
    } else {
      // For other goal types, just select the goal
      this.groupService.setSelectedGoal(goal);
    }
  }

  // Format date to yyyy/mm/dd
  formatDate(dateString: string | Date): string {
    let date: Date;

    if (typeof dateString === 'string') {
      // Handle ISO string or yyyy-mm-dd format
      date = new Date(dateString);
    } else {
      date = dateString;
    }

    // Check if date is valid
    if (isNaN(date.getTime())) {
      return 'Invalid Date';
    }

    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}/${month}/${day}`;
  }

  // Format target value based on goal type
  formatTargetValue(goal: Goal): string {
    switch (goal.goal_type) {
      case 'distance':
        return `${goal.target_value} km`;
      case 'elevation':
        return `${goal.target_value} m`;
      case 'summit_count':
        return `${goal.target_value} peaks`;
      case 'specific_summits':
        return `${goal.target_summits?.length || 0} specific peaks`;
      default:
        return goal.target_value.toString();
    }
  }

  // Get goal type icon/label
  getGoalTypeIcon(goalType: string): string {
    switch (goalType) {
      case 'distance':
        return '🏃';
      case 'elevation':
        return '⛰️';
      case 'summit_count':
        return '🏔️';
      case 'specific_summits':
        return '🎯';
      default:
        return '📊';
    }
  }

  // Check if goal has specific summits to display
  hasSpecificSummits(goal: Goal): boolean {
    return !!(
      goal.goal_type === 'specific_summits' &&
      goal.target_summits &&
      goal.target_summits.length > 0
    );
  }

  // Check if summit details are currently expanded for this goal
  isSummitDetailsExpanded(goal: Goal): boolean {
    return (
      this.hasSpecificSummits(goal) && goal === this.groupService.selectedGoal()
    );
  }
}
