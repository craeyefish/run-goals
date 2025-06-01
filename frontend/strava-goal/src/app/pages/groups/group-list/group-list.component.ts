import { CommonModule } from "@angular/common";
import { Component, signal, WritableSignal } from "@angular/core";
import { Router } from "@angular/router";
import { GroupsCreateFormComponent } from "src/app/components/groups/groups-create-form/groups-create-form.component";
import { GroupsEditFormComponent } from "src/app/components/groups/groups-edit-form/groups-edit-form.component";
import { GroupsJoinFormComponent } from "src/app/components/groups/groups-join-form/groups-join-form.component";
import { GroupsLeaveFormComponent } from "src/app/components/groups/groups-leave-form/groups-leave-form.component";
import { GroupsTableComponent } from "src/app/components/groups/groups-table/groups-table.component";
import { AuthService } from "src/app/services/auth.service";
import { BreadcrumbService } from "src/app/services/breadcrumb.service";
import { CreateGroupMemberRequest, CreateGroupRequest, Group, GroupService, LeaveGroupRequest, UpdateGroupRequest } from "src/app/services/groups.service";

@Component({
  selector: 'group-list-page',
  standalone: true,
  imports: [
    CommonModule,
    GroupsTableComponent,
    GroupsCreateFormComponent,
    GroupsEditFormComponent,
    GroupsJoinFormComponent,
    GroupsLeaveFormComponent
  ],
  templateUrl: './group-list.component.html',
  styleUrls: ['./group-list.component.scss'],
})
export class GroupsListPageComponent {

  constructor(
    private groupService: GroupService,
    private authService: AuthService,
    private breadcrumbService: BreadcrumbService,
    private router: Router,
  ) {
    this.groupService.loadGroups()
  }

  groups = this.groupService.groups;

  createGroupFormSignal: WritableSignal<{ show: boolean, data: Group | null }> = signal({ show: false, data: null });
  joinGroupFormSignal: WritableSignal<{ show: boolean, code: string | null }> = signal({ show: false, code: null });
  editGroupFormSignal: WritableSignal<{ show: boolean, data: Group | null }> = signal({ show: false, data: null });
  leaveGroupFormSignal: WritableSignal<{ show: boolean, data: Group | null }> = signal({ show: false, data: null });

  openCreateGroupForm = () => this.createGroupFormSignal.set({ show: true, data: null });
  openJoinGroupForm = () => this.joinGroupFormSignal.set({ show: true, code: null });
  openEditGroupForm = (group: Group) => this.editGroupFormSignal.set({ show: true, data: group });
  openLeaveGroupForm = (group: Group) => this.leaveGroupFormSignal.set({ show: true, data: group });

  ngOnInit() {
    this.breadcrumbService.setItems(
      [
        {
          label: 'Groups', routerLink: '/groups', callback: () => {
            this.groupService.selectedGroup.set(null);
            this.router.navigate(['/groups']);
          }
        }
      ]
    )
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
        this.groupService.notifyGroupUpdate();
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
        this.groupService.notifyGroupUpdate();
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
    const requestPayload: CreateGroupMemberRequest = {
      group_code: data.code,
      user_id: this.authService.getUserID()!,
      role: 'member',
    };

    this.groupService.createGroupMember(requestPayload).subscribe({
      next: () => {
        console.log('Group Joined');
        this.joinGroupFormSignal.set({ show: false, code: null });
        this.groupService.notifyGroupUpdate();
      },
      error: (err) => {
        console.error('Error joining group:', err)
      }
    })
  }

  onLeaveGroupFormSubmit = (data: Group) => {
    const requestPayload: LeaveGroupRequest = {
      groupID: data.id,
      userID: this.authService.getUserID()!,
    };

    this.groupService.leaveGroup(requestPayload).subscribe({
      next: () => {
        console.log('Left Group: ', data.name);
        this.leaveGroupFormSignal.set({ show: false, data: null });
        this.groupService.notifyGroupUpdate();
      },
      error: (err) => {
        console.error('Error leaving group:', err)
      }
    })
  }

}
