import { Component, Input, signal, WritableSignal } from '@angular/core';
import { Group } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-create-form',
  standalone: true,
  imports: [],
  templateUrl: './groups-create-form.component.html',
  styleUrl: './groups-create-form.component.scss',
})
export class GroupsCreateFormComponent {
  constructor() { }

  @Input({ required: true }) formSignal!: WritableSignal<{ show: boolean, data: Group | null }>;
  @Input({ required: true }) onSubmit!: (data: { name: string }) => void;

  groupName = signal('');

  ngOnInit(): void {

  }

  handleClose() {
    this.formSignal.set({ show: false, data: null });
  }

  handleSubmit() {
    this.onSubmit({
      name: this.groupName(),
    })
  }
}
