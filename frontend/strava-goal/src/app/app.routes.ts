import { Routes } from '@angular/router';
import { HomeComponent } from 'src/app/pages/home/home.component';
import { GroupsComponent } from 'src/app/pages/groups/groups.component';
import { SummitsComponent } from 'src/app/pages/summits/summits.component';
import { ProfileComponent } from 'src/app/pages/profile/profile.component';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { AuthLayoutComponent } from './layout/auth-layout/auth-layout.component';
import { LoginComponent } from 'src/app/components/login/login.component';
import { MainLayoutComponent } from './layout/main-layout/main-layout.component';
import { StravaCallbackComponent } from './components/strava-callback/strava-callback.component';
import { StravaCallbackGuard } from './guards/strava-callback.guard';
import { NoAuthGuard } from './guards/no-auth.guard';
import { AuthPageComponent } from './pages/auth/auth.component';

export const routes: Routes = [
  {
    path: 'auth',
    component: AuthLayoutComponent,
    canActivate: [NoAuthGuard],
    children: [
      { path: '', redirectTo: 'login', pathMatch: 'full' },
      { path: 'login', component: AuthPageComponent },
      { path: 'strava/callback', component: StravaCallbackComponent, canActivate: [StravaCallbackGuard] },
    ]
  },
  {
    path: '',
    component: MainLayoutComponent,
    canActivate: [AuthGuard],
    children: [
      { path: '', redirectTo: 'home', pathMatch: 'full' },
      { path: 'home', component: HomeComponent },
      { path: 'groups', component: GroupsComponent },
      { path: 'summits', component: SummitsComponent },
      { path: 'profile', component: ProfileComponent },
    ],
  },
  { path: '**', redirectTo: '' }
];
