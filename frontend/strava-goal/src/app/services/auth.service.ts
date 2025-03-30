import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private stateKey = "strava_oath_state"
  private tokenKey = 'jwt_token';
  private token: string | null = null;

  constructor(private http: HttpClient, private router: Router) { }

  storeToken(token: string): void {
    this.token = token;
    localStorage.setItem(this.tokenKey, token);
  }

  loadTokenFromStorage(): void {
    const saved = localStorage.getItem(this.tokenKey);
    if (saved) {
      this.token = saved;
    }
  }

  getToken(): string | null {
    return this.token;
  }

  isLoggedIn(): boolean {
    return !!this.token;
  }

  logout(): void {
    this.token = null;
    localStorage.removeItem(this.tokenKey);
  }

  loginWithStravaAuth(code: string): Observable<{ token: string }> {
    return this.http.post<{ token: string }>('/auth/strava/callback', {
      code,
    });
  }

  login(): void {
    // Redirect to Strava’s OAuth page

    // Temporary for skipping login while testing
    // if (1 == 1) {
    //   this.token = 'testToken';
    //   this.router.navigate(['']);
    //   return
    // }

    const clientId = '49851';
    const redirectUri = encodeURIComponent(
      'https://craeyebytes.com/auth/strava/callback'
    );
    const scope = 'read,activity:read_all';
    const state = this.generateState()

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
    localStorage.setItem(this.stateKey, state)
    return state
  }

  validateState(returnedState: string): boolean {
    const storedState = localStorage.getItem(this.stateKey);
    localStorage.removeItem(this.stateKey);
    return storedState === returnedState;
  }
}
