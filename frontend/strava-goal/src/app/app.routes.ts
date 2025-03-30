import { Routes } from '@angular/router';
import { HomeComponent } from 'src/app/pages/home/home.component';
import { GroupsComponent } from 'src/app/pages/groups/groups.component';
import { SummitsComponent } from 'src/app/pages/summits/summits.component';
import { ProfileComponent } from 'src/app/pages/profile/profile.component';
import { AuthGuard } from 'src/app/guards/auth.guard';
import { AuthLayoutComponent } from './layout/auth-layout/auth-layout.component';
import { AuthComponent } from './auth/auth.component';
import { MainLayoutComponent } from './layout/main-layout/main-layout.component';

export const routes: Routes = [
  {
    path: 'auth',
    component: AuthLayoutComponent,
    children: [
      { path: 'login', component: AuthComponent },
      { path: '', redirectTo: 'login', pathMatch: 'full' }
    ]
  },
  {
    path: '',
    component: MainLayoutComponent,
    canActivate: [AuthGuard],
    children: [
      { path: 'home', component: HomeComponent },
      { path: 'groups', component: GroupsComponent },
      { path: 'summits', component: SummitsComponent },
      { path: 'profile', component: ProfileComponent },
    ],
  },
];
