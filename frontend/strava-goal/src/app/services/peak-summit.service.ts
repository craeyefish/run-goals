import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { PeakSummaries } from '../models/peak-summit.model';

@Injectable({
  providedIn: 'root',
})
export class PeakSummitService {
  constructor(private http: HttpClient) {}

  getPeakSummaries(): Observable<PeakSummaries[]> {
    // Adjust the endpoint to match your backend
    return this.http.get<PeakSummaries[]>('/api/peak-summaries');
  }
}
