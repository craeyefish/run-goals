import { Component, effect, inject, signal } from '@angular/core';
import { GroupService, Member } from 'src/app/services/groups.service';
import { GroupsTableComponent } from '../groups-table/groups-table.component';
import { GroupsGoalsTableComponent } from '../groups-goals-table/groups-goals-table.component';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'groups-members-table',
  imports: [CommonModule],
  templateUrl: './groups-members-table.component.html',
  styleUrl: './groups-members-table.component.scss',
})
export class GroupsMembersTableComponent {

  members = signal<Member[]>([]);

  constructor(private groupService: GroupService) {
    effect(() => {
      const group = this.groupService.selectedGroup();
      const goal = this.groupService.selectedGoal();

      if (group && goal) {
        this.groupService.getGoalProgress(group.id, goal.id).subscribe({
          next: (response) => this.members.set(response.members),
          error: (err) => console.error('Failed to load members', err)
        });
      } else {
        this.members.set([]);
      }
    })
  }
}
