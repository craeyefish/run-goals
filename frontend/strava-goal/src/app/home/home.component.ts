import { Component } from "@angular/core";
import { ActivitiesListComponent } from "../components/activities-list/activities-list.component";
import { UserStatsComponent } from "../components/user-stats/user-stats.component";
import { ActivityMapComponent } from "../components/activity-map/activity-map.component";

@Component({
    selector: "app-home",
    imports: [UserStatsComponent, ActivitiesListComponent, ActivityMapComponent],
    templateUrl: "./home.component.html",
    styleUrls: ["./home.component.scss"]
})
export class HomeComponent {}
