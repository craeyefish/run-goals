// Challenge types matching backend models

export type ChallengeType = 'predefined' | 'custom' | 'yearly_goal';
export type GoalType = 'distance' | 'elevation' | 'summit_count' | 'specific_summits';
export type CompetitionMode = 'collaborative' | 'competitive';
export type Visibility = 'private' | 'friends' | 'public';

export interface Challenge {
    id: number;
    name: string;
    description?: string;
    challengeType: ChallengeType;
    goalType: GoalType;
    competitionMode: CompetitionMode;
    visibility: Visibility;
    startDate?: string;
    deadline?: string;
    createdByUserId?: number;
    createdByGroupId?: number;
    targetValue?: number; // For distance/elevation (in meters)
    targetSummitCount?: number; // For summit_count
    region?: string;
    difficulty?: string;
    isFeatured: boolean;
    joinCode: string;
    isLocked: boolean;
    createdAt: string;
    updatedAt: string;
}

export interface ChallengeWithProgress extends Challenge {
    // For specific_summits and summit_count goal types
    totalPeaks: number;
    completedPeaks: number;
    // For distance goal type (meters, but frontend converts to km)
    currentDistance: number;
    // For elevation goal type (meters)
    currentElevation: number;
    // For summit_count goal type
    currentSummitCount: number;
    // Metadata
    isJoined: boolean;
    isCompleted: boolean;
}

export interface ChallengePeak {
    id: number;
    challengeId: number;
    peakId: number;
    sortOrder: number;
    createdAt: string;
}

export interface ChallengePeakWithDetails extends ChallengePeak {
    name: string;
    altName?: string;
    latitude: number;
    longitude: number;
    elevation: number;
    region?: string;
    isSummited: boolean;
}

export interface ChallengeParticipant {
    id: number;
    challengeId: number;
    userId: number;
    joinedAt: string;
    completedAt?: string;
    peaksCompleted: number; // For specific_summits
    totalPeaks: number; // For specific_summits
    totalDistance: number; // For distance (in meters)
    totalElevation: number; // For elevation (in meters)
    totalSummitCount: number; // For summit_count
}

export interface ChallengeParticipantWithUser extends ChallengeParticipant {
    userName: string;
    stravaAthleteId: number;
    profilePicture?: string;
}

export interface LeaderboardEntry {
    rank: number;
    userId: number;
    userName: string;
    stravaAthleteId: number;
    profilePicture?: string;
    peaksCompleted: number;
    totalPeaks: number;
    totalDistance: number;
    totalElevation: number;
    totalSummitCount: number;
    progress: number; // Percentage 0-100
    joinedAt: string;
    completedAt?: string;
}

export interface ChallengeGroup {
    id: number;
    challengeId: number;
    groupId: number;
    startedAt: string;
    completedAt?: string;
    deadlineOverride?: string;
}

export interface ChallengeGroupWithDetails extends ChallengeGroup {
    groupName: string;
    memberCount: number;
}

export interface ChallengeSummitLog {
    id: number;
    challengeId: number;
    userId: number;
    peakId?: number;
    activityId?: number;
    summitedAt: string;
    createdAt: string;
}

export interface ChallengeSummitLogWithDetails extends ChallengeSummitLog {
    peakName?: string;
    peakElevation?: number;
}

// Request/Response types

export interface CreateChallengeRequest {
    name: string;
    description?: string;
    challengeType: ChallengeType;
    goalType: GoalType;
    competitionMode: CompetitionMode;
    visibility: Visibility;
    startDate?: string;
    deadline?: string;
    targetValue?: number; // For distance/elevation (in meters)
    targetSummitCount?: number; // For summit_count
    region?: string;
    difficulty?: string;
    peakIds: number[];
}

export interface UpdateChallengeRequest {
    id: number;
    name: string;
    description?: string;
    challengeType: ChallengeType;
    goalType: GoalType;
    competitionMode: CompetitionMode;
    visibility: Visibility;
    startDate?: string;
    deadline?: string;
    targetValue?: number; // For distance/elevation (in meters)
    targetSummitCount?: number; // For summit_count
    region?: string;
    difficulty?: string;
}

export interface JoinChallengeByCodeRequest {
    joinCode: string;
}

export interface SetChallengePeaksRequest {
    peakIds: number[];
}

export interface AddGroupToChallengeRequest {
    groupId: number;
    deadlineOverride?: string;
}

export interface RecordSummitRequest {
    peakId: number;
    activityId?: number;
    summitedAt: string;
}

export interface CreateChallengeResponse {
    id: number;
}

export interface ChallengeDetailResponse {
    challenge: ChallengeWithProgress;
    peaks: ChallengePeakWithDetails[];
    participants: ChallengeParticipantWithUser[];
}

export interface ChallengeListResponse {
    challenges: ChallengeWithProgress[];
    total: number;
}

export interface PublicChallengeListResponse {
    challenges: Challenge[];
    total: number;
}
