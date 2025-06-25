import {
  Component,
  OnInit,
  OnDestroy,
  ViewChild,
  ElementRef,
  AfterViewInit,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivityService, Activity } from '../../services/activity.service';
import { PeakSummitService } from '../../services/peak-summit.service';
import { GroupService, Group, Goal } from '../../services/groups.service';
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

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent implements OnInit, OnDestroy, AfterViewInit {
  @ViewChild('progressChart') chartCanvas!: ElementRef<HTMLCanvasElement>;

  loading = true;
  error: string | null = null;
  currentYear = new Date().getFullYear();

  // Stats
  totalDistance = 0;
  totalElevation = 0; // New stat
  totalActivities = 0;
  totalSummits = 0;
  avgDistance = 0;

  // Personal goals (keeping the old ones for now)
  distanceGoal = 1000; // 1000km goal
  summitGoal = 20; // 20 peaks goal
  distanceGoalPercentage = 0;
  summitGoalPercentage = 0;

  // Group goals
  groupGoals: GroupGoalDisplay[] = [];

  // Chart
  chart: Chart | null = null;

  private destroy$ = new Subject<void>();

  constructor(
    private activityService: ActivityService,
    private peakSummitService: PeakSummitService,
    private groupService: GroupService
  ) {}

  ngOnInit(): void {
    this.loadDashboardData();
  }

  ngAfterViewInit(): void {
    // Chart will be created after data loads
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
    if (this.chart) {
      this.chart.destroy();
    }
  }

  loadDashboardData(): void {
    console.log('Loading dashboard data...');
    this.activityService.loadActivities();
    this.groupService.loadGroups();

    combineLatest([
      this.activityService.activities$.pipe(
        filter((activities) => activities !== null)
      ),
      this.peakSummitService.getPeakSummaries(),
    ])
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: ([activities, peakSummaries]) => {
          console.log('Dashboard data loaded:', {
            activities: activities?.length,
            peakSummaries: peakSummaries?.length,
          });

          // Ensure we have valid data before calculating stats
          const validActivities = activities || [];
          const validPeakSummaries = peakSummaries || [];

          this.calculateStats(validActivities, validPeakSummaries);
          this.loadGroupGoals();

          // Ensure the view has been initialized before creating chart
          setTimeout(() => {
            this.createProgressChart(validActivities);
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

  loadGroupGoals(): void {
    const groups = this.groupService.groups();
    this.groupGoals = [];

    if (groups.length === 0) {
      return;
    }

    // For each group, load its goals
    groups.forEach((group) => {
      this.groupService.getGroupGoals(group.id).subscribe({
        next: (response) => {
          if (response.goals.length > 0) {
            this.groupGoals.push({
              groupName: group.name,
              goals: response.goals.map((goal) => ({
                ...goal,
                progressPercentage: this.calculateGoalProgress(goal),
              })),
            });
          }
        },
        error: (err) => {
          console.error(`Failed to load goals for group ${group.name}:`, err);
        },
      });
    });
  }

  calculateGoalProgress(goal: Goal): number {
    // This is a simplified calculation - you might want to implement
    // more sophisticated logic based on goal type and current user stats
    const today = new Date();
    const startDate = new Date(goal.start_date || today);
    const endDate = new Date(goal.end_date || today);

    if (startDate > today) return 0; // Goal hasn't started yet
    if (endDate < today) return 100; // Goal has ended

    // Calculate time-based progress as a fallback
    const totalTime = endDate.getTime() - startDate.getTime();
    const elapsedTime = today.getTime() - startDate.getTime();
    return Math.min(100, Math.max(0, (elapsedTime / totalTime) * 100));
  }

  calculateStats(activities: Activity[], peakSummaries: any[]): void {
    const currentYear = new Date().getFullYear();

    // Filter activities for current year
    const yearActivities = activities.filter((activity) => {
      const activityYear = new Date(activity.start_date).getFullYear();
      return activityYear === currentYear;
    });

    // Calculate basic stats
    this.totalActivities = yearActivities.length;
    this.totalDistance = parseFloat(
      yearActivities
        .reduce((total, activity) => {
          return total + activity.distance / 1000;
        }, 0)
        .toFixed(1)
    );

    // Calculate total elevation gain
    this.totalElevation = parseFloat(
      yearActivities
        .reduce((total, activity) => {
          return total + (activity.total_elevation_gain || 0);
        }, 0)
        .toFixed(0)
    );

    this.avgDistance =
      this.totalActivities > 0
        ? parseFloat((this.totalDistance / this.totalActivities).toFixed(1))
        : 0;

    // Calculate summits - add null check
    let summitCount = 0;
    if (peakSummaries && Array.isArray(peakSummaries)) {
      peakSummaries.forEach((peak) => {
        if (peak.activities && Array.isArray(peak.activities)) {
          peak.activities.forEach((activity: any) => {
            const activityYear = new Date(activity.start_date).getFullYear();
            if (activityYear === currentYear) {
              summitCount++;
            }
          });
        }
      });
    }
    this.totalSummits = summitCount;

    // Calculate personal goal percentages - fixed to 2 decimal places
    this.distanceGoalPercentage = parseFloat(
      Math.min(100, (this.totalDistance / this.distanceGoal) * 100).toFixed(2)
    );
    this.summitGoalPercentage = parseFloat(
      Math.min(100, (this.totalSummits / this.summitGoal) * 100).toFixed(2)
    );
  }

  createProgressChart(activities: Activity[]): void {
    if (!this.chartCanvas) return;

    const currentYear = new Date().getFullYear();
    const yearActivities = activities
      .filter(
        (activity) =>
          new Date(activity.start_date).getFullYear() === currentYear
      )
      .sort(
        (a, b) =>
          new Date(a.start_date).getTime() - new Date(b.start_date).getTime()
      );

    console.log('Year activities for chart:', yearActivities.length);

    // Create cumulative distance data
    let cumulativeDistance = 0;
    const chartData = yearActivities.map((activity) => {
      cumulativeDistance += activity.distance / 1000;
      return {
        x: new Date(activity.start_date).getTime(),
        y: parseFloat(cumulativeDistance.toFixed(1)),
        hasSummit: activity.has_summit,
      };
    });

    // Create goal line data (daily progress needed to reach annual goal)
    const startOfYear = new Date(currentYear, 0, 1); // January 1st
    const endOfYear = new Date(currentYear, 11, 31); // December 31st
    const today = new Date();

    // Calculate daily distance needed
    const totalDaysInYear =
      Math.ceil(
        (endOfYear.getTime() - startOfYear.getTime()) / (1000 * 60 * 60 * 24)
      ) + 1;
    const dailyDistanceNeeded = this.distanceGoal / totalDaysInYear;

    // Create goal line points
    const goalLineData = [];

    // Start point
    goalLineData.push({
      x: startOfYear.getTime(),
      y: 0,
    });

    // If we have activities, create points that align with activity dates
    if (yearActivities.length > 0) {
      yearActivities.forEach((activity) => {
        const activityDate = new Date(activity.start_date);
        const daysSinceStart = Math.ceil(
          (activityDate.getTime() - startOfYear.getTime()) /
            (1000 * 60 * 60 * 24)
        );
        const expectedDistance = dailyDistanceNeeded * daysSinceStart;

        goalLineData.push({
          x: activityDate.getTime(),
          y: parseFloat(expectedDistance.toFixed(1)),
        });
      });
    }

    // Add current date point
    const daysSinceStartToday = Math.ceil(
      (today.getTime() - startOfYear.getTime()) / (1000 * 60 * 60 * 24)
    );
    const expectedDistanceToday = dailyDistanceNeeded * daysSinceStartToday;
    goalLineData.push({
      x: today.getTime(),
      y: parseFloat(expectedDistanceToday.toFixed(1)),
    });

    // End point
    goalLineData.push({
      x: endOfYear.getTime(),
      y: this.distanceGoal,
    });

    console.log('Chart data:', chartData);
    console.log('Goal line data:', goalLineData);

    const ctx = this.chartCanvas.nativeElement.getContext('2d');
    if (!ctx) return;

    const config: ChartConfiguration = {
      type: 'line',
      data: {
        datasets: [
          {
            label: 'Cumulative Distance (km)',
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
              return point?.hasSummit ? 8 : 4;
            },
            pointBorderWidth: 2,
            pointBorderColor: '#fff',
            order: 1,
          },
          {
            label: `Goal Pace (${this.distanceGoal}km/year)`,
            data: goalLineData,
            borderColor: '#00d4aa', // Turquoise/teal color
            backgroundColor: 'transparent',
            borderWidth: 2,
            borderDash: [5, 5], // Dashed line
            fill: false,
            tension: 0,
            pointRadius: 0, // No points on goal line
            pointHoverRadius: 0,
            order: 2,
          },
        ],
      },
      options: {
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
              displayFormats: {
                month: 'MMM yyyy',
              },
            },
            grid: {
              color: 'rgba(100, 100, 100, 0.2)',
            },
            ticks: {
              color: '#666',
            },
          },
          y: {
            beginAtZero: true,
            grid: {
              color: 'rgba(100, 100, 100, 0.2)',
            },
            ticks: {
              color: '#666',
              callback: function (value) {
                return value + ' km';
              },
            },
          },
        },
        plugins: {
          legend: {
            labels: {
              color: '#333',
            },
          },
          tooltip: {
            callbacks: {
              title: (context: any) => {
                // Format the date nicely
                const date = new Date(context[0].parsed.x);
                return date.toLocaleDateString('en-US', {
                  year: 'numeric',
                  month: 'short',
                  day: 'numeric',
                });
              },
              afterLabel: (context: any) => {
                if (context.datasetIndex === 0) {
                  // Only for actual distance line
                  const point = chartData[context.dataIndex];
                  return point?.hasSummit ? '🏔️ Peak Summit!' : '';
                }
                return '';
              },
            } as Partial<TooltipCallbacks<'line'>>,
          },
        },
      },
    };

    this.chart = new Chart(ctx, config);
  }
}

// Interface for displaying group goals
interface GroupGoalDisplay {
  groupName: string;
  goals: (Goal & { progressPercentage: number })[];
}
