import { Component, computed, effect, inject } from '@angular/core';
import { GroupService, MemberContribution } from 'src/app/services/groups.service';
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

    const groupTotalUniqueSummits = members.reduce((sum, member) => sum + (member.total_unique_summits ?? 0), 0);

    return members.map(member => {
      const userId = member.user_id;
      const role = member.role;
      const totalActivities = member.total_activities;
      const totalDistance = member.total_distance;
      const totalUniqueSummits = member.total_unique_summits ?? 0;
      const totalSummits = member.total_summits;
      const contributionPercentage = groupTotalUniqueSummits === 0
        ? 0
        : Math.round(((member.total_unique_summits ?? 0) / groupTotalUniqueSummits) * 100);
      return {
        userId,
        role,
        totalActivities,
        totalDistance,
        totalUniqueSummits,
        totalSummits,
        contributionPercentage
      }
    })
  })

  constructor() {
    effect(() => {
      const group = this.groupService.selectedGroup();
      const goal = this.groupService.selectedGoal();
      const selectedGoalChange = this.groupService.selectedGoalChange();
      const memberChange = this.groupService.memberAddedOrRemoved();

      if ((selectedGoalChange || memberChange) && (group && goal)) {
        this.groupService.getGroupMembersGoalContribution(group.id, goal.start_date!, goal.end_date!).subscribe({
          next: (response) => this.groupService.membersContribution.set(response.members),
          error: (err) => console.error('Failed to load members', err)
        });
        this.groupService.resetSelectedGoalChange();
      } else if (group && !goal) {
        this.groupService.getGroupMembers(group.id).subscribe({
          next: (response) => {
            const members = response.members
            const membersContribution: MemberContribution[] = members.map(member => ({
              group_member_id: member.id,
              group_id: member.group_id,
              user_id: member.user_id,
              role: member.role,
              joined_at: member.joined_at,
              total_activities: null,
              total_distance: null,
              total_unique_summits: null,
              total_summits: null,
            }));
            this.groupService.membersContribution.set(membersContribution);
          },
          error: (err) => console.error('Failed to load members', err)
        })
      } else {
        this.groupService.membersContribution.set([]);
      }
    }
    )
  }

  selectedGoal = this.groupService.selectedGoal;
}
