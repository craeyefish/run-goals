import { Component } from '@angular/core';
import { GroupsCreateBtnComponent } from 'src/app/components/groups-create-btn/groups-create-btn.component';
import { GroupsTableComponent } from 'src/app/components/groups-table/groups-table.component';

@Component({
  selector: 'app-groups',
  standalone: true,
  imports: [GroupsCreateBtnComponent, GroupsTableComponent],
  templateUrl: './groups.component.html',
  styleUrls: ['./groups.component.scss'],
})
export class GroupsComponent { }
