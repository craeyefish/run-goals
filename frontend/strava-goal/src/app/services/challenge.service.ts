import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, signal, computed } from '@angular/core';
import { Observable } from 'rxjs';
import {
    Challenge,
    ChallengeWithProgress,
    ChallengePeakWithDetails,
    ChallengeParticipantWithUser,
    ChallengeSummitLogWithDetails,
    LeaderboardEntry,
    CreateChallengeRequest,
    CreateChallengeResponse,
    UpdateChallengeRequest,
    SetChallengePeaksRequest,
    AddGroupToChallengeRequest,
    RecordSummitRequest,
    JoinChallengeByCodeRequest,
    ChallengeDetailResponse,
    ChallengeListResponse,
    PublicChallengeListResponse,
} from '../models/challenge.model';
import { ActivityWithUser } from './activity.service';

@Injectable({
    providedIn: 'root',
})
export class ChallengeService {
    // State signals
    challenges = signal<ChallengeWithProgress[]>([]);
    selectedChallenge = signal<ChallengeWithProgress | null>(null);
    challengePeaks = signal<ChallengePeakWithDetails[]>([]);
    participants = signal<ChallengeParticipantWithUser[]>([]);
    summitLog = signal<ChallengeSummitLogWithDetails[]>([]);
    challengeActivities = signal<ActivityWithUser[]>([]);
    featuredChallenges = signal<Challenge[]>([]);
    publicChallenges = signal<Challenge[]>([]);
    isLoading = signal<boolean>(false);

    // Computed values
    completedPeaks = computed(() =>
        this.challengePeaks().filter(p => p.isSummited)
    );

    progressPercentage = computed(() => {
        const challenge = this.selectedChallenge();
        if (!challenge) return 0;

        switch (challenge.goalType) {
            case 'specific_summits':
                if (challenge.totalPeaks === 0) return 0;
                return Math.round((challenge.completedPeaks / challenge.totalPeaks) * 100);

            case 'distance':
                if (!challenge.targetValue || challenge.targetValue === 0) return 0;
                return Math.min(100, Math.round((challenge.currentDistance / challenge.targetValue) * 100));

            case 'elevation':
                if (!challenge.targetValue || challenge.targetValue === 0) return 0;
                return Math.min(100, Math.round((challenge.currentElevation / challenge.targetValue) * 100));

            case 'summit_count':
                if (!challenge.targetSummitCount || challenge.targetSummitCount === 0) return 0;
                return Math.min(100, Math.round((challenge.currentSummitCount / challenge.targetSummitCount) * 100));

            default:
                return 0;
        }
    });

    constructor(private http: HttpClient) { }

    // ==================== Challenge CRUD ====================

    createChallenge(request: CreateChallengeRequest): Observable<CreateChallengeResponse> {
        return this.http.post<CreateChallengeResponse>('/api/challenges', request);
    }

    getChallenge(id: number): Observable<ChallengeDetailResponse> {
        const params = new HttpParams().set('id', id);
        return this.http.get<ChallengeDetailResponse>('/api/challenge', { params });
    }

    updateChallenge(request: UpdateChallengeRequest): Observable<void> {
        return this.http.put<void>('/api/challenge', request);
    }

    deleteChallenge(id: number): Observable<void> {
        const params = new HttpParams().set('id', id);
        return this.http.delete<void>('/api/challenge', { params });
    }

    // ==================== Discovery ====================

    getUserChallenges(): Observable<ChallengeListResponse> {
        return this.http.get<ChallengeListResponse>('/api/challenges');
    }

    getFeaturedChallenges(): Observable<PublicChallengeListResponse> {
        return this.http.get<PublicChallengeListResponse>('/api/challenges/featured');
    }

    getPublicChallenges(region?: string, limit?: number, offset?: number): Observable<PublicChallengeListResponse> {
        let params = new HttpParams();
        if (region) params = params.set('region', region);
        if (limit) params = params.set('limit', limit);
        if (offset) params = params.set('offset', offset);
        return this.http.get<PublicChallengeListResponse>('/api/challenges/public', { params });
    }

