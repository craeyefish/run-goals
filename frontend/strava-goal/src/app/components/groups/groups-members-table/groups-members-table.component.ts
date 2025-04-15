import { Component, effect, inject } from '@angular/core';
import { GroupService } from 'src/app/services/groups.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'groups-members-table',
  imports: [CommonModule],
  templateUrl: './groups-members-table.component.html',
  styleUrl: './groups-members-table.component.scss',
})
export class GroupsMembersTableComponent {
  public groupService = inject(GroupService);

  constructor() {
    effect(() => {
      const group = this.groupService.selectedGroup();
      const goal = this.groupService.selectedGoal();
      const selectedGoalChange = this.groupService.selectedGoalChange();
      const memberChange = this.groupService.memberAddedOrRemoved();

      if ((selectedGoalChange || memberChange) && (group && goal)) {
        this.groupService.getGroupMembersGoalContribution(group.id, goal.start_date, goal.end_date).subscribe({
          next: (response) => this.groupService.membersContribution.set(response.members),
          error: (err) => console.error('Failed to load members', err)
        });
        this.groupService.resetSelectedGoalChange();
      } else {
        this.groupService.membersContribution.set([]);
      }
    })
  }
}
