import { Component } from '@angular/core';
import { GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-progress-bar',
  imports: [],
  templateUrl: './groups-progress-bar.component.html',
  styleUrl: './groups-progress-bar.component.scss',
})
export class GroupsProgressBarComponent {
  constructor(private groupService: GroupService) { }
}
