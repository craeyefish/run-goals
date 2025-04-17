import { Component, computed, effect, inject } from '@angular/core';
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

  membersComputed = computed(() => {
    const members = this.groupService.membersContribution();

    if (!members) return [];

    const totalUniqueSummits = members.reduce((sum, member) => sum + member.total_unique_summits, 0);
    if (totalUniqueSummits === 0) {
      return members.map(member => ({
        ...member,
        contributionPercentage: 0
      }))
    } else {
      return members.map(member => ({
        ...member,
        contributionPercentage: Math.round((member.total_unique_summits / totalUniqueSummits) * 100)
      }))
    }
  })

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
