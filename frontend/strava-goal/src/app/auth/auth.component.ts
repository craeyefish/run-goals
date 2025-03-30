import { Component } from "@angular/core";
import { AuthService } from "src/app/services/auth.service";
import { LoginComponent } from "../components/login/login.component";
import { StravaCallbackComponent } from "../components/strava-callback/strava-callback.component";

@Component({
  selector: "app-auth",
  standalone: true,
  imports: [LoginComponent, StravaCallbackComponent],
  templateUrl: "auth.component.html",
  styleUrl: "auth.component.scss",
})
export class AuthComponent {
  username: string = "";
  password: string = "";

  constructor(private authService: AuthService) {}

  login() {}
}
