import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { tap } from 'rxjs/operators';
import { environment } from '../../environments/environment';

export interface PersonalYearlyGoal {
    id?: number;
    user_id?: number;
    year: number;
    distance_goal: number;  // km
    elevation_goal: number; // meters
    summit_goal: number;    // count
    target_summits: number[]; // peak IDs
    created_at?: string;
    updated_at?: string;
}

@Injectable({
    providedIn: 'root'
})
export class PersonalGoalsService {
    private apiUrl = `${environment.baseUrl}/api`;

    private currentGoalSubject = new BehaviorSubject<PersonalYearlyGoal | null>(null);
    currentGoal$ = this.currentGoalSubject.asObservable();

    constructor(private http: HttpClient) { }

    /**
     * Get personal goals for a specific year (defaults to current year)
     */
    getGoals(year?: number): Observable<PersonalYearlyGoal> {
        const params = year ? `?year=${year}` : '';
        return this.http.get<PersonalYearlyGoal>(`${this.apiUrl}/personal-goals${params}`).pipe(
            tap(goal => this.currentGoalSubject.next(goal))
        );
    }

    /**
     * Get the current year's goals
     */
    getCurrentYearGoals(): Observable<PersonalYearlyGoal> {
        return this.getGoals(new Date().getFullYear());
    }

    /**
     * Save personal goals (creates or updates)
     */
    saveGoals(goal: PersonalYearlyGoal): Observable<PersonalYearlyGoal> {
        return this.http.post<PersonalYearlyGoal>(`${this.apiUrl}/personal-goals`, goal).pipe(
            tap(savedGoal => this.currentGoalSubject.next(savedGoal))
        );
    }

    /**
     * Quick update for individual goal values
     */
    updateDistanceGoal(distance: number): Observable<PersonalYearlyGoal> {
        const currentGoal = this.currentGoalSubject.value;
        if (!currentGoal) {
            return this.saveGoals({
                year: new Date().getFullYear(),
                distance_goal: distance,
                elevation_goal: 50000,
                summit_goal: 20,
                target_summits: []
            });
        }
        return this.saveGoals({ ...currentGoal, distance_goal: distance });
    }

    updateElevationGoal(elevation: number): Observable<PersonalYearlyGoal> {
        const currentGoal = this.currentGoalSubject.value;
        if (!currentGoal) {
            return this.saveGoals({
                year: new Date().getFullYear(),
                distance_goal: 1000,
                elevation_goal: elevation,
                summit_goal: 20,
                target_summits: []
            });
        }
        return this.saveGoals({ ...currentGoal, elevation_goal: elevation });
    }

    updateSummitGoal(count: number): Observable<PersonalYearlyGoal> {
        const currentGoal = this.currentGoalSubject.value;
        if (!currentGoal) {
            return this.saveGoals({
                year: new Date().getFullYear(),
                distance_goal: 1000,
                elevation_goal: 50000,
                summit_goal: count,
                target_summits: []
            });
        }
        return this.saveGoals({ ...currentGoal, summit_goal: count });
    }

    /**
     * Add a peak to the target summits list
     */
    addTargetSummit(peakId: number): Observable<PersonalYearlyGoal> {
        const currentGoal = this.currentGoalSubject.value;
        if (!currentGoal) {
            return this.saveGoals({
                year: new Date().getFullYear(),
                distance_goal: 1000,
                elevation_goal: 50000,
                summit_goal: 20,
                target_summits: [peakId]
            });
        }

        // Don't add duplicates
        if (currentGoal.target_summits.includes(peakId)) {
            return new Observable(observer => {
                observer.next(currentGoal);
                observer.complete();
            });
        }

        return this.saveGoals({
            ...currentGoal,
            target_summits: [...currentGoal.target_summits, peakId]
        });
    }

    /**
     * Remove a peak from the target summits list
     */
    removeTargetSummit(peakId: number): Observable<PersonalYearlyGoal> {
        const currentGoal = this.currentGoalSubject.value;
        if (!currentGoal) {
            return new Observable(observer => {
                observer.next(null as any);
                observer.complete();
            });
        }

        return this.saveGoals({
            ...currentGoal,
            target_summits: currentGoal.target_summits.filter(id => id !== peakId)
        });
    }

    /**
     * Get the cached current goal value
     */
    getCurrentGoalValue(): PersonalYearlyGoal | null {
        return this.currentGoalSubject.value;
    }
}
