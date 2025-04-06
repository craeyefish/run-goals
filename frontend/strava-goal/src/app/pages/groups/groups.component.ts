import { Component } from '@angular/core';
import { GroupsCreateBtnComponent } from 'src/app/components/groups-create-btn/groups-create-btn.component';
import { GroupsMembersTableComponent } from 'src/app/components/groups-members-table/groups-members-table.component';
import { GroupsProgressBarComponent } from 'src/app/components/groups-progress-bar/groups-progress-bar.component';
import { GroupsTableComponent } from 'src/app/components/groups-table/groups-table.component';

@Component({
  selector: 'app-groups',
  standalone: true,
  imports: [GroupsCreateBtnComponent, GroupsTableComponent, GroupsMembersTableComponent, GroupsProgressBarComponent],
  templateUrl: './groups.component.html',
  styleUrls: ['./groups.component.scss'],
})
export class GroupsComponent { }
