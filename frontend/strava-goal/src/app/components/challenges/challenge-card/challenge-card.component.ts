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
        if (!c || c.totalPeaks === 0) return 0;
        return Math.round((c.completedPeaks / c.totalPeaks) * 100);
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
}
