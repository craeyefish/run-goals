import { Component } from '@angular/core';
import { GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-create-btn',
  imports: [],
  templateUrl: './groups-create-btn.component.html',
  styleUrl: './groups-create-btn.component.scss',
})
export class GroupsCreateBtnComponent {
  constructor(private groupService: GroupService) { }

  onCreateClick() {
    this.groupService.create();
  }
}
