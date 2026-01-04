package models

import "time"

// ChallengeType determines challenge behavior
type ChallengeType string

const (
	ChallengeTypePredefined ChallengeType = "predefined"  // Admin-created, discoverable
	ChallengeTypeCustom     ChallengeType = "custom"      // User-created
	ChallengeTypeYearlyGoal ChallengeType = "yearly_goal" // Special yearly summit count goal
)

// CompetitionMode determines how participants interact
type CompetitionMode string

const (
	CompetitionModeCollaborative CompetitionMode = "collaborative" // Work together
	CompetitionModeCompetitive   CompetitionMode = "competitive"   // Leaderboard
)

// Visibility determines who can see/join the challenge
type Visibility string

const (
	VisibilityPrivate Visibility = "private" // Only creator and invited
	VisibilityFriends Visibility = "friends" // Friends can see
	VisibilityPublic  Visibility = "public"  // Anyone can see
)

// GoalType determines what participants are trying to achieve
type GoalType string

const (
	GoalTypeDistance        GoalType = "distance"         // Total distance (km)
	GoalTypeElevation       GoalType = "elevation"        // Total elevation gain (m)
	GoalTypeSummitCount     GoalType = "summit_count"     // Number of summits
	GoalTypeSpecificSummits GoalType = "specific_summits" // Specific list of peaks
)

// Challenge represents a summit challenge
type Challenge struct {
	ID                 int64           `json:"id" db:"id"`
	Name               string          `json:"name" db:"name"`
	Description        *string         `json:"description" db:"description"`
	ChallengeType      ChallengeType   `json:"challengeType" db:"challenge_type"`
	GoalType           GoalType        `json:"goalType" db:"goal_type"`
	CompetitionMode    CompetitionMode `json:"competitionMode" db:"competition_mode"`
	Visibility         Visibility      `json:"visibility" db:"visibility"`
	StartDate          *time.Time      `json:"startDate" db:"start_date"`
	Deadline           *time.Time      `json:"deadline" db:"deadline"`
	CreatedByUserID    *int64          `json:"createdByUserId" db:"created_by_user_id"`
	CreatedByGroupID   *int64          `json:"createdByGroupId" db:"created_by_group_id"`
	TargetValue        *float64        `json:"targetValue" db:"target_value"`             // For distance/elevation (in meters)
	TargetSummitCount  *int            `json:"targetSummitCount" db:"target_summit_count"` // For summit_count
	Region             *string         `json:"region" db:"region"`
	Difficulty         *string         `json:"difficulty" db:"difficulty"`
	IsFeatured         bool            `json:"isFeatured" db:"is_featured"`
	JoinCode           string          `json:"joinCode" db:"join_code"`
	IsLocked           bool            `json:"isLocked" db:"is_locked"`
	CreatedAt          time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time       `json:"updatedAt" db:"updated_at"`
}

// ChallengeWithProgress includes computed progress fields
type ChallengeWithProgress struct {
	Challenge
	// For specific_summits and summit_count goal types
	TotalPeaks     int `json:"totalPeaks"`
	CompletedPeaks int `json:"completedPeaks"`
	// For distance goal type (meters, but frontend converts to km)
	CurrentDistance float64 `json:"currentDistance"`
	// For elevation goal type (meters)
	CurrentElevation float64 `json:"currentElevation"`
	// For summit_count goal type
	CurrentSummitCount int `json:"currentSummitCount"`
	// Metadata
	IsJoined    bool `json:"isJoined"`
	IsCompleted bool `json:"isCompleted"`
}

