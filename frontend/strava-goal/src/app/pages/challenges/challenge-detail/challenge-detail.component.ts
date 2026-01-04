import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy, signal, computed, effect } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subject, takeUntil } from 'rxjs';
import { ChallengeService } from 'src/app/services/challenge.service';
import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
import { ProfileService } from 'src/app/services/profile.service';
import { ChallengePeakWithDetails, ChallengeParticipantWithUser, LeaderboardEntry } from 'src/app/models/challenge.model';
import { DataTableComponent, TableColumn } from 'src/app/components/shared/data-table/data-table.component';

@Component({
    selector: 'challenge-detail-page',
    standalone: true,
    imports: [CommonModule, DataTableComponent],
    templateUrl: './challenge-detail.component.html',
    styleUrls: ['./challenge-detail.component.scss'],
})
export class ChallengeDetailComponent implements OnInit, OnDestroy {
    private destroy$ = new Subject<void>();
    private challengeId: number | null = null;

    constructor(
        public challengeService: ChallengeService,
        private profileService: ProfileService,
        private breadcrumbService: BreadcrumbService,
        private route: ActivatedRoute,
        private router: Router
    ) {
        this.loadUserProfile();
    }

    challenge = this.challengeService.selectedChallenge;
    peaks = this.challengeService.challengePeaks;
    participants = this.challengeService.participants;
    summitLog = this.challengeService.summitLog;
    activities = this.challengeService.challengeActivities;
    isLoading = this.challengeService.isLoading;
    progressPercentage = this.challengeService.progressPercentage;

    // Leaderboard for competitive challenges
    leaderboard = signal<LeaderboardEntry[]>([]);
    loadingLeaderboard = signal(false);

    // Join code display
    showJoinCode = signal(false);

    // Current user
    currentUserId = signal<number | undefined>(undefined);

    activeTab: 'peaks' | 'participants' | 'activity' = 'participants';

    // Table column definitions for Peaks tab
    peakColumns: TableColumn<ChallengePeakWithDetails>[] = [
        {
            header: 'Peak Name',
            field: 'name',
            type: 'link',
            sortable: true,
            linkFn: (row) => `/explore?peakId=${row.peakId}`
        },
        {
            header: 'Elevation',
            field: 'elevation',
            type: 'number',
            sortable: true,
            formatter: (value) => `${value}m`,
            align: 'right'
        },
        {
            header: 'Region',
            field: 'region',
            type: 'text',
            sortable: true
        },
        {
            header: 'Status',
            field: 'isSummited',
            type: 'badge',
            badgeClass: (value) => value ? 'badge-success' : 'badge-default',
            formatter: (value) => value ? '✓ Summited' : 'Pending'
        }
    ];

    // Table column definitions for Participants tab (Collaborative)
    collaborativeParticipantsColumns = computed<TableColumn<ChallengeParticipantWithUser>[]>(() => {
        const challenge = this.challenge();
        if (!challenge) return [];

        return [
            {
                header: 'Member',
                field: 'userName',
                type: 'link',
                sortable: true,
                linkFn: (row) => `https://www.strava.com/athletes/${row.stravaAthleteId}`,
                linkExternal: true,
                formatter: (value, row) => value || row.stravaAthleteId.toString()
            },
            {
                header: 'Contributed',
                field: 'totalDistance',
                type: 'text',
                sortable: true,
                formatter: (value, row) => {
                    if (challenge.goalType === 'distance') return `${(row.totalDistance / 1000).toFixed(1)} km`;
                    if (challenge.goalType === 'elevation') return `${Math.round(row.totalElevation)} m`;
                    if (challenge.goalType === 'summit_count') return `${row.totalSummitCount} summits`;
                    return `${row.peaksCompleted} peaks`;
                },
                align: 'right'
            },
            {
                header: 'Joined',
                field: 'joinedAt',
                type: 'date'
            }
        ];
    });

