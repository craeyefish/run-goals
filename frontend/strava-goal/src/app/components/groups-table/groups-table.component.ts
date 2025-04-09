import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Group, GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-table',
  imports: [CommonModule],
  templateUrl: './groups-table.component.html',
  styleUrl: './groups-table.component.scss',
})
export class GroupsTableComponent {
  groups: Group[] = [];

  constructor(private groupService: GroupService) { }

  ngOnInit() {
    const userID = 1;  // todo: get all user data from authservice signal or storage ?
    this.groupService.getGroups(userID).subscribe({
      next: (response) => {
        this.groups = response.groups;
      },
      error: (err) => {
        console.error('Failed to load groups', err)
      }
    })
  }
}
