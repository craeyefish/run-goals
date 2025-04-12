import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';

export interface Peak {
  id: number;
  osm_id: number;
  latitude: number;
  longitude: number;
  name: string;
  elevation_meters: number;
  is_summited: boolean; // new boolean to indicate if this peak was visited by the user/group
}

@Injectable({
  providedIn: 'root',
})
export class PeakService {
  private peaksSubject = new BehaviorSubject<Peak[] | null>(null);
  peaks$ = this.peaksSubject.asObservable();
  private loading = false;

  constructor(private http: HttpClient) {}

  // If you have a bounding-box approach, you'd pass minLat etc. as query params.
  // For now, assume we want all peaks:
  loadPeaks(): void {
    if (this.peaksSubject.value) {
      return;
    }

    if (this.loading) {
      return;
    }

    this.loading = true;

    this.http.get<Peak[]>('/api/peaks').subscribe({
      next: (acts) => {
        this.peaksSubject.next(acts);
        this.loading = false;
      },
      error: (err) => {
        console.error('Failed to load peaks', err);
        this.loading = false;
      },
    });
  }
}
