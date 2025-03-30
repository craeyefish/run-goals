import { Route } from "@angular/router";
import { AuthComponent } from "./auth.component";

export const AUTH_ROUTES: Route[] = [
  {
    path: "login",
    component: AuthComponent,
  },
  {
    path: "",
    redirectTo: "login",
    pathMatch: "full",
  },
];
