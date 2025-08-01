import { Routes } from '@angular/router';
import { HomeComponent } from 'src/app/pages/home/home.component';
import { GroupsComponent } from 'src/app/pages/groups/groups.component';
import { SummitsPageComponent } from 'src/app/pages/summits/summits-page.component';
import { ProfileComponent } from 'src/app/pages/profile/profile.component';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { AuthLayoutComponent } from './layout/auth-layout/auth-layout.component';
import { MainLayoutComponent } from './layout/main-layout/main-layout.component';
import { StravaCallbackComponent } from './components/strava-callback/strava-callback.component';
import { StravaCallbackGuard } from './guards/strava-callback.guard';
import { NoAuthGuard } from './guards/no-auth.guard';
import { AuthPageComponent } from './pages/auth/auth.component';
import { HikeGangActivitiesComponent } from './hg/hike-gang-activities/hike-gang-activities.component';
import { HikeGangCoverComponent } from './hg/hike-gang-cover/hike-gang-cover.component';
import { HikeGangHomeComponent } from './hg/hike-gang-home/hike-gang-home.component';
import { HikeGangBadgesComponent } from './hg/hike-gang-badges/hike-gang-badges.component';
import { GroupsDetailsPageComponent } from './pages/groups/group-details/group-details.component';
import { GroupsListPageComponent } from './pages/groups/group-list/group-list.component';
import { MapPageComponent } from './pages/map-page/map-page.component';
import { ActivitiesPageComponent } from './pages/activities-page/activities-page.component';
import { SupportComponent } from './pages/support/support.component';

export const routes: Routes = [
  {
    path: 'login',
    component: AuthLayoutComponent,
    canActivate: [NoAuthGuard],
    children: [
      { path: '', component: AuthPageComponent },
      {
        path: 'strava/callback',
        component: StravaCallbackComponent,
        canActivate: [StravaCallbackGuard],
      },
    ],
  },
  {
    path: 'hg',
    component: HikeGangHomeComponent,
    children: [
      {
        path: '',
        component: HikeGangCoverComponent,
      },
      {
        path: 'activities',
        component: HikeGangActivitiesComponent,
      },
      {
        path: 'badges',
        component: HikeGangBadgesComponent,
      },
    ],
  },
  {
    path: '',
    component: MainLayoutComponent,
    canActivate: [AuthGuard],
    children: [
      { path: '', redirectTo: 'home', pathMatch: 'full' },
      { path: 'home', component: HomeComponent },
      {
        path: 'groups',
        component: GroupsComponent,
        children: [
          { path: '', component: GroupsListPageComponent },
          { path: ':code', component: GroupsDetailsPageComponent },
        ],
      },
      { path: 'activities', component: ActivitiesPageComponent },
      { path: 'summits', component: SummitsPageComponent },
      { path: 'map', component: MapPageComponent },
      { path: 'profile', component: ProfileComponent },
      { path: 'support', component: SupportComponent },
    ],
  },
  { path: '**', redirectTo: '' },
];
