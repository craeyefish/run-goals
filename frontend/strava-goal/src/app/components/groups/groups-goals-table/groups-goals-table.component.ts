import { CommonModule } from '@angular/common';
import { Component, effect, inject, Input, signal } from '@angular/core';
import { Goal, GroupService } from 'src/app/services/groups.service';
import { GroupsProgressBarComponent } from '../groups-progress-bar/groups-progress-bar.component';

@Component({
  selector: 'groups-goals-table',
  imports: [CommonModule, GroupsProgressBarComponent],
  templateUrl: './groups-goals-table.component.html',
  styleUrl: './groups-goals-table.component.scss',
})
export class GroupsGoalsTableComponent {
  public groupService = inject(GroupService);

  @Input() onEditGoalClick?: (goal: Goal) => void;

  constructor() {
    effect(() => {
      const selectedGroup = this.groupService.selectedGroup();
      if (selectedGroup) {
        this.groupService.loadGoals(selectedGroup.id);
      }
    })

    effect(() => {
      if (this.groupService.goalCreated()) {
        const selectedGroup = this.groupService.selectedGroup();
        if (selectedGroup) {
          this.groupService.loadGoals(selectedGroup.id);
        }
      }
    })
  }

  selectGoal(goal: Goal) {
    this.groupService.selectedGoal.set(goal);
    this.groupService.notifySelectedGoalChange();
  }
}
