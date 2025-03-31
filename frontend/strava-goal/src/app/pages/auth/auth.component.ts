import { Component } from "@angular/core";
import { LoginComponent } from "src/app/components/login/login.component";

@Component({
  selector: "auth-page",
  standalone: true,
  imports: [LoginComponent],
  templateUrl: "./auth.component.html",
  styleUrls: ["./auth.component.scss"],
})
export class AuthPageComponent { }
