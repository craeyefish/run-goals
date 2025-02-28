import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Peak {
  id: number;
  osm_id: number;
  lat: number;
  lon: number;
  name: string;
  elev_m: number;
  is_summited: boolean; // new boolean to indicate if this peak was visited by the user/group
}

@Injectable({
  providedIn: 'root',
})
export class PeakService {
  constructor(private http: HttpClient) {}

  // If you have a bounding-box approach, you'd pass minLat etc. as query params.
  // For now, assume we want all peaks:
  getPeaks(): Observable<Peak[]> {
    return this.http.get<Peak[]>('/api/peaks'); // Adjust path if needed
  }
}
