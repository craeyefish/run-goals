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
}
