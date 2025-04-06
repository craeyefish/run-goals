import { Component } from '@angular/core';
import { GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-table',
  imports: [],
  templateUrl: './groups-table.component.html',
  styleUrl: './groups-table.component.scss',
})
export class GroupsTableComponent {
  constructor(private groupService: GroupService) { }
}
