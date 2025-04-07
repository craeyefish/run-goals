import { CommonModule } from '@angular/common';
import { Component, signal, WritableSignal } from '@angular/core';
import { GroupsCreateBtnComponent } from 'src/app/components/groups-create-btn/groups-create-btn.component';
import { GroupsCreateFormComponent } from 'src/app/components/groups-create-form/groups-create-form.component';
import { GroupsMembersTableComponent } from 'src/app/components/groups-members-table/groups-members-table.component';
import { GroupsProgressBarComponent } from 'src/app/components/groups-progress-bar/groups-progress-bar.component';
import { GroupsTableComponent } from 'src/app/components/groups-table/groups-table.component';
import { CreateGroupRequest, GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'app-groups',
  standalone: true,
  imports: [
    CommonModule,
    GroupsCreateBtnComponent,
    GroupsTableComponent,
    GroupsMembersTableComponent,
    GroupsProgressBarComponent,
    GroupsCreateFormComponent
  ],
  templateUrl: './groups.component.html',
  styleUrls: ['./groups.component.scss'],
})
export class GroupsComponent {

  constructor(private groupService: GroupService) { }

  showCreateGroupForm: WritableSignal<boolean> = signal(false);

  openCreateGroupForm = () => this.showCreateGroupForm.set(true)

  onCreateGroupFormSubmit = (data: { name: string }) => {
    const requestPayload: CreateGroupRequest = {
      name: data.name,
      created_by: 1,
    };

    this.groupService.createGroup(requestPayload).subscribe({
      next: (response) => {
        console.log('Group Created:', data);
        this.showCreateGroupForm.set(false);
      },
      error: (err) => {
        console.error('oops, a error', err)
      }
    })
  }
}
