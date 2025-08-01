import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-support',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './support.component.html',
  styleUrls: ['./support.component.scss'],
})
export class SupportComponent {
  stravaAthleteId: string = '';
  isDeleting: boolean = false;
  deleteSuccess: boolean = false;
  deleteError: string = '';

  constructor(
    private http: HttpClient,
    private router: Router,
    private authService: AuthService
  ) {}

  onDeleteAccount() {
    if (!this.stravaAthleteId.trim()) {
      this.deleteError = 'Please enter your Strava Athlete ID';
      return;
    }

    if (
      !confirm(
        '⚠️ Are you absolutely sure you want to delete all your data? This action cannot be undone!'
      )
    ) {
      return;
    }

    this.isDeleting = true;
    this.deleteError = '';

    this.http
      .delete(`/support/delete-account/${this.stravaAthleteId}`)
      .subscribe({
        next: () => {
          this.deleteSuccess = true;
          this.isDeleting = false;
          this.stravaAthleteId = '';

          // Automatically log out the user since their account is deleted
          this.authService.logout();

          // Redirect to login page after a brief delay to show success message
          setTimeout(() => {
            this.router.navigate(['/login']);
          }, 2000);
        },
        error: (error) => {
          this.deleteError =
            error.error?.message ||
            'Failed to delete account. Please try again.';
          this.isDeleting = false;
        },
      });
  }

  navigateToHome() {
    // If account was deleted, user should go to login instead
    if (this.deleteSuccess) {
      this.router.navigate(['/login']);
    } else {
      this.router.navigate(['/home']);
    }
  }
}
