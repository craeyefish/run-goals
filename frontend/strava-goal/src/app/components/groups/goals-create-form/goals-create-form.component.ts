import { Component, Input, signal, WritableSignal } from '@angular/core';

@Component({
  selector: 'goals-create-form',
  standalone: true,
  imports: [],
  templateUrl: './goals-create-form.component.html',
  styleUrl: './goals-create-form.component.scss',
})
export class GoalsCreateFormComponent {
  constructor() { }

  @Input({ required: true }) formVisible!: WritableSignal<boolean>;
  @Input({ required: true }) onSubmit!: (
    data: {
      name: string,
      targetValue: string,
      startDate: string,
      endDate: string
    }
  ) => void;

  goalName = signal('');
  targetValue = signal('');
  startDate = signal('');
  endDate = signal('');

  handleClose() {
    this.formVisible.set(false);
  }

  handleSubmit() {
    this.onSubmit({
      name: this.goalName(),
      targetValue: this.targetValue(),
      startDate: this.startDate(),
      endDate: this.endDate()
    });
  }
}
