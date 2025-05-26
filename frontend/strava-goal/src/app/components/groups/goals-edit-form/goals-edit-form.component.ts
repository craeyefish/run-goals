import { Component, Input, signal, WritableSignal } from '@angular/core';
import { Goal } from 'src/app/services/groups.service';

@Component({
  selector: 'goals-edit-form',
  standalone: true,
  imports: [],
  templateUrl: './goals-edit-form.component.html',
  styleUrl: './goals-edit-form.component.scss',
})
export class GoalsEditFormComponent {
  constructor() { }

  @Input({ required: true }) formSignal!: WritableSignal<{ show: boolean, data: Goal | null }>;
  @Input({ required: true }) onSubmit!: (data: Goal) => void;


  goalName = signal('');
  targetValue = signal('');
  startDate = signal('');
  endDate = signal('');
  isEditMode = signal(false);

  ngOnInit(): void {
    const data = this.formSignal().data
    if (data) {
      this.goalName.set(data.name);
      this.targetValue.set(data.target_value);
      this.startDate.set(data.start_date!);
      this.endDate.set(data.end_date!);
    }
  }

  handleClose() {
    this.formSignal.set({ show: false, data: null });
  }

  handleSubmit() {
    const data = this.formSignal().data;
    if (data) {
      data.name = this.goalName();
      data.target_value = this.targetValue();
      data.start_date = this.startDate();
      data.end_date = this.endDate();
      this.onSubmit(data);
    } else {
      console.error('error submitting goal update - data not found');
      return;
    }
  }
}
