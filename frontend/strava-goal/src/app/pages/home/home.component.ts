import {
  Component,
  OnInit,
  OnDestroy,
  ViewChild,
  ElementRef,
  AfterViewInit,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ActivityService, Activity } from '../../services/activity.service';
import { PeakSummitService } from '../../services/peak-summit.service';
import { PeakService, Peak } from '../../services/peak.service';
import { PersonalGoalsService, PersonalYearlyGoal } from '../../services/personal-goals.service';
import { SummitFavouritesService } from '../../services/summit-favourites.service';
import { PeakPickerComponent, SelectedPeak } from '../../components/peak-picker/peak-picker.component';
import { YearlyGoalsComponent } from '../../components/yearly-goals/yearly-goals.component';
import { Subject, combineLatest } from 'rxjs';
import { takeUntil, filter } from 'rxjs/operators';
import {
  Chart,
  registerables,
  ScriptableContext,
  TooltipCallbacks,
  ChartConfiguration,
} from 'chart.js';
import 'chartjs-adapter-date-fns';

Chart.register(...registerables);

interface RecentSummit {
  peakId: number;
  peakName: string;
  summitedAt: Date;
  activityId: number;
}

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, FormsModule, PeakPickerComponent, YearlyGoalsComponent],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent implements OnInit, OnDestroy, AfterViewInit {
  @ViewChild('distanceChart') distanceChartCanvas!: ElementRef<HTMLCanvasElement>;
  @ViewChild('elevationChart') elevationChartCanvas!: ElementRef<HTMLCanvasElement>;

  loading = true;
  error: string | null = null;
  currentYear = new Date().getFullYear();
  selectedYear = this.currentYear;
  availableYears: number[] = [];

  // Stats
  totalDistance = 0;
  totalElevation = 0;
  totalActivities = 0;
  totalSummits = 0;

  // Personal goals from API
  personalGoals: PersonalYearlyGoal | null = null;

  // Summit favourites (new year-independent wishlist)
  favouriteIds: number[] = [];
  favouritePeaks: any[] = [];
  completedFavourites: any[] = [];
  incompleteFavourites: any[] = [];
  wishlistTab: 'incomplete' | 'complete' = 'incomplete';
  wishlistCollapsed = true; // Start collapsed by default

  // Recent summits
  recentSummits: RecentSummit[] = [];

  // Modal states
  showGoalEditor = false;
  editingGoalType: 'distance' | 'elevation' = 'distance';
  editingGoalValue = 0;
  showSummitPicker = false;

  // Charts
  distanceChart: Chart | null = null;
  elevationChart: Chart | null = null;

  // Store activities for chart creation
  yearActivities: Activity[] = [];
  private allActivities: Activity[] = [];
  private peakSummaries: any[] = [];

  private destroy$ = new Subject<void>();

  constructor(
    private activityService: ActivityService,
    private peakSummitService: PeakSummitService,
    private peakService: PeakService,
    private personalGoalsService: PersonalGoalsService,
    private summitFavouritesService: SummitFavouritesService,
    private router: Router
  ) { }

  ngOnInit(): void {
    // Load wishlist collapse preference from localStorage
    const savedCollapsed = localStorage.getItem('wishlistCollapsed');
    if (savedCollapsed !== null) {
      this.wishlistCollapsed = savedCollapsed === 'true';
    }

    this.loadDashboardData();
  }

  ngAfterViewInit(): void {
    // Charts will be created after data loads
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
    if (this.distanceChart) {
      this.distanceChart.destroy();
    }
    if (this.elevationChart) {
      this.elevationChart.destroy();
    }
  }

  loadDashboardData(): void {
    this.activityService.loadActivities();

    // Load personal goals
    this.personalGoalsService.getCurrentYearGoals().pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (goals) => {
        this.personalGoals = goals;
      },
      error: (err) => console.error('Error loading personal goals:', err)
    });

    // Load summit favourites
    this.summitFavouritesService.getFavourites().pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (favouriteIds) => {
        this.favouriteIds = favouriteIds;
        this.loadFavouritePeaks();
      },
      error: (err) => console.error('Error loading summit favourites:', err)
    });

    // Load active challenges
    // Removed active challenges section from home page

    combineLatest([
      this.activityService.activities$.pipe(
        filter((activities) => activities !== null)
      ),
      this.peakSummitService.getPeakSummaries(),
    ])
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: ([activities, peakSummaries]) => {
          this.allActivities = activities || [];
          this.peakSummaries = peakSummaries || [];

          // Extract available years from activities
          this.extractAvailableYears();

          // Calculate stats for selected year
          this.calculateStatsForYear(this.selectedYear);
          this.loadRecentSummits(this.peakSummaries);

          // Store for chart creation
          this.yearActivities = this.allActivities.filter(
            (activity) => new Date(activity.start_date).getFullYear() === this.selectedYear
          ).sort(
            (a, b) => new Date(a.start_date).getTime() - new Date(b.start_date).getTime()
          );

          setTimeout(() => {
            this.createDistanceChart();
            this.createElevationChart();
          }, 100);

          this.loading = false;
        },
        error: (err) => {
          this.error = 'Failed to load dashboard data';
          this.loading = false;
          console.error('Error loading dashboard:', err);
        },
      });
  }

  loadFavouritePeaks(): void {
    if (!this.favouriteIds.length) {
      this.favouritePeaks = [];
      this.completedFavourites = [];
      this.incompleteFavourites = [];
      return;
    }

    // Load peaks data
    this.peakService.loadPeaks();

    combineLatest([
      this.peakService.peaks$.pipe(filter(peaks => peaks !== null)),
      this.peakSummitService.getPeakSummaries()
    ]).pipe(
      takeUntil(this.destroy$)
    ).subscribe(([allPeaks, summaries]) => {
      this.favouritePeaks = this.favouriteIds.map(peakId => {
        const peakInfo = allPeaks!.find((p: Peak) => p.id === peakId);
        const summitData = (summaries || []).find((s: any) => s.peak_id === peakId);
        const latestActivity = summitData?.activities?.[0]; // Most recent summit

        return {
          id: peakId,
          name: peakInfo?.name || summitData?.peak_name || 'Unknown Peak',
          elevation: peakInfo?.elevation_meters || 0,
          completed: !!summitData?.activities?.length,
          summitedAt: latestActivity ? new Date(latestActivity.start_date) : undefined
        };
      });

      // Split into completed and incomplete
      this.completedFavourites = this.favouritePeaks
        .filter(p => p.completed)
        .sort((a, b) => {
          if (!a.summitedAt || !b.summitedAt) return 0;
          return b.summitedAt.getTime() - a.summitedAt.getTime(); // Most recent first
        });

      this.incompleteFavourites = this.favouritePeaks
        .filter(p => !p.completed)
        .sort((a, b) => a.name.localeCompare(b.name));
    });
  }

  loadRecentSummits(peakSummaries: any[]): void {
    const allSummits: RecentSummit[] = [];

    // Extract all summits from the selected year
    peakSummaries.forEach((peak: any) => {
      if (peak.activities && Array.isArray(peak.activities)) {
        peak.activities.forEach((activity: any) => {
          const activityYear = new Date(activity.start_date).getFullYear();
          if (activityYear === this.selectedYear) {
            allSummits.push({
              peakId: peak.peak_id,
              peakName: peak.peak_name,
              summitedAt: new Date(activity.start_date),
              activityId: activity.id
            });
          }
        });
      }
    });

    // Sort by date descending - keep all summits for the expanded view
    this.recentSummits = allSummits
      .sort((a, b) => b.summitedAt.getTime() - a.summitedAt.getTime());
  }

  extractAvailableYears(): void {
    const years = new Set<number>();
    years.add(this.currentYear); // Always include current year

    this.allActivities.forEach(activity => {
      const year = new Date(activity.start_date).getFullYear();
      years.add(year);
    });

    // Sort descending (most recent first)
    this.availableYears = Array.from(years).sort((a, b) => b - a);
  }

  selectYear(year: number): void {
    this.selectedYear = year;
    this.calculateStatsForYear(year);
    this.loadRecentSummits(this.peakSummaries);

    // Update year activities for charts
    this.yearActivities = this.allActivities.filter(
      (activity) => new Date(activity.start_date).getFullYear() === year
    ).sort(
      (a, b) => new Date(a.start_date).getTime() - new Date(b.start_date).getTime()
    );

    // Recreate charts
    if (this.distanceChart) this.distanceChart.destroy();
    if (this.elevationChart) this.elevationChart.destroy();
    setTimeout(() => {
      this.createDistanceChart();
      this.createElevationChart();
    }, 100);
  }

  calculateStatsForYear(year: number): void {
    const yearActivities = this.allActivities.filter((activity) => {
      const activityYear = new Date(activity.start_date).getFullYear();
      return activityYear === year;
    });

    this.totalActivities = yearActivities.length;
    this.totalDistance = parseFloat(
      yearActivities
        .reduce((total, activity) => total + activity.distance / 1000, 0)
        .toFixed(1)
    );

    this.totalElevation = parseFloat(
      yearActivities
        .reduce((total, activity) => total + (activity.total_elevation_gain || 0), 0)
        .toFixed(0)
    );

    let summitCount = 0;
    if (this.peakSummaries && Array.isArray(this.peakSummaries)) {
      this.peakSummaries.forEach((peak) => {
        if (peak.activities && Array.isArray(peak.activities)) {
          peak.activities.forEach((activity: any) => {
            const activityYear = new Date(activity.start_date).getFullYear();
            if (activityYear === year) {
              summitCount++;
            }
          });
        }
      });
    }
    this.totalSummits = summitCount;
  }

  calculateStats(activities: Activity[], peakSummaries: any[]): void {
    this.calculateStatsForYear(this.selectedYear);
  }

  // Goal getters with defaults
  get distanceGoal(): number {
    return this.personalGoals?.distance_goal || 1000;
  }

  get elevationGoal(): number {
    return this.personalGoals?.elevation_goal || 25000;
  }

  get summitGoal(): number {
    return this.personalGoals?.summit_goal || 20;
  }

  get distanceProgress(): number {
    return Math.min(100, (this.totalDistance / this.distanceGoal) * 100);
  }

  get elevationProgress(): number {
    return Math.min(100, (this.totalElevation / this.elevationGoal) * 100);
  }

  // Modal close handlers - only close if mousedown started on overlay
  onOverlayMouseDown(event: MouseEvent): void {
    if (event.target === event.currentTarget) {
      this.showGoalEditor = false;
    }
  }

  onSummitOverlayMouseDown(event: MouseEvent): void {
    if (event.target === event.currentTarget) {
      this.showSummitPicker = false;
    }
  }

  // Goal editing
  openGoalEditor(type: 'distance' | 'elevation'): void {
    this.editingGoalType = type;
    this.editingGoalValue = type === 'distance' ? this.distanceGoal : this.elevationGoal;
    this.showGoalEditor = true;
  }

  saveGoal(): void {
    if (!this.personalGoals) return;

    const updates: PersonalYearlyGoal = {
      ...this.personalGoals,
      distance_goal: this.editingGoalType === 'distance' ? this.editingGoalValue : this.personalGoals.distance_goal,
      elevation_goal: this.editingGoalType === 'elevation' ? this.editingGoalValue : this.personalGoals.elevation_goal,
    };

    this.personalGoalsService.saveGoals(updates).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (updated: PersonalYearlyGoal) => {
        this.personalGoals = updated;
        this.showGoalEditor = false;
        // Recreate charts with new goals
        if (this.distanceChart) this.distanceChart.destroy();
        if (this.elevationChart) this.elevationChart.destroy();
        this.createDistanceChart();
        this.createElevationChart();
      },
      error: (err: any) => console.error('Error saving goal:', err)
    });
  }

  // Summit wishlist management
  openSummitPicker(): void {
    this.showSummitPicker = true;
  }

  closeSummitPicker(): void {
    this.showSummitPicker = false;
  }

  onPeakAdded(peak: SelectedPeak): void {
    this.summitFavouritesService.addFavourite(peak.id).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (favouriteIds: number[]) => {
        this.favouriteIds = favouriteIds;
        this.loadFavouritePeaks();
      },
      error: (err: any) => console.error('Error adding favourite:', err)
    });
  }

  onPeakRemoved(peak: SelectedPeak): void {
    this.summitFavouritesService.removeFavourite(peak.id).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (favouriteIds: number[]) => {
        this.favouriteIds = favouriteIds;
        this.loadFavouritePeaks();
      },
      error: (err: any) => console.error('Error removing favourite:', err)
    });
  }

  removeFavourite(summitId: number): void {
    this.summitFavouritesService.removeFavourite(summitId).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (favouriteIds: number[]) => {
        this.favouriteIds = favouriteIds;
        this.loadFavouritePeaks();
      },
      error: (err: any) => console.error('Error removing favourite:', err)
    });
  }

  // Navigation methods
  viewWishlistOnMap(): void {
    this.router.navigate(['/explore'], { queryParams: { filter: 'wishlist' } });
  }

  // Wishlist collapse/expand
  toggleWishlist(): void {
    this.wishlistCollapsed = !this.wishlistCollapsed;
    // Optionally save preference to localStorage
    localStorage.setItem('wishlistCollapsed', String(this.wishlistCollapsed));
  }

  createDistanceChart(): void {
    if (!this.distanceChartCanvas) return;

    const chartData = this.buildCumulativeData(this.yearActivities, 'distance');
    const goalLineData = this.buildGoalLine(this.distanceGoal);

    const ctx = this.distanceChartCanvas.nativeElement.getContext('2d');
    if (!ctx) return;

    const config: ChartConfiguration = {
      type: 'line',
      data: {
        datasets: [
          {
            label: 'Distance (km)',
            data: chartData,
            borderColor: '#fc4c02',
            backgroundColor: 'rgba(252, 76, 2, 0.1)',
            borderWidth: 3,
            fill: true,
            tension: 0.1,
            pointBackgroundColor: (context: ScriptableContext<'line'>) => {
              const point = chartData[context.dataIndex];
              return point?.hasSummit ? '#ff6b35' : '#fc4c02';
            },
            pointRadius: (context: ScriptableContext<'line'>) => {
              const point = chartData[context.dataIndex];
              return point?.hasSummit ? 6 : 3;
            },
            pointBorderWidth: 1,
            pointBorderColor: '#fff',
            order: 1,
          },
          {
            label: `Goal (${this.distanceGoal} km)`,
            data: goalLineData,
            borderColor: '#00d4aa',
            backgroundColor: 'transparent',
            borderWidth: 2,
            borderDash: [5, 5],
            fill: false,
            tension: 0,
            pointRadius: 0,
            order: 2,
          },
        ],
      },
      options: this.getChartOptions('km'),
    };

    this.distanceChart = new Chart(ctx, config);
  }

  createElevationChart(): void {
    if (!this.elevationChartCanvas) return;

    const chartData = this.buildCumulativeData(this.yearActivities, 'elevation');
    const goalLineData = this.buildGoalLine(this.elevationGoal);

    const ctx = this.elevationChartCanvas.nativeElement.getContext('2d');
    if (!ctx) return;

    const config: ChartConfiguration = {
      type: 'line',
      data: {
        datasets: [
          {
            label: 'Elevation (m)',
            data: chartData,
            borderColor: '#6b5b95',
            backgroundColor: 'rgba(107, 91, 149, 0.1)',
            borderWidth: 3,
            fill: true,
            tension: 0.1,
            pointBackgroundColor: (context: ScriptableContext<'line'>) => {
              const point = chartData[context.dataIndex];
              return point?.hasSummit ? '#8b7fb5' : '#6b5b95';
            },
            pointRadius: (context: ScriptableContext<'line'>) => {
              const point = chartData[context.dataIndex];
              return point?.hasSummit ? 6 : 3;
            },
            pointBorderWidth: 1,
            pointBorderColor: '#fff',
            order: 1,
          },
          {
            label: `Goal (${this.elevationGoal.toLocaleString()} m)`,
            data: goalLineData,
            borderColor: '#00d4aa',
            backgroundColor: 'transparent',
            borderWidth: 2,
            borderDash: [5, 5],
            fill: false,
            tension: 0,
            pointRadius: 0,
            order: 2,
          },
        ],
      },
      options: this.getChartOptions('m'),
    };

    this.elevationChart = new Chart(ctx, config);
  }

  private buildCumulativeData(activities: Activity[], type: 'distance' | 'elevation'): any[] {
    let cumulative = 0;
    return activities.map(activity => {
      if (type === 'distance') {
        cumulative += activity.distance / 1000;
      } else {
        cumulative += activity.total_elevation_gain || 0;
      }
      return {
        x: new Date(activity.start_date).getTime(),
        y: parseFloat(cumulative.toFixed(type === 'distance' ? 1 : 0)),
        hasSummit: activity.has_summit,
      };
    });
  }

  private buildGoalLine(goal: number): any[] {
    const startOfYear = new Date(this.currentYear, 0, 1);
    const endOfYear = new Date(this.currentYear, 11, 31);
    const today = new Date();

    const totalDays = Math.ceil((endOfYear.getTime() - startOfYear.getTime()) / (1000 * 60 * 60 * 24)) + 1;
    const dailyNeeded = goal / totalDays;
    const daysSoFar = Math.ceil((today.getTime() - startOfYear.getTime()) / (1000 * 60 * 60 * 24));

    return [
      { x: startOfYear.getTime(), y: 0 },
      { x: today.getTime(), y: parseFloat((dailyNeeded * daysSoFar).toFixed(1)) },
      { x: endOfYear.getTime(), y: goal },
    ];
  }

  private getChartOptions(unit: string): any {
    return {
      responsive: true,
      maintainAspectRatio: false,
      interaction: {
        intersect: false,
        mode: 'index',
      },
      scales: {
        x: {
          type: 'time',
          time: {
            unit: 'month',
            displayFormats: { month: 'MMM' },
          },
          grid: { color: 'rgba(100, 100, 100, 0.2)' },
          ticks: { color: '#666' },
        },
        y: {
          beginAtZero: true,
          grid: { color: 'rgba(100, 100, 100, 0.2)' },
          ticks: {
            color: '#666',
            callback: (value: any) => value.toLocaleString() + ' ' + unit,
          },
        },
      },
      plugins: {
        legend: { display: false },
        tooltip: {
          callbacks: {
            title: (context: any) => {
              const date = new Date(context[0].parsed.x);
              return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
            },
          } as Partial<TooltipCallbacks<'line'>>,
        },
      },
    };
  }
}
