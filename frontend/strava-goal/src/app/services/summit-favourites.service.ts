import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { tap } from 'rxjs/operators';
import { environment } from '../../environments/environment';

@Injectable({
    providedIn: 'root'
})
export class SummitFavouritesService {
    private apiUrl = `${environment.baseUrl}/api`;

    private favouritesSubject = new BehaviorSubject<number[]>([]);
    favourites$ = this.favouritesSubject.asObservable();

    constructor(private http: HttpClient) { }

    /**
     * Get all favourite peak IDs for the current user
     */
    getFavourites(): Observable<number[]> {
        return this.http.get<number[]>(`${this.apiUrl}/summit-favourites`).pipe(
            tap(favourites => this.favouritesSubject.next(favourites))
        );
    }

    /**
     * Add a peak to favourites
     */
    addFavourite(peakId: number): Observable<number[]> {
        return this.http.post<number[]>(`${this.apiUrl}/summit-favourites`, { peak_id: peakId }).pipe(
            tap(favourites => this.favouritesSubject.next(favourites))
        );
    }

    /**
     * Remove a peak from favourites
     */
    removeFavourite(peakId: number): Observable<number[]> {
        return this.http.delete<number[]>(`${this.apiUrl}/summit-favourites?peak_id=${peakId}`).pipe(
            tap(favourites => this.favouritesSubject.next(favourites))
        );
    }

    /**
     * Check if a peak is in favourites
     */
    isFavourite(peakId: number): boolean {
        return this.favouritesSubject.value.includes(peakId);
    }

    /**
     * Get the current favourites value
     */
    getCurrentFavourites(): number[] {
        return this.favouritesSubject.value;
    }
}
