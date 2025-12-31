import { Component, OnInit, OnDestroy, Input, Output, EventEmitter, ViewChild, ElementRef, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { PersonalGoalsService, PersonalYearlyGoal } from '../../services/personal-goals.service';
import { Activity } from '../../services/activity.service';
import { Chart, registerables, ChartConfiguration } from 'chart.js';
import 'chartjs-adapter-date-fns';

Chart.register(...registerables);

interface YearlyProgress {
    distance: { current: number; goal: number; percentage: number };
    elevation: { current: number; goal: number; percentage: number };
    summits: { current: number; goal: number; percentage: number };
}

export interface RecentSummit {
    peakId: number;
    peakName: string;
    summitedAt: Date;
    activityId: number;
}

@Component({
    selector: 'app-yearly-goals',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './yearly-goals.component.html',
    styleUrls: ['./yearly-goals.component.scss']
})
export class YearlyGoalsComponent implements OnInit, OnDestroy, AfterViewInit {
    @ViewChild('distanceChart') distanceChartCanvas!: ElementRef<HTMLCanvasElement>;
    @ViewChild('elevationChart') elevationChartCanvas!: ElementRef<HTMLCanvasElement>;

    @Input() currentDistance = 0;
    @Input() currentElevation = 0;
    @Input() currentSummits = 0;
    @Input() yearActivities: Activity[] = [];
    @Input() recentSummits: RecentSummit[] = [];
    @Output() openPeakPicker = new EventEmitter<void>();

    private destroy$ = new Subject<void>();

    currentYear = new Date().getFullYear();
    selectedYear = this.currentYear;

    // Current goals
    goals: PersonalYearlyGoal | null = null;
    progress: YearlyProgress = {
        distance: { current: 0, goal: 1000, percentage: 0 },
        elevation: { current: 0, goal: 50000, percentage: 0 },
        summits: { current: 0, goal: 20, percentage: 0 }
    };

    // Expanded sections
    expandedSection: 'distance' | 'elevation' | 'summits' | null = null;

    // Charts - use 'any' type to avoid complex Chart.js generics
    distanceChart: Chart | null = null;
    elevationChart: Chart | null = null;

    // Historical goals
    historicalGoals: PersonalYearlyGoal[] = [];
    showHistory = false;

    // Modal states
    showGoalEditor = false;
    editingGoalType: 'distance' | 'elevation' | 'summits' = 'distance';
    editingGoalValue = 0;

    // Quick set modal
    showQuickSetModal = false;
    quickSetGoals = {
        distance: 1000,
        elevation: 50000,
        summits: 20
    };

    constructor(
        private personalGoalsService: PersonalGoalsService
    ) { }

    ngOnInit(): void {
        this.loadGoals();
        this.loadHistoricalGoals();
    }

    ngAfterViewInit(): void {
        // Charts created when section expands
    }

    ngOnDestroy(): void {
        this.destroy$.next();
        this.destroy$.complete();
        this.destroyCharts();
    }

    ngOnChanges(): void {
        this.updateProgress();
        // Recreate charts if expanded and data changed
        if (this.expandedSection === 'distance' && this.distanceChartCanvas) {
            setTimeout(() => this.createDistanceChart(), 50);
        }
        if (this.expandedSection === 'elevation' && this.elevationChartCanvas) {
            setTimeout(() => this.createElevationChart(), 50);
        }
    }

    destroyCharts(): void {
        if (this.distanceChart) {
            this.distanceChart.destroy();
            this.distanceChart = null;
        }
        if (this.elevationChart) {
            this.elevationChart.destroy();
            this.elevationChart = null;
        }
    }

    loadGoals(): void {
        this.personalGoalsService.getGoals(this.selectedYear).pipe(
            takeUntil(this.destroy$)
        ).subscribe({
            next: (goals) => {
                this.goals = goals;
                this.updateProgress();
            },
            error: (err) => console.error('Error loading goals:', err)
        });
    }

    loadHistoricalGoals(): void {
        this.personalGoalsService.getAllGoals().pipe(
            takeUntil(this.destroy$)
        ).subscribe({
            next: (goals) => {
                this.historicalGoals = goals.filter(g => g.year !== this.currentYear);
            },
            error: (err) => console.error('Error loading historical goals:', err)
        });
    }

    updateProgress(): void {
        const distanceGoal = this.goals?.distance_goal || 1000;
        const elevationGoal = this.goals?.elevation_goal || 50000;
        const summitGoal = this.goals?.summit_goal || 20;

        this.progress = {
            distance: {
                current: this.currentDistance,
                goal: distanceGoal,
                percentage: Math.min(100, (this.currentDistance / distanceGoal) * 100)
            },
            elevation: {
                current: this.currentElevation,
                goal: elevationGoal,
                percentage: Math.min(100, (this.currentElevation / elevationGoal) * 100)
            },
            summits: {
                current: this.currentSummits,
                goal: summitGoal,
                percentage: Math.min(100, (this.currentSummits / summitGoal) * 100)
            }
        };
    }

    // Toggle expanded section
    toggleSection(section: 'distance' | 'elevation' | 'summits', event: MouseEvent): void {
        // Prevent toggle when clicking edit button
        if ((event.target as HTMLElement).closest('.edit-btn')) {
            return;
        }

        if (this.expandedSection === section) {
            this.expandedSection = null;
            this.destroyCharts();
        } else {
            this.expandedSection = section;
            // Create chart after DOM updates
            setTimeout(() => {
                if (section === 'distance') {
                    this.createDistanceChart();
                } else if (section === 'elevation') {
                    this.createElevationChart();
                }
            }, 100);
        }
    }

    // Get milestone status
    getMilestone(percentage: number): string {
        if (percentage >= 100) return '🏆';
        if (percentage >= 75) return '🔥';
        if (percentage >= 50) return '💪';
        if (percentage >= 25) return '🚀';
        return '🎯';
    }

    // Open goal editor for specific type
    openGoalEditor(type: 'distance' | 'elevation' | 'summits', event: MouseEvent): void {
        event.stopPropagation();
        this.editingGoalType = type;
        this.editingGoalValue = type === 'distance' ? this.progress.distance.goal :
            type === 'elevation' ? this.progress.elevation.goal :
                this.progress.summits.goal;
        this.showGoalEditor = true;
    }

    closeGoalEditor(): void {
        this.showGoalEditor = false;
    }

    saveGoal(): void {
        if (!this.goals) {
            this.goals = {
                year: this.currentYear,
                distance_goal: this.progress.distance.goal,
                elevation_goal: this.progress.elevation.goal,
                summit_goal: this.progress.summits.goal,
                target_summits: []
            };
        }

        if (this.editingGoalType === 'distance') {
            this.goals.distance_goal = this.editingGoalValue;
        } else if (this.editingGoalType === 'elevation') {
            this.goals.elevation_goal = this.editingGoalValue;
        } else {
            this.goals.summit_goal = this.editingGoalValue;
        }

        this.personalGoalsService.saveGoals(this.goals).pipe(
            takeUntil(this.destroy$)
        ).subscribe({
            next: (saved) => {
                this.goals = saved;
                this.updateProgress();
                this.closeGoalEditor();
                // Recreate chart if section is expanded
                if (this.expandedSection === 'distance') {
                    setTimeout(() => this.createDistanceChart(), 50);
                } else if (this.expandedSection === 'elevation') {
                    setTimeout(() => this.createElevationChart(), 50);
                }
            },
            error: (err) => console.error('Error saving goal:', err)
        });
    }

    // Quick set all goals at once
    openQuickSetModal(): void {
        this.quickSetGoals = {
            distance: this.goals?.distance_goal || 1000,
            elevation: this.goals?.elevation_goal || 50000,
            summits: this.goals?.summit_goal || 20
        };
        this.showQuickSetModal = true;
    }

    closeQuickSetModal(): void {
        this.showQuickSetModal = false;
    }

    saveAllGoals(): void {
        const goalData: PersonalYearlyGoal = {
            year: this.currentYear,
            distance_goal: this.quickSetGoals.distance,
            elevation_goal: this.quickSetGoals.elevation,
            summit_goal: this.quickSetGoals.summits,
            target_summits: this.goals?.target_summits || []
        };

        this.personalGoalsService.saveGoals(goalData).pipe(
            takeUntil(this.destroy$)
        ).subscribe({
            next: (saved) => {
                this.goals = saved;
                this.updateProgress();
                this.closeQuickSetModal();
            },
            error: (err) => console.error('Error saving goals:', err)
        });
    }

    toggleHistory(): void {
        this.showHistory = !this.showHistory;
        if (this.showHistory) {
            this.expandedSection = null;
            this.destroyCharts();
        }
    }

    onOverlayClick(event: MouseEvent): void {
        if ((event.target as HTMLElement).classList.contains('modal-overlay')) {
            this.showGoalEditor = false;
            this.showQuickSetModal = false;
        }
    }

    // ==================== CHARTS ====================

    createDistanceChart(): void {
        if (!this.distanceChartCanvas?.nativeElement) return;

        if (this.distanceChart) {
            this.distanceChart.destroy();
        }

        const ctx = this.distanceChartCanvas.nativeElement.getContext('2d');
        if (!ctx) return;

        const goal = this.progress.distance.goal;
        const cumulativeData = this.getCumulativeData('distance');
        const goalLine = this.getGoalLine(goal);

        this.distanceChart = new Chart(ctx, {
            type: 'line',
            data: {
                datasets: [
                    {
                        label: 'Actual',
                        data: cumulativeData,
                        borderColor: '#00d4aa',
                        backgroundColor: 'rgba(0, 212, 170, 0.1)',
                        fill: true,
                        tension: 0.3,
                        pointRadius: 0,
                        pointHoverRadius: 4,
                    },
                    {
                        label: 'Goal pace',
                        data: goalLine,
                        borderColor: 'rgba(0, 0, 0, 0.3)',
                        borderDash: [5, 5],
                        fill: false,
                        tension: 0,
                        pointRadius: 0,
                    }
                ]
            },
            options: this.getChartOptions('km')
        });
    }

    createElevationChart(): void {
        if (!this.elevationChartCanvas?.nativeElement) return;

        if (this.elevationChart) {
            this.elevationChart.destroy();
        }

        const ctx = this.elevationChartCanvas.nativeElement.getContext('2d');
        if (!ctx) return;

        const goal = this.progress.elevation.goal;
        const cumulativeData = this.getCumulativeData('elevation');
        const goalLine = this.getGoalLine(goal);

        this.elevationChart = new Chart(ctx, {
            type: 'line',
            data: {
                datasets: [
                    {
                        label: 'Actual',
                        data: cumulativeData,
                        borderColor: '#ff6b6b',
                        backgroundColor: 'rgba(255, 107, 107, 0.1)',
                        fill: true,
                        tension: 0.3,
                        pointRadius: 0,
                        pointHoverRadius: 4,
                    },
                    {
                        label: 'Goal pace',
                        data: goalLine,
                        borderColor: 'rgba(0, 0, 0, 0.3)',
                        borderDash: [5, 5],
                        fill: false,
                        tension: 0,
                        pointRadius: 0,
                    }
                ]
            },
            options: this.getChartOptions('m')
        });
    }

    private getCumulativeData(type: 'distance' | 'elevation'): { x: number; y: number }[] {
        if (!this.yearActivities?.length) return [];

        let cumulative = 0;
        return this.yearActivities.map(activity => {
            if (type === 'distance') {
                cumulative += activity.distance / 1000;
            } else {
                cumulative += activity.total_elevation_gain || 0;
            }
            return {
                x: new Date(activity.start_date).getTime(),
                y: Math.round(cumulative)
            };
        });
    }

    private getGoalLine(goal: number): { x: number; y: number }[] {
        const startOfYear = new Date(this.currentYear, 0, 1).getTime();
        const endOfYear = new Date(this.currentYear, 11, 31).getTime();
        return [
            { x: startOfYear, y: 0 },
            { x: endOfYear, y: goal }
        ];
    }

    private getChartOptions(unit: string): ChartConfiguration['options'] {
        return {
            responsive: true,
            maintainAspectRatio: false,
            interaction: {
                intersect: false,
                mode: 'index'
            },
            plugins: {
                legend: {
                    display: true,
                    position: 'top',
                    labels: {
                        boxWidth: 12,
                        padding: 10,
                        font: { size: 11 }
                    }
                },
                tooltip: {
                    callbacks: {
                        label: (context) => `${context.dataset.label}: ${context.parsed.y.toLocaleString()} ${unit}`
                    }
                }
            },
            scales: {
                x: {
                    type: 'time',
                    time: {
                        unit: 'month',
                        displayFormats: { month: 'MMM' }
                    },
                    min: new Date(this.currentYear, 0, 1).getTime(),
                    max: new Date(this.currentYear, 11, 31).getTime(),
                    grid: { display: false },
                    ticks: { font: { size: 10 } }
                },
                y: {
                    beginAtZero: true,
                    grid: { color: 'rgba(0, 0, 0, 0.05)' },
                    ticks: {
                        font: { size: 10 },
                        callback: (value) => value.toLocaleString()
                    }
                }
            }
        };
    }
}
