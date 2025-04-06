import { Component } from '@angular/core';
import { GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-members-table',
  imports: [],
  templateUrl: './groups-members-table.component.html',
  styleUrl: './groups-members-table.component.scss',
})
export class GroupsMembersTableComponent {
  constructor(private groupService: GroupService) { }
}
