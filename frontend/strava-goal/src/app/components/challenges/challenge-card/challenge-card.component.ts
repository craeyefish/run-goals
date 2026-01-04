import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { Challenge, ChallengeWithProgress } from 'src/app/models/challenge.model';
import { ChallengeService } from 'src/app/services/challenge.service';

@Component({
    selector: 'challenge-card',
    standalone: true,
    imports: [CommonModule],
    templateUrl: './challenge-card.component.html',
    styleUrls: ['./challenge-card.component.scss'],
})
export class ChallengeCardComponent {
    @Input() challenge!: Challenge | ChallengeWithProgress;
    @Input() showJoinButton = false;
    @Input() showManageOptions = false;

    constructor(public challengeService: ChallengeService) { }

    get progressChallenge(): ChallengeWithProgress | null {
        return this.isProgressChallenge(this.challenge) ? this.challenge : null;
    }

    get isJoined(): boolean {
        // Check if it's a ChallengeWithProgress with isJoined flag
        if (this.isProgressChallenge(this.challenge)) {
            return this.challenge.isJoined;
        }
        // For public challenges, check if it exists in user's challenges
        return this.challengeService.challenges().some(c => c.id === this.challenge.id);
    }

    get progressPercentage(): number {
        const c = this.progressChallenge;
        if (!c) return 0;

        switch (c.goalType) {
            case 'specific_summits':
                if (c.totalPeaks === 0) return 0;
                return Math.round((c.completedPeaks / c.totalPeaks) * 100);

            case 'distance':
                if (!c.targetValue || c.targetValue === 0) return 0;
                return Math.min(100, Math.round((c.currentDistance / c.targetValue) * 100));

            case 'elevation':
                if (!c.targetValue || c.targetValue === 0) return 0;
                return Math.min(100, Math.round((c.currentElevation / c.targetValue) * 100));

            case 'summit_count':
                if (!c.targetSummitCount || c.targetSummitCount === 0) return 0;
                return Math.min(100, Math.round((c.currentSummitCount / c.targetSummitCount) * 100));

            default:
                return 0;
        }
    }

    get hasDeadline(): boolean {
        return !!this.challenge.deadline;
    }

    get daysRemaining(): number | null {
        if (!this.challenge.deadline) return null;
        const deadline = new Date(this.challenge.deadline);
        const today = new Date();
        const diff = deadline.getTime() - today.getTime();
        return Math.ceil(diff / (1000 * 60 * 60 * 24));
    }

    get isOverdue(): boolean {
        const days = this.daysRemaining;
        return days !== null && days < 0;
    }

    private isProgressChallenge(c: Challenge | ChallengeWithProgress): c is ChallengeWithProgress {
        return 'totalPeaks' in c;
    }

    onJoinClick(event: Event) {
        event.stopPropagation();
        this.challengeService.joinChallenge(this.challenge.id).subscribe({
            next: () => {
                this.challengeService.loadUserChallenges();
                this.challengeService.loadPublicChallenges(); // Refresh public list too
            },
            error: (err) => {
                console.error('Failed to join challenge', err);
            },
        });
    }

    copyJoinCode(event: Event) {
        event.stopPropagation();
        navigator.clipboard.writeText(this.challenge.joinCode).then(() => {
            // Could add a toast notification here
            console.log('Join code copied to clipboard');
        });
    }
}
