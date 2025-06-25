import {
  Component,
  Input,
  Output,
  EventEmitter,
  OnInit,
  OnDestroy,
  signal,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { PeakService, Peak } from '../../../services/peak.service';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'goals-peak-selector',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './goals-peak-selector.component.html',
  styleUrls: ['./goals-peak-selector.component.css'],
})
export class GoalsPeakSelectorComponent implements OnInit, OnDestroy {
  @Input() selectedPeakIds: number[] = [];
  @Output() selectedPeakIdsChange = new EventEmitter<number[]>();

  searchTerm = signal('');
  allPeaks = signal<Peak[]>([]);
  filteredPeaks = signal<Peak[]>([]);
  selectedPeaks = signal<Peak[]>([]);
  isLoading = signal(false);

  private destroy$ = new Subject<void>();

  constructor(private peakService: PeakService) {}

  ngOnInit() {
    this.loadPeaks();
    // Initialize selected peaks from input
    this.updateSelectedPeaks();
  }

  ngOnChanges() {
    this.updateSelectedPeaks();
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  loadPeaks() {
    this.isLoading.set(true);

    // Use the existing service pattern
    this.peakService.loadPeaks();

    this.peakService.peaks$.pipe(takeUntil(this.destroy$)).subscribe({
      next: (peaks) => {
        if (peaks) {
          this.allPeaks.set(peaks);
          this.filteredPeaks.set(peaks.slice(0, 50)); // Show first 50 peaks initially
          this.isLoading.set(false);
          this.updateSelectedPeaks();
        }
      },
      error: (err) => {
        console.error('Failed to load peaks:', err);
        this.isLoading.set(false);
      },
    });
  }

  onSearchChange() {
    const term = this.searchTerm().toLowerCase();
    if (term.length === 0) {
      this.filteredPeaks.set(this.allPeaks().slice(0, 50));
    } else {
      const filtered = this.allPeaks()
        .filter(
          (peak) =>
            peak.name.toLowerCase().includes(term) ||
            peak.elevation_meters.toString().includes(term)
        )
        .slice(0, 20); // Limit search results
      this.filteredPeaks.set(filtered);
    }
  }

  togglePeakSelection(peak: Peak) {
    const currentIds = [...this.selectedPeakIds];
    const index = currentIds.indexOf(peak.id);

    if (index > -1) {
      currentIds.splice(index, 1);
    } else {
      currentIds.push(peak.id);
    }

    this.selectedPeakIdsChange.emit(currentIds);
  }

  isPeakSelected(peakId: number): boolean {
    return this.selectedPeakIds.includes(peakId);
  }

  private updateSelectedPeaks() {
    const selected = this.allPeaks().filter((peak) =>
      this.selectedPeakIds.includes(peak.id)
    );
    this.selectedPeaks.set(selected);
  }

  removePeak(peakId: number) {
    const currentIds = this.selectedPeakIds.filter((id) => id !== peakId);
    this.selectedPeakIdsChange.emit(currentIds);
  }
}
