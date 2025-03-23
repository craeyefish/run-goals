import { Component } from "@angular/core";

@Component({
  selector: "app-user-stats",
  standalone: true,
  templateUrl: "./user-stats.component.html",
  styleUrls: ["./user-stats.component.scss"],
})
export class UserStatsComponent {
  username = "John Doe";
  totalKm = 999;
  summits = 99;
}
