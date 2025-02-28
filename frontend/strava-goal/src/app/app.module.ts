import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { GoalProgressComponent } from './goal-progress/goal-progress.component';
import { HttpClientModule } from '@angular/common/http';
import { JoinChallengeComponent } from './join-challenge/join-challenge.component';
import { ActivityMapComponent } from './activity-map/activity-map.component';
import { PeakSummitTableComponent } from './peak-summit-table/peak-summit-table.component';

@NgModule({
  declarations: [
    AppComponent,
    GoalProgressComponent,
    JoinChallengeComponent,
    ActivityMapComponent,
    PeakSummitTableComponent,
  ],
  imports: [BrowserModule, AppRoutingModule, HttpClientModule, FormsModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
