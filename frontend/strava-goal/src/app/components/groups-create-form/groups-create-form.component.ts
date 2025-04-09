import { Component, Input, signal, WritableSignal } from '@angular/core';

@Component({
  selector: 'groups-create-form',
  standalone: true,
  imports: [],
  templateUrl: './groups-create-form.component.html',
  styleUrl: './groups-create-form.component.scss',
})
export class GroupsCreateFormComponent {
  constructor() { }

  @Input({ required: true }) formVisible!: WritableSignal<boolean>;
  @Input({ required: true }) onSubmit!: (data: { name: string }) => void;

  groupName = signal('');

  handleClose() {
    this.formVisible.set(false);
  }

  handleSubmit() {
    this.onSubmit({
      name: this.groupName()
    });
  }
}
