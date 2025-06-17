import { Component, OnInit } from '@angular/core';
import { PeakSummitService } from 'src/app/services/peak-summit.service';
import { PeakSummaries } from 'src/app/models/peak-summit.model';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-peak-summit-table',
  imports: [CommonModule],
  templateUrl: './peak-summit-table.component.html',
  styleUrls: ['./peak-summit-table.component.scss'],
})
export class PeakSummitTableComponent implements OnInit {
  peakSummaries: PeakSummaries[] = [];
  expandedPeaks: Set<number> = new Set();

  constructor(private peakSummitService: PeakSummitService) {}

  ngOnInit(): void {
    this.fetchPeakSummaries();
  }

  fetchPeakSummaries(): void {
    this.peakSummitService.getPeakSummaries().subscribe({
      next: (data) => {
        this.peakSummaries = data;
      },
      error: (err) => {
        console.error('Error fetching peak summaries:', err);
      },
    });
  }

  toggleExpanded(peakId: number): void {
    if (this.expandedPeaks.has(peakId)) {
      this.expandedPeaks.delete(peakId);
    } else {
      this.expandedPeaks.add(peakId);
    }
  }

  isExpanded(peakId: number): boolean {
    return this.expandedPeaks.has(peakId);
  }

  formatMovingTime(movingTime: number): string {
    const hours = Math.floor(movingTime / 3600);
    const minutes = Math.floor((movingTime % 3600) / 60);
    const seconds = Math.floor(movingTime % 60);

    if (hours > 0) {
      return `${hours}:${minutes.toString().padStart(2, '0')}:${seconds
        .toString()
        .padStart(2, '0')}`;
    } else {
      return `${minutes}:${seconds.toString().padStart(2, '0')}`;
    }
  }

  formatDistance(distance: number): string {
    const km = distance / 1000;
    return `${km.toFixed(2)} km`;
  }
}
