import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { BreadcrumbComponent } from 'src/app/components/breadcrumb/breadcrumb.component';
import { ChallengeService } from 'src/app/services/challenge.service';

@Component({
    selector: 'app-challenges',
    standalone: true,
    imports: [
        CommonModule,
        RouterOutlet,
        BreadcrumbComponent,
    ],
    templateUrl: './challenges.component.html',
    styleUrls: ['./challenges.component.scss'],
})
export class ChallengesComponent {
    constructor(private challengeService: ChallengeService) { }

    challenges = this.challengeService.challenges;
    selectedChallenge = this.challengeService.selectedChallenge;
}
