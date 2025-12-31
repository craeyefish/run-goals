import { CommonModule } from '@angular/common';
import { Component, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { ChallengeService } from 'src/app/services/challenge.service';
import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
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
        private breadcrumbService: BreadcrumbService,
        private router: Router
    ) {
        this.challengeService.loadUserChallenges();
        this.challengeService.loadFeaturedChallenges();
    }

    challenges = this.challengeService.challenges;
    featuredChallenges = this.challengeService.featuredChallenges;
    isLoading = this.challengeService.isLoading;

    showCreateForm = signal(false);
    activeTab = signal<'my' | 'discover'>('my');

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

    setTab(tab: 'my' | 'discover') {
        this.activeTab.set(tab);
        if (tab === 'discover') {
            this.challengeService.loadPublicChallenges();
        }
    }

    get myChallenges(): ChallengeWithProgress[] {
        return this.challenges() || [];
    }

    get activeChallenges(): ChallengeWithProgress[] {
        return this.myChallenges.filter(c => !c.isCompleted);
    }

    get completedChallenges(): ChallengeWithProgress[] {
        return this.myChallenges.filter(c => c.isCompleted);
    }
}
