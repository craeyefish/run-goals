import { Component, Input, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Goal } from '../../../services/groups.service';
import { PeakService, Peak } from '../../../services/peak.service';

interface SummitStatus {
  peak: Peak;
  isCompleted: boolean;
  completedBy: string[]; // User names who completed this summit
  completedDate?: string;
}

@Component({
  selector: 'groups-summit-details',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './groups-summit-details.component.html',
  styleUrls: ['./groups-summit-details.component.css'],
})
export class GroupsSummitDetailsComponent implements OnInit {
  @Input() goal!: Goal;

  summitStatuses = signal<SummitStatus[]>([]);
  loading = signal(true);

  constructor(private peakService: PeakService) {}

  ngOnInit() {
    this.loadSummitStatuses();
  }

  async loadSummitStatuses() {
    if (!this.goal.target_summits || this.goal.target_summits.length === 0) {
      this.loading.set(false);
      return;
    }

    try {
      // Load all peaks first
      this.peakService.loadPeaks();

      this.peakService.peaks$.subscribe((peaks) => {
        if (peaks) {
          const targetPeaks = peaks.filter((peak) =>
            this.goal.target_summits?.includes(peak.id)
          );

          // For each peak, check completion status
          const statuses: SummitStatus[] = targetPeaks.map((peak) => ({
            peak,
            isCompleted: peak.is_summited, // This assumes the peak has user context
            completedBy: [], // You'll need to implement group member checking
            completedDate: undefined,
          }));

          this.summitStatuses.set(statuses);
          this.loading.set(false);
        }
      });
    } catch (error) {
      console.error('Error loading summit statuses:', error);
      this.loading.set(false);
    }
  }

  getCompletionPercentage(): number {
    const statuses = this.summitStatuses();
    if (statuses.length === 0) return 0;

    const completed = statuses.filter((s) => s.isCompleted).length;
    return Math.round((completed / statuses.length) * 100);
  }
}