    searchChallenges(query: string, limit?: number): Observable<PublicChallengeListResponse> {
        let params = new HttpParams().set('q', query);
        if (limit) params = params.set('limit', limit);
        return this.http.get<PublicChallengeListResponse>('/api/challenges/search', { params });
    }

    // ==================== Peaks ====================

    getChallengePeaks(challengeId: number): Observable<ChallengePeakWithDetails[]> {
        const params = new HttpParams().set('challengeId', challengeId);
        return this.http.get<ChallengePeakWithDetails[]>('/api/challenge-peaks', { params });
    }

    setChallengePeaks(challengeId: number, request: SetChallengePeaksRequest): Observable<void> {
        const params = new HttpParams().set('challengeId', challengeId);
        return this.http.put<void>('/api/challenge-peaks', request, { params });
    }

    // ==================== Participation ====================

    joinChallenge(challengeId: number): Observable<void> {
        const params = new HttpParams().set('challengeId', challengeId);
        return this.http.post<void>('/api/challenge-join', {}, { params });
    }

    joinChallengeByCode(joinCode: string): Observable<Challenge> {
        const request: JoinChallengeByCodeRequest = { joinCode };
        return this.http.post<Challenge>('/api/challenge-join-by-code', request);
    }

    leaveChallenge(challengeId: number): Observable<void> {
        const params = new HttpParams().set('challengeId', challengeId);
        return this.http.delete<void>('/api/challenge-leave', { params });
    }

    lockChallenge(challengeId: number): Observable<void> {
        const params = new HttpParams().set('challengeId', challengeId);
        return this.http.post<void>('/api/challenge-lock', {}, { params });
    }

    getParticipants(challengeId: number): Observable<ChallengeParticipantWithUser[]> {
        const params = new HttpParams().set('challengeId', challengeId);
        return this.http.get<ChallengeParticipantWithUser[]>('/api/challenge-participants', { params });
    }

    getLeaderboard(challengeId: number): Observable<LeaderboardEntry[]> {
        const params = new HttpParams().set('challengeId', challengeId);
        return this.http.get<LeaderboardEntry[]>('/api/challenge-leaderboard', { params });
    }

    // ==================== Progress ====================

    getSummitLog(challengeId: number, userId?: number): Observable<ChallengeSummitLogWithDetails[]> {
        let params = new HttpParams().set('challengeId', challengeId);
        if (userId) params = params.set('userId', userId);
        return this.http.get<ChallengeSummitLogWithDetails[]>('/api/challenge-summit-log', { params });
    }

    recordSummit(challengeId: number, request: RecordSummitRequest): Observable<void> {
        const params = new HttpParams().set('challengeId', challengeId);
        return this.http.post<void>('/api/challenge-summit', request, { params });
    }

    // ==================== Activities ====================

    getChallengeActivities(challengeId: number): Observable<ActivityWithUser[]> {
        const params = new HttpParams().set('challengeId', challengeId);
        return this.http.get<ActivityWithUser[]>('/api/challenge-activities', { params });
    }

    // ==================== Groups ====================

    addGroupToChallenge(challengeId: number, request: AddGroupToChallengeRequest): Observable<void> {
        const params = new HttpParams().set('challengeId', challengeId);
        return this.http.post<void>('/api/challenge-group', request, { params });
    }

    removeGroupFromChallenge(challengeId: number, groupId: number): Observable<void> {
        const params = new HttpParams()
            .set('challengeId', challengeId)
            .set('groupId', groupId);
        return this.http.delete<void>('/api/challenge-group', { params });
    }

    getGroupChallenges(groupId: number): Observable<PublicChallengeListResponse> {
        const params = new HttpParams().set('groupId', groupId);
        return this.http.get<PublicChallengeListResponse>('/api/group-challenges', { params });
    }

    // ==================== State Management ====================

    loadUserChallenges(): void {
        this.isLoading.set(true);
        this.getUserChallenges().subscribe({
            next: (response) => {
                this.challenges.set(response.challenges || []);
                this.isLoading.set(false);
            },
            error: (err) => {
                console.error('Failed to load challenges', err);
                this.isLoading.set(false);
            },
        });
    }

