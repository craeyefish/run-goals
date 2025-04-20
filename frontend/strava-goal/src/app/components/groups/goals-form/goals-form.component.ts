import { Component, Input, signal, WritableSignal } from '@angular/core';
import { Goal } from 'src/app/services/groups.service';

@Component({
  selector: 'goals-form',
  standalone: true,
  imports: [],
  templateUrl: './goals-form.component.html',
  styleUrl: './goals-form.component.scss',
})
export class GoalsFormComponent {
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
  @Input() initialValues?: Goal | null;
  @Input() mode: 'create' | 'edit' = 'create';


  goalName = signal('');
  targetValue = signal('');
  startDate = signal('');
  endDate = signal('');
  isEditMode = signal(false);

  ngOnInit(): void {
    if (this.initialValues) {
      this.goalName.set(this.initialValues.name);
      this.targetValue.set(this.initialValues.target_value);
      this.startDate.set(this.initialValues.start_date);
      this.endDate.set(this.initialValues.end_date);
      this.isEditMode.set(true);
    }
  }

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
