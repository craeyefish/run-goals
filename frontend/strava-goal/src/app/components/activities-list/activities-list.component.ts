import { CommonModule, NgFor } from '@angular/common';
import { Component } from '@angular/core';
import { filter } from 'rxjs';
import { Activity, ActivityService } from 'src/app/services/activity.service';

@Component({
  selector: 'app-activities-list',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './activities-list.component.html',
  styleUrls: ['./activities-list.component.scss'],
})
export class ActivitiesListComponent {
  activities: Activity[] = [];
  pageSize = 14;
  currentPage = 1;

  constructor(private activityService: ActivityService) {}

  ngOnInit(): void {
    this.activityService.loadActivities(); // Will only load if not already loaded
    this.activityService.activities$
      .pipe(filter((acts) => acts !== null))
      .subscribe((acts) => {
        this.activities = acts!;
      });
  }

  get paginatedActivities(): Activity[] {
    const startIndex = (this.currentPage - 1) * this.pageSize;
    return this.activities.slice(startIndex, startIndex + this.pageSize);
  }

  totalPages(): number {
    return Math.ceil(this.activities.length / this.pageSize);
  }

  goToPage(page: number) {
    this.currentPage = page;
  }

  getDisplayedPages(): number[] {
    const pages: number[] = [];
    const total = this.totalPages();

    if (total <= 7) {
      // Show all pages if small number
      return Array.from({ length: total }, (_, i) => i + 1);
    }

    pages.push(1);

    if (this.currentPage > 4) {
      pages.push(-1); // Will render as "..."
    }

    const startPage = Math.max(2, this.currentPage - 1);
    const endPage = Math.min(total - 1, this.currentPage + 2);

    for (let i = startPage; i <= endPage; i++) {
      pages.push(i);
    }

    if (this.currentPage + 2 < total - 1) {
      pages.push(-1); // Will render as "..."
    }

    pages.push(total);

    return pages;
  }

  truncateName(name: string, maxLength: number = 30): string {
    if (!name) return '';
    return name.length > maxLength ? name.slice(0, maxLength) + '...' : name;
  }

  formatMovingTime(seconds: number): string {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    return [h, m, s].map((v) => (v < 10 ? '0' + v : v)).join(':');
  }

  formatDistance(meters: number): string {
    const km = meters / 1000;
    return km.toFixed(2).replace('.', ',') + 'km';
  }

  formatElevation(meters: number): string {
    return Math.round(meters) + 'm';
  }
}