    // Table column definitions for Leaderboard (Competitive)
    leaderboardColumns = computed<TableColumn<LeaderboardEntry>[]>(() => {
        const challenge = this.challenge();
        if (!challenge) return [];

        return [
            {
                header: 'Rank',
                field: 'rank',
                type: 'text',
                width: '80px',
                formatter: (value) => {
                    if (value === 1) return '🥇';
                    if (value === 2) return '🥈';
                    if (value === 3) return '🥉';
                    return `#${value}`;
                }
            },
            {
                header: 'Member',
                field: 'userName',
                type: 'link',
                sortable: true,
                linkFn: (row) => `https://www.strava.com/athletes/${row.stravaAthleteId}`,
                linkExternal: true,
                formatter: (value, row) => value || row.stravaAthleteId.toString()
            },
            {
                header: 'Progress',
                field: 'progress',
                type: 'progress',
                progressValue: (row) => row.progress,
                progressLabel: (row) => {
                    if (challenge.goalType === 'specific_summits') {
                        return `${row.peaksCompleted}/${row.totalPeaks}`;
                    }
                    if (challenge.goalType === 'distance') {
                        return `${(row.totalDistance / 1000).toFixed(1)}/${(challenge.targetValue! / 1000).toFixed(1)} km`;
                    }
                    if (challenge.goalType === 'elevation') {
                        return `${Math.round(row.totalElevation)}/${challenge.targetValue} m`;
                    }
                    if (challenge.goalType === 'summit_count') {
                        return `${row.totalSummitCount}/${challenge.targetSummitCount}`;
                    }
                    return `${row.progress}%`;
                }
            },
            {
                header: 'Status',
                field: 'completedAt',
                type: 'badge',
                badgeClass: (value) => value ? 'badge-success' : 'badge-default',
                formatter: (value) => value ? '✅ Complete' : 'In Progress'
            }
        ];
    });

    // Table column definitions for Activities tab
    activityColumns = computed<TableColumn<any>[]>(() => {
        const challenge = this.challenge();
        if (!challenge) return [];

        const columns: TableColumn<any>[] = [
            {
                header: 'Activity',
                field: 'name',
                type: 'link',
                linkFn: (row) => `https://www.strava.com/activities/${row.strava_activity_id}`,
                linkExternal: true
            },
            {
                header: 'User',
                field: 'userName',
                type: 'link',
                linkFn: (row) => `https://www.strava.com/athletes/${row.stravaAthleteId}`,
                linkExternal: true,
                formatter: (value, row) => value || row.stravaAthleteId.toString()
            }
        ];

        // Add goal-specific columns
        if (challenge.goalType === 'summit_count' || challenge.goalType === 'specific_summits') {
            columns.push({
                header: 'Peaks',
                field: 'peakNames',
                type: 'text'
            });
        }

        if (challenge.goalType === 'distance' || challenge.goalType === 'summit_count' || challenge.goalType === 'specific_summits') {
            columns.push({
                header: 'Distance',
                field: 'distance',
                type: 'number',
                formatter: (value) => `${(value / 1000).toFixed(1)} km`,
                align: 'right'
            });
        }

        if (challenge.goalType === 'elevation' || challenge.goalType === 'summit_count' || challenge.goalType === 'specific_summits') {
            columns.push({
                header: 'Elevation',
                field: 'total_elevation_gain',
                type: 'number',
                formatter: (value) => `${Math.round(value)} m`,
                align: 'right'
            });
        }

        columns.push({
            header: 'Date',
            field: 'start_date',
            type: 'date',
            sortable: true
        });

        return columns;
    });

    loadUserProfile() {
        this.profileService.getUserProfile().subscribe({
            next: (profile) => {
                this.currentUserId.set(profile.id);
            },
            error: (err) => {
                console.error('Failed to load user profile', err);
            }
        });
    }

    ngOnInit() {
        this.route.params.pipe(takeUntil(this.destroy$)).subscribe((params) => {
            const id = parseInt(params['id'], 10);
            if (!isNaN(id)) {
                this.challengeId = id;
                this.challengeService.loadChallenge(id);
                this.challengeService.loadSummitLog(id);
                this.challengeService.loadChallengeActivities(id);
                this.loadLeaderboard(id);
            }
        });

        // Watch for challenge to load and update breadcrumb
        effect(() => {
            const challenge = this.challengeService.selectedChallenge();
            if (challenge) {
                this.breadcrumbService.addItem({
                    label: challenge.name,
                });
            }
        });
    }

