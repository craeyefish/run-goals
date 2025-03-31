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
}
