import { CommonModule, NgFor } from '@angular/common';
import { Component } from '@angular/core';

interface Activity {
  id: number;
  date: string;
  km: number;
  summits: number;
  location: string;
}

@Component({
  selector: 'app-activities-list',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './activities-list.component.html',
  styleUrls: ['./activities-list.component.scss'],
})
export class ActivitiesListComponent {
  activities: Activity[] = [
    { id: 3, date: '2025-03-01', km: 20, summits: 2, location: 'stellenbosch' },
    { id: 2, date: '2025-02-01', km: 1, summits: 0, location: 'durbanville' },
    { id: 1, date: '2025-01-01', km: 2, summits: 0, location: 'melkbos' },
  ];
}
