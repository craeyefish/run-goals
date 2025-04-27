import { Component, computed, effect } from '@angular/core';
import { GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-progress-bar',
  imports: [],
  templateUrl: './groups-progress-bar.component.html',
  styleUrl: './groups-progress-bar.component.scss',
})
export class GroupsProgressBarComponent {
  constructor(private groupService: GroupService) { }

  selectedGoal = this.groupService.selectedGoal;
  membersContribution = this.groupService.membersContribution;

  totalSummits = computed(() => {
    const members = this.membersContribution();
    return members?.reduce((sum, member) => sum + (member.total_unique_summits ?? 0), 0) ?? 0;
  })

  summitProgress = computed(() => {
    const goal = this.selectedGoal();
    const total = this.totalSummits();
    if (!goal) return 0;
    return Math.min(100, Math.round((total / Number(goal.target_value)) * 100));
  })

}
