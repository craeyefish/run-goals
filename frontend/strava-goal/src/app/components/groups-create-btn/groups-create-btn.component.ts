import { Component, Input } from '@angular/core';
import { GroupService } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-create-btn',
  standalone: true,
  imports: [],
  templateUrl: './groups-create-btn.component.html',
  styleUrl: './groups-create-btn.component.scss',
})
export class GroupsCreateBtnComponent {
  constructor() { }

  @Input() onClick?: () => void;
}
