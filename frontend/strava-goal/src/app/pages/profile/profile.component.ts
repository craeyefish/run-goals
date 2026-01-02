import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ProfileService, UserProfile } from '../../services/profile.service';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-profile',
  imports: [CommonModule, FormsModule],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss'],
})
export class ProfileComponent implements OnInit, OnDestroy {
  profile: UserProfile | null = null;
  loading = true;
  error: string | null = null;
  editingUsername = false;
  newUsername = '';
  updateError: string | null = null;
  updateSuccess = false;
  updating = false;

  private destroy$ = new Subject<void>();

  constructor(private profileService: ProfileService) {}

  ngOnInit(): void {
    this.loadProfile();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  loadProfile(): void {
    this.profileService
      .getUserProfile()
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (profile) => {
          this.profile = profile;
          this.newUsername = profile.username || '';
          this.loading = false;
        },
        error: (err) => {
          this.error = 'Failed to load profile';
          this.loading = false;
          console.error('Error loading profile:', err);
        },
      });
  }

  startEditingUsername(): void {
    this.editingUsername = true;
    this.newUsername = this.profile?.username || '';
    this.updateError = null;
    this.updateSuccess = false;
  }

  cancelEditingUsername(): void {
    this.editingUsername = false;
    this.newUsername = this.profile?.username || '';
    this.updateError = null;
  }

  saveUsername(): void {
    if (!this.newUsername || this.newUsername.trim().length < 3) {
      this.updateError = 'Username must be at least 3 characters';
      return;
    }

    if (this.newUsername.trim().length > 50) {
      this.updateError = 'Username must be 50 characters or less';
      return;
    }

    this.updating = true;
    this.updateError = null;

    this.profileService
      .updateUsername(this.newUsername.trim())
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (profile) => {
          this.profile = profile;
          this.newUsername = profile.username || '';
          this.editingUsername = false;
          this.updating = false;
          this.updateSuccess = true;
          setTimeout(() => {
            this.updateSuccess = false;
          }, 3000);
        },
        error: (err) => {
          this.updateError = err.error?.message || 'Failed to update username';
          this.updating = false;
          console.error('Error updating username:', err);
        },
      });
  }
}
