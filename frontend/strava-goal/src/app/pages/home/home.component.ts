import { Component } from "@angular/core";
import { ActivitiesListComponent } from "src/app/components/activities-list/activities-list.component";
import { UserStatsComponent } from "src/app/components/user-stats/user-stats.component";
import { ActivityMapComponent } from "src/app/components/activity-map/activity-map.component";

@Component({
  selector: "app-home",
  standalone: true,
  imports: [UserStatsComponent, ActivitiesListComponent, ActivityMapComponent],
  templateUrl: "./home.component.html",
  styleUrls: ["./home.component.scss"],
})
export class HomeComponent {}
