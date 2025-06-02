import { Component } from '@angular/core';
import { ActivitiesListComponent } from '../../components/activities-list/activities-list.component';

@Component({
  selector: 'activities-page',
  imports: [ActivitiesListComponent],
  templateUrl: './activities-page.component.html',
  styleUrl: './activities-page.component.css',
})
export class ActivitiesPageComponent {}
