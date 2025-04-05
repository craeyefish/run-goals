import { Injectable } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  CanActivate,
  GuardResult,
  MaybeAsync,
  Router,
  RouterStateSnapshot,
} from '@angular/router';
import { AuthService } from '../services/auth.service';

@Injectable({
  providedIn: 'root',
})
export class StravaCallbackGuard implements CanActivate {
  constructor(private router: Router, private authService: AuthService) {}

  canActivate(route: ActivatedRouteSnapshot): boolean {
    const code = route.queryParams['code'];
    const state = route.queryParams['state'];

    // Ensure auth code provided
    if (!code) {
      this.router.navigate(['/login']);
      return false;
    }

    // Ensure strava request by checking state
    if (!state || !this.authService.validateState(state)) {
      this.router.navigate(['/login']);
      return false;
    }

    return true;
  }
}
