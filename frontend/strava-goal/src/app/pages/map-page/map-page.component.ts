import { Component } from '@angular/core';
import { ActivityMapComponent } from 'src/app/components/activity-map/activity-map.component';

@Component({
  selector: 'map-page',
  imports: [ActivityMapComponent],
  templateUrl: './map-page.component.html',
  styleUrl: './map-page.component.css',
})
export class MapPageComponent {}
