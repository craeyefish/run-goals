import { CommonModule } from "@angular/common";
import { Component, signal, WritableSignal } from "@angular/core";
import { GroupsCreateFormComponent } from "src/app/components/groups/groups-create-form/groups-create-form.component";
import { GroupsEditFormComponent } from "src/app/components/groups/groups-edit-form/groups-edit-form.component";
import { GroupsJoinFormComponent } from "src/app/components/groups/groups-join-form/groups-join-form.component";
import { GroupsTableComponent } from "src/app/components/groups/groups-table/groups-table.component";
import { CreateGroupRequest, Group, GroupService, UpdateGroupRequest } from "src/app/services/groups.service";

@Component({
  selector: 'group-list-page',
  standalone: true,
  imports: [
    CommonModule,
    GroupsTableComponent,
    GroupsCreateFormComponent,
    GroupsEditFormComponent,
    GroupsJoinFormComponent
  ],
  templateUrl: './group-list.component.html',
  styleUrls: ['./group-list.component.scss'],
})
export class GroupsListPageComponent {

  constructor(private groupService: GroupService) {
    this.groupService.loadGroups()
  }

  groups = this.groupService.groups;

  createGroupFormSignal: WritableSignal<{ show: boolean, data: Group | null }> = signal({ show: false, data: null });
  showJoinGroupForm: WritableSignal<boolean> = signal(false);
  editGroupFormSignal: WritableSignal<{ show: boolean, data: Group | null }> = signal({ show: false, data: null });

  openCreateGroupForm = () => this.createGroupFormSignal.set({ show: true, data: null });
  openJoinGroupForm = () => this.showJoinGroupForm.set(true);
  openEditGroupForm = (group: Group) => {
    this.editGroupFormSignal.set({ show: true, data: group });
  }

  onCreateGroupFormSubmit = (data: { name: string }) => {
    const requestPayload: CreateGroupRequest = {
      name: data.name,
      created_by: 1,
    };

    this.groupService.createGroup(requestPayload).subscribe({
      next: (response) => {
        console.log('Group Created:', response);
        this.createGroupFormSignal.set({ show: false, data: null });
        this.groupService.notifyGroupUpdate(response.group_id);
      },
      error: (err) => {
        console.error('Error creating group:', err)
      }
    })
  }

  onEditGroupFormSubmit = (data: Group) => {
    if (!data) {
      console.log('No group data');
      return;
    }
    const requestPayload: UpdateGroupRequest = {
      id: data.id,
      name: data.name,
      created_by: data.created_by,
      created_at: data.created_at,
    };

    this.groupService.updateGroup(requestPayload).subscribe({
      next: () => {
        console.log('Group Updated: ', data.id);
        this.editGroupFormSignal.set({ show: false, data: null });
        this.groupService.notifyGroupUpdate(data.id);
      },
      error: (err) => {
        console.error('Error updating group:', err)
      }
    })
  }

  onJoinGroupFormSubmit = (
    data: {
      code: string
    }
  ) => {

  }

}
