import { Component } from '@angular/core';
import { GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-goals-table',
  imports: [],
  templateUrl: './groups-goals-table.component.html',
  styleUrl: './groups-goals-table.component.scss',
})
export class GroupsGoalsTableComponent {
  constructor(private groupService: GroupService) { }
}
