import { Component, OnInit } from '@angular/core';
import {
  ProgressService,
  GoalProgress,
  UserContribution,
} from '../services/progress.service';

@Component({
  selector: 'app-goal-progress',
  templateUrl: './goal-progress.component.html',
  styleUrls: ['./goal-progress.component.scss'],
})
export class GoalProgressComponent implements OnInit {
  goalProgress: GoalProgress | null = null;

  constructor(private progressService: ProgressService) {}

  ngOnInit(): void {
    this.progressService.getProgress().subscribe({
      next: (data) => {
        this.goalProgress = data;
      },
      error: (err) => {
        console.error('Error fetching progress', err);
      },
    });
  }
}
