import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private tokenKey = 'jwt_token';
  private token: string | null = null;

  constructor(private http: HttpClient) {}

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
}
