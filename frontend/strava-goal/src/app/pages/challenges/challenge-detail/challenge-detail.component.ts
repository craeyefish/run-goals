import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subject, takeUntil } from 'rxjs';
import { ChallengeService } from 'src/app/services/challenge.service';
import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
import { ChallengePeakWithDetails, ChallengeParticipantWithUser, LeaderboardEntry } from 'src/app/models/challenge.model';

@Component({
    selector: 'challenge-detail-page',
    standalone: true,
    imports: [CommonModule],
    templateUrl: './challenge-detail.component.html',
    styleUrls: ['./challenge-detail.component.scss'],
})
export class ChallengeDetailComponent implements OnInit, OnDestroy {
    private destroy$ = new Subject<void>();
    private challengeId: number | null = null;

    constructor(
        public challengeService: ChallengeService,
        private breadcrumbService: BreadcrumbService,
        private route: ActivatedRoute,
        private router: Router
    ) { }

    challenge = this.challengeService.selectedChallenge;
    peaks = this.challengeService.challengePeaks;
    participants = this.challengeService.participants;
    summitLog = this.challengeService.summitLog;
    isLoading = this.challengeService.isLoading;
    progressPercentage = this.challengeService.progressPercentage;

    // Leaderboard for competitive challenges
    leaderboard = signal<LeaderboardEntry[]>([]);
    loadingLeaderboard = signal(false);

    activeTab: 'peaks' | 'participants' | 'activity' = 'peaks';

    ngOnInit() {
        this.route.params.pipe(takeUntil(this.destroy$)).subscribe((params) => {
            const id = parseInt(params['id'], 10);
            if (!isNaN(id)) {
                this.challengeId = id;
                this.challengeService.loadChallenge(id);
                this.challengeService.loadSummitLog(id);
                this.loadLeaderboard(id);
            }
        });

        // Update breadcrumbs when challenge loads
        this.challengeService.selectedChallenge
        // Watch for challenge changes - using effect would be cleaner but this works
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
}
