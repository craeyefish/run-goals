import { CommonModule } from '@angular/common';
import { Component, effect, inject, Input, signal } from '@angular/core';
import { Router } from '@angular/router';
import { Group, GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-table',
  imports: [CommonModule],
  templateUrl: './groups-table.component.html',
  styleUrl: './groups-table.component.scss',
})
export class GroupsTableComponent {
  public groupService = inject(GroupService);

  @Input() onEditGroupClick?: (group: Group) => void;
  @Input() onLeaveGroupClick?: (group: Group) => void;

  constructor(private router: Router) {
    effect(() => {
      if (this.groupService.groupUpdate()) {
        this.groupService.loadGroups();
      }
    })
  }

  selectGroup(group: Group) {
    this.groupService.selectedGroup.set(group);
    this.router.navigate(['/groups', group.code])
  }
}
