package dto

import (
	"run-goals/models"
	"time"
)

// ==================== Challenge Requests ====================

type CreateChallengeRequest struct {
	Name            string                  `json:"name"`
	Description     *string                 `json:"description"`
	ChallengeType   models.ChallengeType    `json:"challengeType"`
	CompetitionMode models.CompetitionMode  `json:"competitionMode"`
	Visibility      models.Visibility       `json:"visibility"`
	StartDate       *time.Time              `json:"startDate"`
	Deadline        *time.Time              `json:"deadline"`
	TargetCount     *int                    `json:"targetCount"`
	Region          *string                 `json:"region"`
	Difficulty      *string                 `json:"difficulty"`
	PeakIDs         []int64                 `json:"peakIds"`
}

type UpdateChallengeRequest struct {
	ID              int64                   `json:"id"`
	Name            string                  `json:"name"`
	Description     *string                 `json:"description"`
	ChallengeType   models.ChallengeType    `json:"challengeType"`
	CompetitionMode models.CompetitionMode  `json:"competitionMode"`
	Visibility      models.Visibility       `json:"visibility"`
	StartDate       *time.Time              `json:"startDate"`
	Deadline        *time.Time              `json:"deadline"`
	TargetCount     *int                    `json:"targetCount"`
	Region          *string                 `json:"region"`
	Difficulty      *string                 `json:"difficulty"`
}

type SetChallengePeaksRequest struct {
	PeakIDs []int64 `json:"peakIds"`
}

type AddGroupToChallengeRequest struct {
	GroupID          int64      `json:"groupId"`
	DeadlineOverride *time.Time `json:"deadlineOverride"`
}

type RecordSummitRequest struct {
	PeakID     int64     `json:"peakId"`
	ActivityID *int64    `json:"activityId"`
	SummitedAt time.Time `json:"summitedAt"`
}

// ==================== Challenge Responses ====================

type CreateChallengeResponse struct {
	ID int64 `json:"id"`
}

type ChallengeDetailResponse struct {
	Challenge    models.ChallengeWithProgress       `json:"challenge"`
	Peaks        []models.ChallengePeakWithDetails  `json:"peaks"`
	Participants []models.ChallengeParticipantWithUser `json:"participants"`
}

type ChallengeListResponse struct {
	Challenges []models.ChallengeWithProgress `json:"challenges"`
	Total      int                            `json:"total"`
}

type PublicChallengeListResponse struct {
	Challenges []models.Challenge `json:"challenges"`
	Total      int                `json:"total"`
}
