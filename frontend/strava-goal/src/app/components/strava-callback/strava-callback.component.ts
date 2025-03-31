import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'strava-callback',
  imports: [],
  templateUrl: './strava-callback.component.html',
  styleUrl: './strava-callback.component.css',
})
export class StravaCallbackComponent implements OnInit {
  constructor(
    private route: ActivatedRoute,
    private authService: AuthService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.route.queryParams.subscribe((params) => {
      const code = params['code'];
      if (code) {
        // Exchange code for JWT by calling the backend
        this.authService.loginWithStravaAuth(code).subscribe({
          next: (res) => {
            this.authService.storeToken(res.token);
            // Now the user is logged in, navigate to the main page
            this.router.navigate(['/']);
          },
          error: (err) => {
            console.error('Error logging in', err);
            this.router.navigate(['/login']);
          },
        });
      } else {
        // No code present? Go back to login
        this.router.navigate(['/login']);
      }
    });
  }
}
