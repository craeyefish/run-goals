import { CommonModule } from '@angular/common';
import { Component, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { ChallengeService } from 'src/app/services/challenge.service';
import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
import { ProfileService } from 'src/app/services/profile.service';
import { Challenge, ChallengeWithProgress, CreateChallengeRequest } from 'src/app/models/challenge.model';
import { ChallengeCardComponent } from 'src/app/components/challenges/challenge-card/challenge-card.component';
import { ChallengeCreateFormComponent } from 'src/app/components/challenges/challenge-create-form/challenge-create-form.component';

@Component({
    selector: 'challenge-list-page',
    standalone: true,
    imports: [
        CommonModule,
        FormsModule,
        ChallengeCardComponent,
        ChallengeCreateFormComponent,
    ],
    templateUrl: './challenge-list.component.html',
    styleUrls: ['./challenge-list.component.scss'],
})
export class ChallengeListComponent implements OnInit {
    constructor(
        public challengeService: ChallengeService,
        private profileService: ProfileService,
        private breadcrumbService: BreadcrumbService,
        private router: Router
    ) {
        this.challengeService.loadUserChallenges();
        this.challengeService.loadFeaturedChallenges();
        this.loadUserProfile();
    }

    challenges = this.challengeService.challenges;
    featuredChallenges = this.challengeService.featuredChallenges;
    isLoading = this.challengeService.isLoading;
    currentUserId = signal<number | undefined>(undefined);

    showCreateForm = signal(false);
    showJoinByCodeForm = signal(false);
    joinCodeInput = signal('');
    joinCodeError = signal('');
    activeTab = signal<'joined' | 'created' | 'discover'>('joined');

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
        this.breadcrumbService.setItems([
            {
                label: 'Challenges',
                routerLink: '/challenges',
                callback: () => {
                    this.challengeService.clearSelection();
                    this.router.navigate(['/challenges']);
                },
            },
        ]);
    }

    openCreateForm() {
        this.showCreateForm.set(true);
    }

    closeCreateForm() {
        this.showCreateForm.set(false);
    }

    onChallengeCreated(request: CreateChallengeRequest) {
        this.challengeService.createChallenge(request).subscribe({
            next: (response) => {
                this.showCreateForm.set(false);
                this.challengeService.loadUserChallenges();
                this.router.navigate(['/challenges', response.id]);
            },
            error: (err) => {
                console.error('Failed to create challenge', err);
            },
        });
    }

    onChallengeClick(challenge: Challenge | ChallengeWithProgress) {
        this.router.navigate(['/challenges', challenge.id]);
    }

    setTab(tab: 'joined' | 'created' | 'discover') {
        this.activeTab.set(tab);
        if (tab === 'discover') {
            this.challengeService.loadPublicChallenges();
        }
    }

    openJoinByCodeForm() {
        this.showJoinByCodeForm.set(true);
        this.joinCodeInput.set('');
        this.joinCodeError.set('');
    }

    closeJoinByCodeForm() {
        this.showJoinByCodeForm.set(false);
        this.joinCodeInput.set('');
        this.joinCodeError.set('');
    }

    joinByCode() {
        const code = this.joinCodeInput().trim().toUpperCase();
        if (!code) {
            this.joinCodeError.set('Please enter a join code');
            return;
        }
        if (code.length !== 6) {
            this.joinCodeError.set('Join code must be 6 characters');
            return;
        }

        this.challengeService.joinChallengeByCode(code).subscribe({
            next: (challenge) => {
                this.closeJoinByCodeForm();
                this.challengeService.loadUserChallenges();
                this.router.navigate(['/challenges', challenge.id]);
            },
            error: (err) => {
                if (err.status === 404) {
                    this.joinCodeError.set('Invalid join code');
                } else if (err.status === 409) {
                    this.joinCodeError.set('You have already joined this challenge');
                } else {
                    this.joinCodeError.set('Failed to join challenge');
                }
            },
        });
    }

    get myChallenges(): ChallengeWithProgress[] {
        return this.challenges() || [];
    }

    get joinedChallenges(): ChallengeWithProgress[] {
        // Challenges I've joined (not created by me, and currently joined)
        const userId = this.currentUserId();
        if (!userId) return [];
        return this.myChallenges.filter(c => c.isJoined && c.createdByUserId !== userId);
    }

    get createdChallenges(): ChallengeWithProgress[] {
        // Challenges I've created
        const userId = this.currentUserId();
        if (!userId) return [];
        return this.myChallenges.filter(c => c.createdByUserId === userId);
    }
}
