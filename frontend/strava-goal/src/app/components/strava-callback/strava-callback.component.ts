import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { CommonModule } from '@angular/common';
import { interval, Subscription } from 'rxjs';

@Component({
  selector: 'strava-callback',
  imports: [],
  templateUrl: './strava-callback.component.html',
  styleUrl: './strava-callback.component.css',
})
export class StravaCallbackComponent implements OnInit, OnDestroy {
  currentStatus = 'Initializing connection...';
  progressPercentage = 0;
  currentStep = 0;
  currentFact = '';

  private progressSubscription?: Subscription;
  private factSubscription?: Subscription;

  private facts = [
    'Mount Everest is 8,848.86 meters tall!',
    'The average runner burns 100 calories per mile',
    'Trail running can burn 50% more calories than road running',
    'Your heart rate can tell you a lot about your fitness level',
    'Running uphill engages 9% more muscles than flat running',
    'The longest recorded run was 350 miles in 80 hours!',
    'Mountain air has less oxygen, making runs more challenging',
    'Running releases endorphins - natural mood boosters!',
  ];

  private statusMessages = [
    'Connecting to Strava...',
    'Verifying your credentials...',
    'Setting up your profile...',
    'Downloading your activities...',
    'Analyzing your routes...',
    'Detecting mountain peaks...',
    'Calculating your stats...',
    'Almost ready...',
  ];

  constructor(
    private route: ActivatedRoute,
    private authService: AuthService,
    private router: Router
  ) { }

  // userID = this.authService.userID;

  ngOnInit(): void {
    this.startLoadingAnimation();
    this.startFactRotation();

    this.route.queryParams.subscribe((params) => {
      const code = params['code'];
      if (code) {
        this.currentStep = 1;
        this.currentStatus = 'Authenticating with Strava...';

        // Exchange code for JWT by calling the backend
        this.authService.loginWithStravaAuth(code).subscribe({
          next: (res) => {
            this.currentStep = 4;
            this.progressPercentage = 100;
            this.currentStatus = 'Success! Redirecting to dashboard...';

            this.authService.storeAccessToken(res.accessToken);
            this.authService.storeRefreshToken(res.refreshToken);
            this.authService.userID.set(res.userID);

            // Wait a moment to show completion, then navigate
            setTimeout(() => {
              this.router.navigate(['/']);
            }, 1500);
          },
          error: (err) => {
            this.currentStatus = 'Authentication failed. Redirecting...';
            console.error('Error logging in', err);
            setTimeout(() => {
              this.router.navigate(['/login']);
            }, 2000);
          },
        });
      } else {
        // No code present? Go back to login
        this.currentStatus = 'No authorization code found. Redirecting...';
        setTimeout(() => {
          this.router.navigate(['/login']);
        }, 2000);
      }
    });
  }

  ngOnDestroy(): void {
    this.progressSubscription?.unsubscribe();
    this.factSubscription?.unsubscribe();
  }

  private startLoadingAnimation(): void {
    let messageIndex = 0;

    this.progressSubscription = interval(800).subscribe(() => {
      if (this.progressPercentage < 90) {
        // Simulate realistic loading progress
        const increment = Math.random() * 15 + 5; // 5-20% increments
        this.progressPercentage = Math.min(
          90,
          (this.progressPercentage + increment).toPrecision(
            0
          ) as unknown as number
        );

        // Update status message
        if (messageIndex < this.statusMessages.length - 1) {
          this.currentStatus = this.statusMessages[messageIndex];
          messageIndex++;
        }

        // Update step based on progress
        if (this.progressPercentage > 20)
          this.currentStep = Math.max(this.currentStep, 2);
        if (this.progressPercentage > 50)
          this.currentStep = Math.max(this.currentStep, 3);
        if (this.progressPercentage > 80)
          this.currentStep = Math.max(this.currentStep, 4);
      }
    });
  }

  private startFactRotation(): void {
    this.currentFact = this.facts[0];
    let factIndex = 0;

    this.factSubscription = interval(6000).subscribe(() => {
      factIndex = (factIndex + 1) % this.facts.length;
      this.currentFact = this.facts[factIndex];
    });
  }
}