// ChallengePeak links a peak to a challenge
type ChallengePeak struct {
	ID          int64     `json:"id" db:"id"`
	ChallengeID int64     `json:"challengeId" db:"challenge_id"`
	PeakID      int64     `json:"peakId" db:"peak_id"`
	SortOrder   int       `json:"sortOrder" db:"sort_order"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

// ChallengePeakWithDetails includes peak information
type ChallengePeakWithDetails struct {
	ChallengePeak
	Name       string  `json:"name" db:"name"`
	AltName    *string `json:"altName" db:"alt_name"`
	Latitude   float64 `json:"latitude" db:"latitude"`
	Longitude  float64 `json:"longitude" db:"longitude"`
	Elevation  float64 `json:"elevation" db:"elevation"`
	Region     *string `json:"region" db:"region"`
	IsSummited bool    `json:"isSummited"` // Computed per user
}

// ChallengeParticipant tracks user participation
type ChallengeParticipant struct {
	ID               int64      `json:"id" db:"id"`
	ChallengeID      int64      `json:"challengeId" db:"challenge_id"`
	UserID           int64      `json:"userId" db:"user_id"`
	JoinedAt         time.Time  `json:"joinedAt" db:"joined_at"`
	CompletedAt      *time.Time `json:"completedAt" db:"completed_at"`
	PeaksCompleted   int        `json:"peaksCompleted" db:"peaks_completed"`     // For specific_summits
	TotalPeaks       int        `json:"totalPeaks" db:"total_peaks"`             // For specific_summits
	TotalDistance    float64    `json:"totalDistance" db:"total_distance"`       // For distance (in meters)
	TotalElevation   float64    `json:"totalElevation" db:"total_elevation"`     // For elevation (in meters)
	TotalSummitCount int        `json:"totalSummitCount" db:"total_summit_count"` // For summit_count
}

// ChallengeParticipantWithUser includes user information
type ChallengeParticipantWithUser struct {
	ChallengeParticipant
	UserName        string  `json:"userName" db:"user_name"`
	StravaAthleteID int64   `json:"stravaAthleteId" db:"strava_athlete_id"`
	ProfilePicture  *string `json:"profilePicture" db:"profile_picture"`
}

// ChallengeGroup tracks group participation
type ChallengeGroup struct {
	ID               int64      `json:"id" db:"id"`
	ChallengeID      int64      `json:"challengeId" db:"challenge_id"`
	GroupID          int64      `json:"groupId" db:"group_id"`
	StartedAt        time.Time  `json:"startedAt" db:"started_at"`
	CompletedAt      *time.Time `json:"completedAt" db:"completed_at"`
	DeadlineOverride *time.Time `json:"deadlineOverride" db:"deadline_override"`
}

// ChallengeGroupWithDetails includes group information
type ChallengeGroupWithDetails struct {
	ChallengeGroup
	GroupName   string `json:"groupName" db:"group_name"`
	MemberCount int    `json:"memberCount" db:"member_count"`
}

// ChallengeSummitLog tracks summit credits toward challenges
type ChallengeSummitLog struct {
	ID          int64     `json:"id" db:"id"`
	ChallengeID int64     `json:"challengeId" db:"challenge_id"`
	UserID      int64     `json:"userId" db:"user_id"`
	PeakID      *int64    `json:"peakId" db:"peak_id"`
	ActivityID  *int64    `json:"activityId" db:"activity_id"`
	SummitedAt  time.Time `json:"summitedAt" db:"summited_at"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

// ChallengeSummitLogWithDetails includes peak and activity info
type ChallengeSummitLogWithDetails struct {
	ChallengeSummitLog
	PeakName      *string  `json:"peakName" db:"peak_name"`
	PeakElevation *float64 `json:"peakElevation" db:"peak_elevation"`
}

// LeaderboardEntry represents a ranked participant in a competitive challenge
type LeaderboardEntry struct {
	Rank            int        `json:"rank"`
	UserID          int64      `json:"userId"`
	UserName        string     `json:"userName"`
	StravaAthleteID int64      `json:"stravaAthleteId"`
	ProfilePicture  *string    `json:"profilePicture"`
	PeaksCompleted  int        `json:"peaksCompleted"`
	TotalPeaks      int        `json:"totalPeaks"`
	TotalDistance   float64    `json:"totalDistance"`
	TotalElevation  float64    `json:"totalElevation"`
	TotalSummitCount int       `json:"totalSummitCount"`
	Progress        float64    `json:"progress"` // Percentage 0-100
	JoinedAt        time.Time  `json:"joinedAt"`
	CompletedAt     *time.Time `json:"completedAt"`
}

// ProposalStatus represents the review status of a challenge proposal
type ProposalStatus string

const (
	ProposalStatusPending  ProposalStatus = "pending"
	ProposalStatusApproved ProposalStatus = "approved"
	ProposalStatusRejected ProposalStatus = "rejected"
)

// ChallengeProposal represents a user-submitted challenge idea for admin review
type ChallengeProposal struct {
	ID                 int64           `json:"id" db:"id"`
	ProposedByUserID   int64           `json:"proposedByUserId" db:"proposed_by_user_id"`
	Name               string          `json:"name" db:"name"`
	Description        *string         `json:"description" db:"description"`
	GoalType           GoalType        `json:"goalType" db:"goal_type"`
	CompetitionMode    CompetitionMode `json:"competitionMode" db:"competition_mode"`
	TargetValue        *float64        `json:"targetValue" db:"target_value"`             // For distance/elevation (in meters)
	TargetSummitCount  *int            `json:"targetSummitCount" db:"target_summit_count"` // For summit_count
	PeakIDs            []int64         `json:"peakIds" db:"peak_ids"`                      // For specific_summits
	Region             *string         `json:"region" db:"region"`
	Difficulty         *string         `json:"difficulty" db:"difficulty"`
	Status             ProposalStatus  `json:"status" db:"status"`
	AdminNotes         *string         `json:"adminNotes" db:"admin_notes"`
	CreatedAt          time.Time       `json:"createdAt" db:"created_at"`
	ReviewedAt         *time.Time      `json:"reviewedAt" db:"reviewed_at"`
	ReviewedByUserID   *int64          `json:"reviewedByUserId" db:"reviewed_by_user_id"`
}
