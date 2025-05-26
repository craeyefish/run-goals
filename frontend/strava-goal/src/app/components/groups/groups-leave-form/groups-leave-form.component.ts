import { Component, Input, signal, WritableSignal } from "@angular/core";
import { Group } from "src/app/services/groups.service";

@Component({
  selector: 'groups-leave-form',
  standalone: true,
  imports: [],
  templateUrl: './groups-leave-form.component.html',
  styleUrl: './groups-leave-form.component.scss',
})
export class GroupsLeaveFormComponent {
  constructor() { }

  @Input({ required: true }) formSignal!: WritableSignal<{ show: boolean, data: Group | null }>;
  @Input({ required: true }) onSubmit!: (data: Group) => void;

  groupID = signal('');

  handleClose() {
    this.formSignal.set({ show: false, data: null });
  }

  handleSubmit() {
    const data = this.formSignal().data;
    if (data) {
      this.onSubmit(data)
    } else {
      console.log('error submitting leaving of group - data not found');
      return;
    }
  }





}
