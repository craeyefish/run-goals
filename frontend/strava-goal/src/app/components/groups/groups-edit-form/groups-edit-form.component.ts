import { Component, Input, signal, WritableSignal } from '@angular/core';
import { Group } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-edit-form',
  standalone: true,
  imports: [],
  templateUrl: './groups-edit-form.component.html',
  styleUrl: './groups-edit-form.component.scss',
})
export class GroupsEditFormComponent {
  constructor() { }

  @Input({ required: true }) formSignal!: WritableSignal<{ show: boolean, data: Group | null }>;
  @Input({ required: true }) onSubmit!: (data: Group) => void;

  groupName = signal('');

  ngOnInit(): void {
    const data = this.formSignal().data;
    if (data) {
      this.groupName.set(data.name);
    }
  }

  handleClose() {
    this.formSignal.set({ show: false, data: null });
  }

  handleSubmit() {
    const data = this.formSignal().data;
    if (data) {
      data.name = this.groupName();
      this.onSubmit(data)
    } else {
      console.log('error submitting group update - data not found');
      return;
    }
  }
}
