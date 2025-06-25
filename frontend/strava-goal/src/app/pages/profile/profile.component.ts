import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ProfileService, UserProfile } from '../../services/profile.service';
import { ActivityService } from '../../services/activity.service';
import { PeakSummitService } from '../../services/peak-summit.service';
import { Subject, combineLatest } from 'rxjs';
import { takeUntil, filter } from 'rxjs/operators';

@Component({
  selector: 'app-profile',
  imports: [CommonModule],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss'],
})
export class ProfileComponent implements OnInit, OnDestroy {
  profile: UserProfile | null = null;
  loading = true;
  error: string | null = null;
  totalDistance = 0;
  totalSummits = 0;

  private destroy$ = new Subject<void>();

  constructor(
    private profileService: ProfileService,
    private activityService: ActivityService,
    private peakSummitService: PeakSummitService
  ) {}

  ngOnInit(): void {
    this.loadProfile();
    this.loadYearStats();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  loadProfile(): void {
    this.profileService
      .getUserProfile()
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (profile) => {
          this.profile = profile;
          this.loading = false;
        },
        error: (err) => {
          this.error = 'Failed to load profile';
          this.loading = false;
          console.error('Error loading profile:', err);
        },
      });
  }

  loadYearStats(): void {
    // Trigger loading for both services
    this.activityService.loadActivities();

    // Combine both observables and wait for both to have data
    combineLatest([
      this.activityService.activities$.pipe(
        filter((activities) => activities !== null)
      ),
      this.peakSummitService.getPeakSummaries(),
    ])
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: ([activities, peakSummaries]) => {
          const currentYear = new Date().getFullYear();

          // Calculate total distance for current year
          const yearActivities = activities!.filter((activity) => {
            const activityYear = new Date(activity.start_date).getFullYear();
            return activityYear === currentYear;
          });

          this.totalDistance = parseFloat(
            yearActivities
              .reduce((total, activity) => {
                return total + activity.distance / 1000; // Convert to km
              }, 0)
              .toFixed(1)
          );

          // Calculate total summits for current year
          let summitCount = 0;
          peakSummaries.forEach((peak) => {
            peak.activities.forEach((activity) => {
              const activityYear = new Date(activity.start_date).getFullYear();
              if (activityYear === currentYear) {
                summitCount++;
              }
            });
          });

          this.totalSummits = summitCount;
        },
        error: (err) => {
          console.error('Error loading year stats:', err);
        },
      });
  }

  getDaysSinceLastUpdate(): number {
    if (!this.profile?.last_updated) return 0;
    const lastUpdate = new Date(this.profile.last_updated);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - lastUpdate.getTime());
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  }
}
