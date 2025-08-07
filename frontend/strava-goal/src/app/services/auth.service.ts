import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, Observable, of, tap } from 'rxjs';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  // selectedUserID signal?

  private stateKey = 'strava_oath_state';
  private accessTokenKey = 'jwt_access_token';
  private refreshTokenKey = 'jwt_refresh_token';
  private accessToken: string | null = null;
  private refreshToken: string | null = null;
  userID = signal<number | null>(null);

  constructor(private http: HttpClient, private router: Router) {}

  storeAccessToken(accessToken: string): void {
    this.accessToken = accessToken;
    localStorage.setItem(this.accessTokenKey, accessToken);
  }

  storeRefreshToken(refreshToken: string): void {
    this.refreshToken = refreshToken;
    localStorage.setItem(this.refreshTokenKey, refreshToken);
  }

  loadAccessTokenFromStorage(): void {
    const saved = localStorage.getItem(this.accessTokenKey);
    if (saved) {
      this.accessToken = saved;
    }
  }

  loadRefreshTokenFromStorage(): void {
    const saved = localStorage.getItem(this.refreshTokenKey);
    if (saved) {
      this.refreshToken = saved;
    }
  }

  getAccessToken(): string | null {
    return this.accessToken;
  }

  doRefresh(): Observable<string> {
    return this.http
      .post<{ accessToken: string }>(
        '/auth/refresh',
        {},
        {
          headers: { Authorization: `Bearer ${this.refreshToken}` },
        }
      )
      .pipe(
        tap((res) => {
          this.storeAccessToken(res.accessToken);
        }),
        map((res) => res.accessToken)
      );
  }

  isLoggedIn(): boolean {
    return !!this.accessToken;
  }

  logout(): void {
    this.accessToken = null;
    localStorage.removeItem(this.accessTokenKey);
    localStorage.removeItem(this.refreshTokenKey);
  }

  loginWithStravaAuth(
    code: string
  ): Observable<{ accessToken: string; refreshToken: string; userID: number }> {
    return this.http.post<{
      accessToken: string;
      refreshToken: string;
      userID: number;
    }>('/auth/strava/callback', {
      code,
    });
  }

  login(): void {
    // Redirect to Strava’s OAuth page

    const clientId = '49851';
    const redirectUri = encodeURIComponent(
      // 'https://summitseekers.co.za/login/strava/callback'
      'http://localhost:4200/login/strava/callback'
    );
    const scope = 'read,activity:read_all';
    const state = this.generateState();

    const stravaAuthUrl =
      `https://www.strava.com/oauth/authorize` +
      `?client_id=${clientId}` +
      `&redirect_uri=${redirectUri}` +
      `&response_type=code` +
      `&scope=${scope}` +
      `&state=${state}`;

    window.location.href = stravaAuthUrl;
  }

  generateState(): string {
    const state = Math.random().toString(36).substring(2, 15);
    localStorage.setItem(this.stateKey, state);
    return state;
  }

  validateState(returnedState: string): boolean {
    const storedState = localStorage.getItem(this.stateKey);
    localStorage.removeItem(this.stateKey);
    return storedState === returnedState;
  }

  getUserID(): number | null {
    return this.userID();
  }
}
