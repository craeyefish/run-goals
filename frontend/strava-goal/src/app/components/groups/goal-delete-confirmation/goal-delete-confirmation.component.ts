import { Component, EventEmitter, Output, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Goal } from '../../../services/groups.service';

@Component({
  selector: 'goal-delete-confirmation',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './goal-delete-confirmation.component.html',
  styleUrls: ['./goal-delete-confirmation.component.css'],
})
export class GoalDeleteConfirmationComponent {
  @Output() onConfirm = new EventEmitter<void>();
  @Output() onCancel = new EventEmitter<void>();
  @Input() goal!: Goal;

  // Track if we're currently dragging to prevent modal close on drag
  private isDragging = false;

  // Handle overlay click with drag detection
  onOverlayClick(event: MouseEvent) {
    if (!this.isDragging && event.target === event.currentTarget) {
      this.handleCancel();
    }
    this.isDragging = false;
  }

  onMouseDown() {
    this.isDragging = false;
  }

  onMouseMove() {
    this.isDragging = true;
  }

  onMouseUp() {
    setTimeout(() => {
      this.isDragging = false;
    }, 10);
  }

  handleConfirm() {
    this.onConfirm.emit();
  }

  handleCancel() {
    this.onCancel.emit();
  }

  getGoalTypeLabel(): string {
    switch (this.goal.goal_type) {
      case 'distance':
        return 'Distance Goal';
      case 'elevation':
        return 'Elevation Goal';
      case 'summit_count':
        return 'Summit Count Goal';
      case 'specific_summits':
        return 'Specific Summits Goal';
      default:
        return 'Goal';
    }
  }

  getTargetDescription(): string {
    switch (this.goal.goal_type) {
      case 'distance':
        return `${this.goal.target_value} km`;
      case 'elevation':
        return `${this.goal.target_value} m elevation`;
      case 'summit_count':
        return `${this.goal.target_value} summits`;
      case 'specific_summits':
        return `${this.goal.target_summits?.length || 0} specific summits`;
      default:
        return this.goal.target_value.toString();
    }
  }
}
