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
import { HeaderComponent } from './header/header.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { HomeComponent } from './home/home.component';
import { UserStatsComponent } from './user-stats/user-stats.component';
import { ActivitiesListComponent } from './activities-list/activities-list.component';

@NgModule({
  declarations: [
    AppComponent,
    GoalProgressComponent,
    JoinChallengeComponent,
    ActivityMapComponent,
    PeakSummitTableComponent,
    HeaderComponent,
    SidebarComponent,
    HomeComponent,
    UserStatsComponent,
    ActivitiesListComponent,
  ],
  imports: [BrowserModule, AppRoutingModule, HttpClientModule, FormsModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