    loadChallenge(id: number): void {
        this.isLoading.set(true);
        this.getChallenge(id).subscribe({
            next: (response) => {
                this.selectedChallenge.set(response.challenge);
                this.challengePeaks.set(response.peaks || []);
                this.participants.set(response.participants || []);
                this.isLoading.set(false);
            },
            error: (err) => {
                console.error('Failed to load challenge', err);
                this.isLoading.set(false);
            },
        });
    }

    loadFeaturedChallenges(): void {
        this.getFeaturedChallenges().subscribe({
            next: (response) => {
                this.featuredChallenges.set(response.challenges || []);
            },
            error: (err) => {
                console.error('Failed to load featured challenges', err);
            },
        });
    }

    loadPublicChallenges(region?: string): void {
        this.getPublicChallenges(region).subscribe({
            next: (response) => {
                this.publicChallenges.set(response.challenges || []);
            },
            error: (err) => {
                console.error('Failed to load public challenges', err);
            },
        });
    }

    loadSummitLog(challengeId: number): void {
        this.getSummitLog(challengeId).subscribe({
            next: (log) => {
                this.summitLog.set(log || []);
            },
            error: (err) => {
                console.error('Failed to load summit log', err);
            },
        });
    }

    loadChallengeActivities(challengeId: number): void {
        this.getChallengeActivities(challengeId).subscribe({
            next: (activities) => {
                this.challengeActivities.set(activities || []);
            },
            error: (err) => {
                console.error('Failed to load challenge activities', err);
            },
        });
    }

    selectChallenge(challenge: ChallengeWithProgress | null): void {
        this.selectedChallenge.set(challenge);
        if (challenge) {
            this.loadChallenge(challenge.id);
        } else {
            this.challengePeaks.set([]);
            this.participants.set([]);
            this.summitLog.set([]);
        }
    }

    clearSelection(): void {
        this.selectedChallenge.set(null);
        this.challengePeaks.set([]);
        this.participants.set([]);
        this.summitLog.set([]);
    }

    // ==================== Helper Methods ====================

    getChallengeTypeLabel(type: string): string {
        switch (type) {
            case 'predefined': return 'Featured';
            case 'custom': return 'Custom';
            case 'yearly_goal': return 'Yearly Goal';
            default: return type;
        }
    }

    getGoalTypeLabel(goalType: string): string {
        switch (goalType) {
            case 'distance': return 'Distance';
            case 'elevation': return 'Elevation Gain';
            case 'summit_count': return 'Summit Count';
            case 'specific_summits': return 'Specific Peaks';
            default: return goalType;
        }
    }

    getGoalTypeIcon(goalType: string): string {
        switch (goalType) {
            case 'distance': return '📏';
            case 'elevation': return '⛰️';
            case 'summit_count': return '🏔️';
            case 'specific_summits': return '📍';
            default: return '🎯';
        }
    }

    getCompetitionModeLabel(mode: string): string {
        switch (mode) {
            case 'collaborative': return 'Collaborative';
            case 'competitive': return 'Competitive';
            default: return mode;
        }
    }

    getVisibilityLabel(visibility: string): string {
        switch (visibility) {
            case 'private': return 'Private';
            case 'friends': return 'Friends Only';
            case 'public': return 'Public';
            default: return visibility;
        }
    }

    getDifficultyColor(difficulty: string | undefined): string {
        switch (difficulty?.toLowerCase()) {
            case 'easy': return '#4CAF50';
            case 'moderate': return '#FF9800';
            case 'hard': return '#f44336';
            case 'expert': return '#9C27B0';
            default: return '#9E9E9E';
        }
    }

    formatTargetValue(challenge: Challenge): string {
        switch (challenge.goalType) {
            case 'distance':
                return challenge.targetValue
                    ? `${(challenge.targetValue / 1000).toFixed(1)} km`
                    : '';
            case 'elevation':
                return challenge.targetValue
                    ? `${challenge.targetValue.toLocaleString()} m`
                    : '';
            case 'summit_count':
                return challenge.targetSummitCount
                    ? `${challenge.targetSummitCount} summits`
                    : '';
            case 'specific_summits':
                return 'Specific peaks list';
            default:
                return '';
        }
    }
}
