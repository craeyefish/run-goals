import { Component } from "@angular/core";
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: "app-header",
  standalone: true,
  templateUrl: "./header.component.html",
  styleUrls: ["./header.component.scss"],
})
export class HeaderComponent {
  title = "Summit Seekers";
  constructor(private authService: AuthService) { }

  onLogoutClick() {
    this.authService.logout();
  }

  reloadPage(): void {
    window.location.reload();
  }
}
