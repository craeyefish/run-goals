import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { JoinChallengeComponent } from './join-challenge/join-challenge.component';
import { GoalProgressComponent } from './goal-progress/goal-progress.component';
import { PeakSummitTableComponent } from './peak-summit-table/peak-summit-table.component';

const routes: Routes = [
  { path: '', component: GoalProgressComponent },
  { path: 'join-challenge', component: JoinChallengeComponent },
  { path: 'peak-summits', component: PeakSummitTableComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
