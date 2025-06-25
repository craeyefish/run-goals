import {
  Component,
  EventEmitter,
  Output,
  Input,
  signal,
  OnInit,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Goal, UpdateGoalRequest } from '../../../services/groups.service';
import { GoalsPeakSelectorComponent } from '../goals-peak-selector/goals-peak-selector.component';

@Component({
  selector: 'goals-edit-form',
  standalone: true,
  imports: [CommonModule, FormsModule, GoalsPeakSelectorComponent],
  templateUrl: './goals-edit-form.component.html',
  styleUrls: ['./goals-edit-form.component.scss'],
})
export class GoalsEditFormComponent implements OnInit {
  @Output() onSubmit = new EventEmitter<UpdateGoalRequest>();
  @Output() onCancel = new EventEmitter<void>();
  @Input() goal!: Goal;

  // Track if we're currently dragging
  private isDragging = false;

  goalName = signal('');
  goalType = signal<
    'distance' | 'elevation' | 'summit_count' | 'specific_summits'
  >('distance');
  targetValue = signal<number>(0);
  targetSummits = signal<number[]>([]);
  startDate = signal('');
  endDate = signal('');
  description = signal('');

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

  ngOnInit(): void {
    // Populate form with existing goal data
    if (this.goal) {
      this.goalName.set(this.goal.name);
      this.goalType.set(this.goal.goal_type);
      this.targetValue.set(this.goal.target_value);
      this.description.set(this.goal.description || '');

      // Convert dates from backend format to yyyy-mm-dd format for HTML inputs
      this.startDate.set(this.formatDateForInput(this.goal.start_date));
      this.endDate.set(this.formatDateForInput(this.goal.end_date));

      if (
        this.goal.goal_type === 'specific_summits' &&
        this.goal.target_summits
      ) {
        this.targetSummits.set(this.goal.target_summits);
      } else {
        this.targetSummits.set([]);
      }
    }
  }

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
    // Reset target value when changing types
    this.targetValue.set(0);
    this.targetSummits.set([]);
  }

  onTargetSummitsChange(summitIds: number[]) {
    this.targetSummits.set(summitIds);
    // For specific summits, set target value to the number of selected summits
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

  // Convert date from backend format to yyyy-mm-dd for HTML input
  private formatDateForInput(dateInput: string | Date): string {
    if (!dateInput) return '';

    const date =
      typeof dateInput === 'string' ? new Date(dateInput) : dateInput;
    if (isNaN(date.getTime())) return '';

    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  }

  // Helper function to format date for backend
  private formatDateForBackend(dateString: string): string {
    if (!dateString) return '';

    // Create a date object from the input string (yyyy-mm-dd)
    const date = new Date(dateString + 'T00:00:00.000Z');

    // Return in ISO format that the backend expects
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

    // Validate end date is after start date
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

    const updateData: UpdateGoalRequest = {
      id: this.goal.id,
      group_id: this.goal.group_id,
      name: this.goalName().trim(),
      description: this.description().trim(),
      goal_type: this.goalType(),
      target_value: this.targetValue(),
      start_date: this.formatDateForBackend(this.startDate()),
      end_date: this.formatDateForBackend(this.endDate()),
    };

    // Add target summits for specific summits goal
    if (this.goalType() === 'specific_summits') {
      updateData.target_summits = this.targetSummits();
    }

    console.log('Submitting goal update:', updateData); // Debug log

    this.onSubmit.emit(updateData);
  }
}
