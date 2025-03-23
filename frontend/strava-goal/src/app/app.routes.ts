import { Routes } from "@angular/router";
import { HomeComponent } from "./home/home.component";
import { GroupsComponent } from "./groups/groups.component";
import { SummitsComponent } from "./summits/summits.component";
import { ProfileComponent } from "./profile/profile.component";

export const routes: Routes = [
  // todo: auth :(
  { path: "", component: HomeComponent },
  { path: "groups", component: GroupsComponent },
  { path: "summits", component: SummitsComponent },
  { path: "profile", component: ProfileComponent },
];
