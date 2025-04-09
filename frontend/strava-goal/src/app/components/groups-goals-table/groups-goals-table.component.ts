import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Goal, GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-goals-table',
  imports: [CommonModule],
  templateUrl: './groups-goals-table.component.html',
  styleUrl: './groups-goals-table.component.scss',
})
export class GroupsGoalsTableComponent {
  goals: Goal[] = [];

  constructor(private groupService: GroupService) { }

  ngOnInit() {
    const groupID = 1; // todo: get selected groupID
    this.groupService.getGroupGoals(groupID).subscribe({
      next: (response) => {
        this.goals = response.goals;
      },
      error: (err) => {
        console.error('Failed to load group goals', err)
      }
    })
  }
}
