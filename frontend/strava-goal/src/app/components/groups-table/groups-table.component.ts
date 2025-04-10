import { CommonModule } from '@angular/common';
import { Component, signal } from '@angular/core';
import { Group, GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-table',
  imports: [CommonModule],
  templateUrl: './groups-table.component.html',
  styleUrl: './groups-table.component.scss',
})
export class GroupsTableComponent {

  groups = signal<Group[]>([]);

  constructor(public groupService: GroupService) { }

  ngOnInit() {
    const userID = 1;  // todo: get all user data from authservice signal or storage ?
    this.groupService.getGroups(userID).subscribe({
      next: (response) => {
        this.groups.set(response.groups);
        if (response.groups.length > 0) {
          this.groupService.selectedGroup.set(response.groups[0]);
          console.log(this.groupService.selectedGroup())
        }
      },
      error: (err) => {
        console.error('Failed to load groups', err)
      }
    })
  }

  selectGroup(group: Group) {
    this.groupService.selectedGroup.set(group);
  }
}
