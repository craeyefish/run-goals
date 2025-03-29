import { Component } from '@angular/core';

@Component({
  selector: 'login',
  imports: [],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css',
})
export class LoginComponent {
  constructor() {}

  loginWithStrava() {
    // Redirect to Strava’s OAuth page
    const clientId = '49851';
    const redirectUri = encodeURIComponent(
      'https://craeyebytes.com/strava/callback'
    );
    const scope = 'read,activity:read_all';

    const stravaAuthUrl =
      `https://www.strava.com/oauth/authorize` +
      `?client_id=${clientId}` +
      `&redirect_uri=${redirectUri}` +
      `&response_type=code` +
      `&scope=${scope}`;

    window.location.href = stravaAuthUrl;
  }
}