    loadLeaderboard(challengeId: number): void {
        this.loadingLeaderboard.set(true);
        this.challengeService.getLeaderboard(challengeId).pipe(
            takeUntil(this.destroy$)
        ).subscribe({
            next: (entries) => {
                this.leaderboard.set(entries);
                this.loadingLeaderboard.set(false);
            },
            error: (err) => {
                console.error('Failed to load leaderboard:', err);
                this.loadingLeaderboard.set(false);
            }
        });
    }

    ngOnDestroy() {
        this.destroy$.next();
        this.destroy$.complete();
        this.challengeService.clearSelection();
    }

    setTab(tab: 'peaks' | 'participants' | 'activity') {
        this.activeTab = tab;
    }

    get completedPeaks(): ChallengePeakWithDetails[] {
        return this.peaks().filter(p => p.isSummited);
    }

    get remainingPeaks(): ChallengePeakWithDetails[] {
        return this.peaks().filter(p => !p.isSummited);
    }

    get sortedParticipants(): ChallengeParticipantWithUser[] {
        return [...this.participants()].sort((a, b) => b.peaksCompleted - a.peaksCompleted);
    }

    onJoinChallenge() {
        if (!this.challengeId) return;
        this.challengeService.joinChallenge(this.challengeId).subscribe({
            next: () => {
                this.challengeService.loadChallenge(this.challengeId!);
            },
            error: (err) => console.error('Failed to join', err),
        });
    }

    onLeaveChallenge() {
        if (!this.challengeId) return;
        if (!confirm('Are you sure you want to leave this challenge?')) return;

        this.challengeService.leaveChallenge(this.challengeId).subscribe({
            next: () => {
                this.router.navigate(['/challenges']);
            },
            error: (err) => console.error('Failed to leave', err),
        });
    }

    onDeleteChallenge() {
        if (!this.challengeId) return;
        if (!confirm('Are you sure you want to delete this challenge? This cannot be undone.')) return;

        this.challengeService.deleteChallenge(this.challengeId).subscribe({
            next: () => {
                this.router.navigate(['/challenges']);
            },
            error: (err) => console.error('Failed to delete', err),
        });
    }

    onLockChallenge() {
        if (!this.challengeId) return;
        if (!confirm('Are you sure you want to lock this challenge? Once locked, it cannot be edited.')) return;

        this.challengeService.lockChallenge(this.challengeId).subscribe({
            next: () => {
                this.challengeService.loadChallenge(this.challengeId!);
            },
            error: (err) => console.error('Failed to lock challenge', err),
        });
    }

    copyJoinCode() {
        const challenge = this.challenge();
        if (!challenge) return;

        navigator.clipboard.writeText(challenge.joinCode).then(() => {
            // Could add a toast notification here
            console.log('Join code copied to clipboard');
        });
    }

    toggleJoinCode() {
        this.showJoinCode.update(v => !v);
    }

    goBack() {
        this.router.navigate(['/challenges']);
    }

    formatDate(dateStr: string | undefined): string {
        if (!dateStr) return '';
        return new Date(dateStr).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
        });
    }

    getParticipantProgressPercentage(participant: ChallengeParticipantWithUser): number {
        const challenge = this.challenge();
        if (!challenge) return 0;

        switch (challenge.goalType) {
            case 'specific_summits':
                if (participant.totalPeaks === 0) return 0;
                return Math.round((participant.peaksCompleted / participant.totalPeaks) * 100);

            case 'distance':
                if (!challenge.targetValue || challenge.targetValue === 0) return 0;
                return Math.min(100, Math.round((participant.totalDistance / challenge.targetValue) * 100));

            case 'elevation':
                if (!challenge.targetValue || challenge.targetValue === 0) return 0;
                return Math.min(100, Math.round((participant.totalElevation / challenge.targetValue) * 100));

            case 'summit_count':
                if (!challenge.targetSummitCount || challenge.targetSummitCount === 0) return 0;
                return Math.min(100, Math.round((participant.totalSummitCount / challenge.targetSummitCount) * 100));

            default:
                return 0;
        }
    }
}
