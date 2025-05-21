import { Component, computed, effect, Input } from '@angular/core';
import { Goal, GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-progress-bar',
  imports: [],
  templateUrl: './groups-progress-bar.component.html',
  styleUrl: './groups-progress-bar.component.scss',
})
export class GroupsProgressBarComponent {
  @Input({ required: true }) goal!: Goal;

  constructor(private groupService: GroupService) { }

  membersContribution = this.groupService.membersContribution;

  totalSummits = computed(() => {
    const members = this.membersContribution();
    return members?.reduce((sum, member) => sum + (member.total_unique_summits ?? 0), 0) ?? 0;
  })

  summitProgress = computed(() => {
    const total = this.totalSummits();
    if (!this.goal) return 0;
    return Math.min(100, Math.round((total / Number(this.goal.target_value)) * 100));
  })

}
