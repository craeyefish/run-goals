import { Component, Input, signal, WritableSignal } from '@angular/core';
import { Group } from 'src/app/services/groups.service';

@Component({
  selector: 'groups-form',
  standalone: true,
  imports: [],
  templateUrl: './groups-form.component.html',
  styleUrl: './groups-form.component.scss',
})
export class GroupsFormComponent {
  constructor() { }

  @Input({ required: true }) formVisible!: WritableSignal<boolean>;
  @Input({ required: true }) onSubmit!: (data: { name: string }) => void;
  @Input() initialValues?: Group | null;
  @Input() mode: 'create' | 'edit' = 'create';

  groupName = signal('');
  isEditMode = signal(false);

  ngOnInit(): void {
    if (this.initialValues) {
      this.groupName.set(this.initialValues.name);
      this.isEditMode.set(true);
    }
  }

  handleClose() {
    this.formVisible.set(false);
  }

  handleSubmit() {
    this.onSubmit({
      name: this.groupName()
    });
  }
}
