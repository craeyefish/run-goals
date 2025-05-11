import { Component, Input, signal, WritableSignal } from "@angular/core";

@Component({
  selector: 'groups-join-form',
  standalone: true,
  imports: [],
  templateUrl: './groups-join-form.component.html',
  styleUrl: './groups-join-form.component.scss',
})
export class GroupsJoinFormComponent {
  constructor() { }

  @Input({ required: true }) formVisible!: WritableSignal<boolean>;
  @Input({ required: true }) onSubmit!: (data: { code: string }) => void;

  code = signal('');

  handleClose() {
    this.formVisible.set(false);
  }

  handleSubmit() {
    this.onSubmit({
      code: this.code()
    });
  }





}
