import { Component } from '@angular/core';

@Component({
  selector: 'app-join-challenge',
  templateUrl: './join-challenge.component.html',
})
export class JoinChallengeComponent {
  password = '';
  correctPassword = 'secret'; // Example only! Use a real approach to hide this.
  mydomain = 'https://craeyebytes.com';

  joinIfPasswordCorrect() {
    if (this.password === this.correctPassword) {
      // Redirect to Strava’s OAuth page
      const clientId = '49851';
      const redirectUri = encodeURIComponent(
        this.mydomain + '/auth/strava/callback'
      );
      const scope = 'read,activity:read_all';

      const stravaAuthUrl =
        `https://www.strava.com/oauth/authorize` +
        `?client_id=${clientId}` +
        `&redirect_uri=${redirectUri}` +
        `&response_type=code` +
        `&scope=${scope}`;

      window.location.href = stravaAuthUrl;
    } else {
      alert('Invalid password!');
    }
  }
}
