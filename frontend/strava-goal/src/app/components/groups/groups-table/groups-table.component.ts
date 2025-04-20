import { CommonModule } from '@angular/common';
import { Component, effect, inject, Input, signal } from '@angular/core';
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

  constructor() {
    effect(() => {
      if (this.groupService.groupCreated()) {
        this.groupService.loadGroups();
      }
    })
  }

  ngOnInit() {
    this.groupService.loadGroups()
  }

  selectGroup(group: Group) {
    this.groupService.selectedGroup.set(group);
  }
}
