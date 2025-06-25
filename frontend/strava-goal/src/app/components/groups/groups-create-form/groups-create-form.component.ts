import { Component, EventEmitter, Output, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { CreateGroupRequest } from '../../../services/groups.service';

@Component({
  selector: 'groups-create-form',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './groups-create-form.component.html',
  styleUrls: ['./groups-create-form.component.scss'],
})
export class GroupsCreateFormComponent {
  @Output() onSubmit = new EventEmitter<CreateGroupRequest>();
  @Output() onCancel = new EventEmitter<void>();

  groupName = signal('');

  handleClose() {
    this.onCancel.emit();
  }

  handleSubmit() {
    // Validation
    if (!this.groupName().trim()) {
      alert('Please enter a group name');
      return;
    }

    const groupData: CreateGroupRequest = {
      name: this.groupName().trim(),
    };

    this.onSubmit.emit(groupData);
  }
}
