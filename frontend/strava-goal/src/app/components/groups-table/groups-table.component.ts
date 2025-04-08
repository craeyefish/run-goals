import { Component } from '@angular/core';
import { GetGroupsResponse, Group, GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-table',
  imports: [],
  templateUrl: './groups-table.component.html',
  styleUrl: './groups-table.component.scss',
})
export class GroupsTableComponent {
  groups: Group[] = [];

  constructor(private groupService: GroupService) { }

  ngOnInit() {
    const userID = 1;
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
