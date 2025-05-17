import { Component, Input, signal, WritableSignal } from '@angular/core';
import { Goal } from 'src/app/services/groups.service';

@Component({
  selector: 'goals-create-form',
  standalone: true,
  imports: [],
  templateUrl: './goals-create-form.component.html',
  styleUrl: './goals-create-form.component.scss',
})
export class GoalsCreateFormComponent {
  constructor() { }

  @Input({ required: true }) formSignal!: WritableSignal<{ show: boolean, data: Goal | null }>;
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
  isEditMode = signal(false);

  ngOnInit(): void {

  }

  handleClose() {
    this.formSignal.set({ show: false, data: null });
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
