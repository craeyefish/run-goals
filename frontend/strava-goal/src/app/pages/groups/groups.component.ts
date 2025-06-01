import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { GroupService } from 'src/app/services/groups.service';
import { Router, RouterOutlet } from '@angular/router';
import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
import { BreadcrumbComponent } from 'src/app/components/breadcrumb/breadcrumb.component';


@Component({
  selector: 'app-groups',
  standalone: true,
  imports: [
    CommonModule,
    RouterOutlet,
    BreadcrumbComponent
  ],
  templateUrl: './groups.component.html',
  styleUrls: ['./groups.component.scss'],
})
export class GroupsComponent {

  constructor(
    private groupService: GroupService,
    private router: Router,
  ) { }

  groups = this.groupService.groups;
  selectedGroup = this.groupService.selectedGroup;

  resetSelectedGroup() {
    this.selectedGroup.set(null);
    this.router.navigate(['/groups']);
  }
}
