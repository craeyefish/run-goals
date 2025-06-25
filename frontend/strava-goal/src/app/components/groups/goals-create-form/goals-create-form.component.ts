import { Component, EventEmitter, Output, Input, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { CreateGoalRequest } from '../../../services/groups.service';
import { GoalsPeakSelectorComponent } from '../goals-peak-selector/goals-peak-selector.component';

@Component({
  selector: 'goals-create-form',
  standalone: true,
  imports: [CommonModule, FormsModule, GoalsPeakSelectorComponent],
  templateUrl: './goals-create-form.component.html',
  styleUrls: ['./goals-create-form.component.scss'],
})
export class GoalsCreateFormComponent {
  @Output() onSubmit = new EventEmitter<CreateGoalRequest>();
  @Output() onCancel = new EventEmitter<void>();
  @Input() groupId!: number;

  // Track if we're currently dragging
  private isDragging = false;

  goalName = signal('');
  goalType = signal<
    'distance' | 'elevation' | 'summit_count' | 'specific_summits'
  >('distance');
  targetValue = signal<number>(0);
  targetSummits = signal<number[]>([]);
  startDate = signal<string>('');
  endDate = signal<string>('');

  goalTypes = [
    {
      value: 'distance',
      label: 'Distance Goal (km)',
      placeholder: 'e.g., 500',
    },
    {
      value: 'elevation',
      label: 'Elevation Goal (m)',
      placeholder: 'e.g., 10000',
    },
    {
      value: 'summit_count',
      label: 'Number of Summits',
      placeholder: 'e.g., 20',
    },
    {
      value: 'specific_summits',
      label: 'Specific Summits',
      placeholder: 'Select peaks below',
    },
  ];

  // Handle overlay click with drag detection
  onOverlayClick(event: MouseEvent) {
    // Only close if we're not dragging and the click target is the overlay itself
    if (!this.isDragging && event.target === event.currentTarget) {
      this.handleClose();
    }
    // Reset drag state
    this.isDragging = false;
  }

  // Track when dragging starts
  onMouseDown() {
    this.isDragging = false; // Reset on new mouse down
  }

  // Track when mouse moves (indicates dragging)
  onMouseMove() {
    this.isDragging = true;
  }

  // Reset drag state on mouse up
  onMouseUp() {
    // Don't reset immediately - let the click handler check the state first
    setTimeout(() => {
      this.isDragging = false;
    }, 10);
  }

  handleClose() {
    this.onCancel.emit();
  }

  onGoalTypeChange(type: string) {
    this.goalType.set(type as any);
    this.targetValue.set(0);
    this.targetSummits.set([]);
  }

  onTargetSummitsChange(summitIds: number[]) {
    this.targetSummits.set(summitIds);
    if (this.goalType() === 'specific_summits') {
      this.targetValue.set(summitIds.length);
    }
  }

  getTargetValuePlaceholder(): string {
    const selectedType = this.goalTypes.find(
      (type) => type.value === this.goalType()
    );
    return selectedType?.placeholder || '';
  }

  isTargetValueRequired(): boolean {
    return this.goalType() !== 'specific_summits';
  }

  private formatDateForBackend(dateString: string): string {
    if (!dateString) return '';
    const date = new Date(dateString + 'T00:00:00.000Z');
    return date.toISOString();
  }

  handleSubmit() {
    // Validation
    if (!this.goalName().trim()) {
      alert('Please enter a goal name');
      return;
    }

    if (!this.startDate()) {
      alert('Please select a start date');
      return;
    }

    if (!this.endDate()) {
      alert('Please select an end date');
      return;
    }

    if (new Date(this.endDate()) <= new Date(this.startDate())) {
      alert('End date must be after start date');
      return;
    }

    if (this.goalType() === 'specific_summits') {
      if (this.targetSummits().length === 0) {
        alert('Please select at least one summit');
        return;
      }
    } else {
      if (this.targetValue() <= 0) {
        alert('Please enter a valid target value');
        return;
      }
    }

    const goalData: CreateGoalRequest = {
      group_id: this.groupId,
      name: this.goalName().trim(),
      goal_type: this.goalType(),
      target_value: this.targetValue(),
      start_date: this.formatDateForBackend(this.startDate()),
      end_date: this.formatDateForBackend(this.endDate()),
    };

    if (this.goalType() === 'specific_summits') {
      goalData.target_summits = this.targetSummits();
    }

    this.onSubmit.emit(goalData);
  }
}
