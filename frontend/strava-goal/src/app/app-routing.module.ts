import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { JoinChallengeComponent } from './join-challenge/join-challenge.component';
import { PeakSummitTableComponent } from './peak-summit-table/peak-summit-table.component';
import { HomeComponent } from './home/home.component';

const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'join-challenge', component: JoinChallengeComponent },
  { path: 'peak-summits', component: PeakSummitTableComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
