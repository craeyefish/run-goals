package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"run-goals/daos"
	"run-goals/models"
	"strings"
	"time"
)

var (
	ErrChallengeNotFound    = errors.New("challenge not found")
	ErrNotChallengeOwner    = errors.New("user is not the challenge owner")
	ErrAlreadyParticipant   = errors.New("user is already a participant")
	ErrNotParticipant       = errors.New("user is not a participant")
	ErrChallengeTypeInvalid = errors.New("invalid challenge type")
	ErrChallengeNotPublic   = errors.New("challenge is not public")
)

type ChallengeServiceInterface interface {
	// Challenge CRUD
	CreateChallenge(userID int64, challenge models.Challenge, peakIDs []int64) (*models.Challenge, error)
	GetChallenge(id int64, userID *int64) (*models.ChallengeWithProgress, error)
	UpdateChallenge(id int64, userID int64, challenge models.Challenge) error
	DeleteChallenge(id int64, userID int64) error

	// Discovery
	GetUserChallenges(userID int64) ([]models.ChallengeWithProgress, error)
	GetFeaturedChallenges() ([]models.Challenge, error)
	GetPublicChallenges(region *string, limit int, offset int) ([]models.Challenge, error)
	SearchChallenges(query string, limit int) ([]models.Challenge, error)

	// Peaks
	GetChallengePeaks(challengeID int64, userID *int64) ([]models.ChallengePeakWithDetails, error)
	SetChallengePeaks(challengeID int64, userID int64, peakIDs []int64) error

	// Participation
	JoinChallenge(challengeID int64, userID int64) error
	JoinChallengeByCode(joinCode string, userID int64) (*models.Challenge, error)
	LeaveChallenge(challengeID int64, userID int64) error
	LockChallenge(challengeID int64, userID int64) error
	GetParticipants(challengeID int64) ([]models.ChallengeParticipantWithUser, error)
	GetLeaderboard(challengeID int64) ([]models.LeaderboardEntry, error)

	// Progress tracking
	RecordSummit(challengeID int64, userID int64, peakID int64, activityID *int64, summitedAt time.Time) error
	GetSummitLog(challengeID int64, userID *int64) ([]models.ChallengeSummitLogWithDetails, error)
	RefreshParticipantProgress(challengeID int64, userID int64) error
	RefreshAllChallengeProgress() error

	// Activities
	GetChallengeActivities(challengeID int64) ([]models.ActivityWithUser, error)

	// Group challenges
	AddGroupToChallenge(challengeID int64, groupID int64, deadlineOverride *time.Time) error
	RemoveGroupFromChallenge(challengeID int64, groupID int64) error
	GetGroupChallenges(groupID int64) ([]models.Challenge, error)
}

type ChallengeService struct {
	l            *log.Logger
	challengeDao *daos.ChallengeDao
	activityDao  *daos.ActivityDao
}

func NewChallengeService(
	l *log.Logger,
	challengeDao *daos.ChallengeDao,
	activityDao *daos.ActivityDao,
) *ChallengeService {
	return &ChallengeService{
		l:            l,
		challengeDao: challengeDao,
		activityDao:  activityDao,
	}
}

// ==================== Challenge CRUD ====================

// generateJoinCode creates a 6-character alphanumeric join code
func (s *ChallengeService) generateJoinCode() (string, error) {
	bytes := make([]byte, 3) // 3 bytes = 6 hex chars
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

func (s *ChallengeService) CreateChallenge(userID int64, challenge models.Challenge, peakIDs []int64) (*models.Challenge, error) {
	// Set creator
	challenge.CreatedByUserID = &userID
	challenge.CreatedAt = time.Now()
	challenge.UpdatedAt = time.Now()

	// Set defaults
	if challenge.ChallengeType == "" {
		challenge.ChallengeType = models.ChallengeTypeCustom
	}
	if challenge.GoalType == "" {
		challenge.GoalType = models.GoalTypeSpecificSummits
	}
	if challenge.CompetitionMode == "" {
		challenge.CompetitionMode = models.CompetitionModeCollaborative
	}
	if challenge.Visibility == "" {
		challenge.Visibility = models.VisibilityPrivate
	}

	// Generate join code if not provided
	if challenge.JoinCode == "" {
		joinCode, err := s.generateJoinCode()
		if err != nil {
			s.l.Printf("Error generating join code: %v", err)
			return nil, err
		}
		challenge.JoinCode = joinCode
	}

	// Create challenge
	id, err := s.challengeDao.CreateChallenge(challenge)
	if err != nil {
		s.l.Printf("Error creating challenge: %v", err)
		return nil, err
	}
	challenge.ID = *id

	// Add peaks if provided
	if len(peakIDs) > 0 {
		err = s.challengeDao.SetChallengePeaks(*id, peakIDs)
		if err != nil {
			s.l.Printf("Error setting challenge peaks: %v", err)
			// Don't fail the challenge creation, but log the error
		}
	}

	// Auto-join the creator
	err = s.challengeDao.JoinChallenge(*id, userID)
	if err != nil {
		s.l.Printf("Error auto-joining creator to challenge: %v", err)
		// Don't fail the challenge creation
	} else {
		// Calculate initial progress for the creator
		err = s.RefreshParticipantProgress(*id, userID)
		if err != nil {
			s.l.Printf("Warning: Failed to refresh creator's progress: %v", err)
		}
	}

	return &challenge, nil
}

func (s *ChallengeService) GetChallenge(id int64, userID *int64) (*models.ChallengeWithProgress, error) {
	challenge, err := s.challengeDao.GetChallengeByID(id)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrChallengeNotFound
	}

	// Build response with progress
	result := &models.ChallengeWithProgress{
		Challenge: *challenge,
	}

	// If user ID provided, check if they're a participant
	if userID != nil {
		isParticipant, _ := s.challengeDao.IsUserParticipant(id, *userID)
		result.IsJoined = isParticipant
	}

	// For collaborative challenges, show team progress (sum of all participants)
	// For competitive challenges, show individual progress
	if challenge.CompetitionMode == models.CompetitionModeCollaborative {
		// Get all participants and sum their progress
		participants, err := s.challengeDao.GetChallengeParticipants(id)
		if err == nil {
			for _, p := range participants {
				result.CompletedPeaks += p.PeaksCompleted
				result.CurrentDistance += p.TotalDistance
				result.CurrentElevation += p.TotalElevation
				result.CurrentSummitCount += p.TotalSummitCount
			}
		}
	} else {
		// Competitive mode: show individual progress
		if userID != nil && result.IsJoined {
			participant, err := s.challengeDao.GetChallengeParticipantByUserID(id, *userID)
			if err == nil && participant != nil {
				result.CompletedPeaks = participant.PeaksCompleted
				result.CurrentDistance = participant.TotalDistance
				result.CurrentElevation = participant.TotalElevation
				result.CurrentSummitCount = participant.TotalSummitCount
				result.IsCompleted = participant.CompletedAt != nil
			}
		}
	}

	// For specific_summits, get peak count (target)
	if challenge.GoalType == models.GoalTypeSpecificSummits {
		peaks, err := s.challengeDao.GetChallengePeaks(id)
		if err == nil {
			result.TotalPeaks = len(peaks)
		}
	}

	return result, nil
}

func (s *ChallengeService) UpdateChallenge(id int64, userID int64, challenge models.Challenge) error {
	// Check ownership
	existing, err := s.challengeDao.GetChallengeByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrChallengeNotFound
	}
	if existing.CreatedByUserID == nil || *existing.CreatedByUserID != userID {
		return ErrNotChallengeOwner
	}

	challenge.ID = id
	challenge.UpdatedAt = time.Now()
	return s.challengeDao.UpdateChallenge(challenge)
}

func (s *ChallengeService) DeleteChallenge(id int64, userID int64) error {
	// Check ownership
	existing, err := s.challengeDao.GetChallengeByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrChallengeNotFound
	}
	if existing.CreatedByUserID == nil || *existing.CreatedByUserID != userID {
		return ErrNotChallengeOwner
	}

	return s.challengeDao.DeleteChallenge(id)
}

// ==================== Discovery ====================

func (s *ChallengeService) GetUserChallenges(userID int64) ([]models.ChallengeWithProgress, error) {
	return s.challengeDao.GetChallengesByUser(userID)
}

func (s *ChallengeService) GetFeaturedChallenges() ([]models.Challenge, error) {
	return s.challengeDao.GetFeaturedChallenges()
}

func (s *ChallengeService) GetPublicChallenges(region *string, limit int, offset int) ([]models.Challenge, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.challengeDao.GetPublicChallenges(region, limit, offset)
}

func (s *ChallengeService) SearchChallenges(query string, limit int) ([]models.Challenge, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.challengeDao.SearchChallenges(query, limit)
}

// ==================== Peaks ====================

func (s *ChallengeService) GetChallengePeaks(challengeID int64, userID *int64) ([]models.ChallengePeakWithDetails, error) {
	if userID != nil {
		return s.challengeDao.GetChallengePeaksWithUserProgress(challengeID, *userID)
	}
	return s.challengeDao.GetChallengePeaks(challengeID)
}

func (s *ChallengeService) SetChallengePeaks(challengeID int64, userID int64, peakIDs []int64) error {
	// Check ownership
	existing, err := s.challengeDao.GetChallengeByID(challengeID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrChallengeNotFound
	}
	if existing.CreatedByUserID == nil || *existing.CreatedByUserID != userID {
		return ErrNotChallengeOwner
	}

	return s.challengeDao.SetChallengePeaks(challengeID, peakIDs)
}

// ==================== Participation ====================

func (s *ChallengeService) JoinChallenge(challengeID int64, userID int64) error {
	// Check challenge exists and is joinable
	challenge, err := s.challengeDao.GetChallengeByID(challengeID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return ErrChallengeNotFound
	}

	// Check visibility - only public challenges can be joined by anyone
	// Private/friends challenges would need invitation logic (future feature)
	if challenge.Visibility == models.VisibilityPrivate {
		// For now, allow joining private challenges if user somehow has the ID
		// In future, add invitation system
	}

	// Check if already participant
	isParticipant, err := s.challengeDao.IsUserParticipant(challengeID, userID)
	if err != nil {
		return err
	}
	if isParticipant {
		return ErrAlreadyParticipant
	}

	// Join the challenge
	err = s.challengeDao.JoinChallenge(challengeID, userID)
	if err != nil {
		return err
	}

	// Calculate initial progress
	err = s.RefreshParticipantProgress(challengeID, userID)
	if err != nil {
		s.l.Printf("Warning: Failed to refresh participant progress: %v", err)
		// Don't fail the join, just log the error
	}

	return nil
}

func (s *ChallengeService) JoinChallengeByCode(joinCode string, userID int64) (*models.Challenge, error) {
	// Find challenge by join code
	challenge, err := s.challengeDao.GetChallengeByJoinCode(joinCode)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrChallengeNotFound
	}

	// Check if already participant
	isParticipant, err := s.challengeDao.IsUserParticipant(challenge.ID, userID)
	if err != nil {
		return nil, err
	}
	if isParticipant {
		return nil, ErrAlreadyParticipant
	}

	// Join the challenge
	err = s.challengeDao.JoinChallenge(challenge.ID, userID)
	if err != nil {
		return nil, err
	}

	// Calculate initial progress
	err = s.RefreshParticipantProgress(challenge.ID, userID)
	if err != nil {
		s.l.Printf("Warning: Failed to refresh participant progress: %v", err)
		// Don't fail the join, just log the error
	}

	return challenge, nil
}

func (s *ChallengeService) LeaveChallenge(challengeID int64, userID int64) error {
	// Check if participant
	isParticipant, err := s.challengeDao.IsUserParticipant(challengeID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return ErrNotParticipant
	}

	return s.challengeDao.LeaveChallenge(challengeID, userID)
}

func (s *ChallengeService) LockChallenge(challengeID int64, userID int64) error {
	// Check challenge exists
	challenge, err := s.challengeDao.GetChallengeByID(challengeID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return ErrChallengeNotFound
	}

	// Check ownership
	if challenge.CreatedByUserID == nil || *challenge.CreatedByUserID != userID {
		return ErrNotChallengeOwner
	}

	// Lock the challenge (irreversible)
	return s.challengeDao.LockChallenge(challengeID)
}

func (s *ChallengeService) GetParticipants(challengeID int64) ([]models.ChallengeParticipantWithUser, error) {
	return s.challengeDao.GetChallengeParticipants(challengeID)
}

func (s *ChallengeService) GetLeaderboard(challengeID int64) ([]models.LeaderboardEntry, error) {
	return s.challengeDao.GetChallengeLeaderboard(challengeID)
}

// ==================== Progress Tracking ====================

func (s *ChallengeService) RecordSummit(challengeID int64, userID int64, peakID int64, activityID *int64, summitedAt time.Time) error {
	// Check if user is participant
	isParticipant, err := s.challengeDao.IsUserParticipant(challengeID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return ErrNotParticipant
	}

	// Check if already summited this peak for this challenge
	hasSummited, err := s.challengeDao.HasUserSummitedPeakForChallenge(challengeID, userID, peakID)
	if err != nil {
		return err
	}
	if hasSummited {
		return nil // Already logged, no error
	}

	// Log the summit
	logEntry := models.ChallengeSummitLog{
		ChallengeID: challengeID,
		UserID:      userID,
		PeakID:      &peakID,
		ActivityID:  activityID,
		SummitedAt:  summitedAt,
	}
	err = s.challengeDao.LogSummit(logEntry)
	if err != nil {
		return err
	}

	// Refresh progress
	return s.RefreshParticipantProgress(challengeID, userID)
}

func (s *ChallengeService) GetSummitLog(challengeID int64, userID *int64) ([]models.ChallengeSummitLogWithDetails, error) {
	return s.challengeDao.GetChallengeSummitLog(challengeID, userID)
}

func (s *ChallengeService) RefreshParticipantProgress(challengeID int64, userID int64) error {
	// Get the challenge to know its goal type and date range
	challenge, err := s.challengeDao.GetChallengeByID(challengeID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return ErrChallengeNotFound
	}

	var peaksCompleted int
	var totalPeaks int
	var totalDistance float64
	var totalElevation float64
	var totalSummitCount int
	var isCompleted bool

	switch challenge.GoalType {
	case models.GoalTypeSpecificSummits:
		// Get total peaks
		peaks, err := s.challengeDao.GetChallengePeaks(challengeID)
		if err != nil {
			return err
		}
		totalPeaks = len(peaks)

		// Get completed peaks
		summitLog, err := s.challengeDao.GetChallengeSummitLog(challengeID, &userID)
		if err != nil {
			return err
		}
		peaksCompleted = len(summitLog)
		isCompleted = peaksCompleted >= totalPeaks && totalPeaks > 0

	case models.GoalTypeDistance:
		// Get activities within challenge date range and sum distance
		activities, err := s.activityDao.GetActivitiesByUserIDAndDateRange(userID, challenge.StartDate, challenge.Deadline)
		if err != nil {
			return err
		}
		for _, activity := range activities {
			totalDistance += activity.Distance
		}
		// Check completion
		if challenge.TargetValue != nil {
			isCompleted = totalDistance >= *challenge.TargetValue
		}

	case models.GoalTypeElevation:
		// Get activities within challenge date range and sum elevation
		activities, err := s.activityDao.GetActivitiesByUserIDAndDateRange(userID, challenge.StartDate, challenge.Deadline)
		if err != nil {
			return err
		}
		for _, activity := range activities {
			totalElevation += activity.Elevation
		}
		// Check completion
		if challenge.TargetValue != nil {
			isCompleted = totalElevation >= *challenge.TargetValue
		}

	case models.GoalTypeSummitCount:
		// Get summit log within challenge date range
		summitLog, err := s.challengeDao.GetChallengeSummitLog(challengeID, &userID)
		if err != nil {
			return err
		}
		totalSummitCount = len(summitLog)
		// Check completion
		if challenge.TargetSummitCount != nil {
			isCompleted = totalSummitCount >= *challenge.TargetSummitCount
		}
	}

	// Update progress with all fields
	err = s.challengeDao.UpdateParticipantProgressFull(
		challengeID, userID,
		peaksCompleted, totalPeaks,
		totalDistance, totalElevation, totalSummitCount,
	)
	if err != nil {
		return err
	}

	// Mark as completed if applicable
	if isCompleted {
		return s.challengeDao.MarkParticipantCompleted(challengeID, userID)
	}

	return nil
}

// RefreshAllChallengeProgress refreshes progress for all active challenges and their participants
// This should be called after syncing activities to update distance/elevation progress
func (s *ChallengeService) RefreshAllChallengeProgress() error {
	// Get all active (non-completed) challenges
	participants, err := s.challengeDao.GetAllActiveParticipants()
	if err != nil {
		s.l.Printf("Error getting active participants: %v", err)
		return err
	}

	s.l.Printf("Refreshing progress for %d active challenge participants", len(participants))

	for _, participant := range participants {
		err := s.RefreshParticipantProgress(participant.ChallengeID, participant.UserID)
		if err != nil {
			s.l.Printf("Error refreshing progress for challenge %d user %d: %v",
				participant.ChallengeID, participant.UserID, err)
			// Continue with other participants even if one fails
		}
	}

	return nil
}

// ==================== Group Challenges ====================

func (s *ChallengeService) AddGroupToChallenge(challengeID int64, groupID int64, deadlineOverride *time.Time) error {
	// Verify challenge exists
	challenge, err := s.challengeDao.GetChallengeByID(challengeID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return ErrChallengeNotFound
	}

	return s.challengeDao.AddGroupToChallenge(challengeID, groupID, deadlineOverride)
}

func (s *ChallengeService) RemoveGroupFromChallenge(challengeID int64, groupID int64) error {
	return s.challengeDao.RemoveGroupFromChallenge(challengeID, groupID)
}

func (s *ChallengeService) GetGroupChallenges(groupID int64) ([]models.Challenge, error) {
	return s.challengeDao.GetGroupChallenges(groupID)
}

// ==================== Auto-detection hook ====================

// ProcessActivityForChallenges is called when an activity is processed to auto-credit summits
// This should be called from the summit detection workflow
func (s *ChallengeService) ProcessActivityForChallenges(userID int64, peakID int64, activityID int64, summitedAt time.Time) error {
	// Get all challenges the user is participating in
	challenges, err := s.challengeDao.GetChallengesByUser(userID)
	if err != nil {
		s.l.Printf("Error getting user challenges: %v", err)
		return err
	}

	// For each challenge, check if this summit should be credited
	for _, challenge := range challenges {
		// Check if activity is within challenge date range
		if challenge.StartDate != nil && summitedAt.Before(*challenge.StartDate) {
			continue
		}
		if challenge.Deadline != nil && summitedAt.After(*challenge.Deadline) {
			continue
		}

		// Handle based on goal type
		switch challenge.GoalType {
		case models.GoalTypeSpecificSummits:
			// For specific_summits, only credit if peak is in the challenge list
			peaks, err := s.challengeDao.GetChallengePeaks(challenge.ID)
			if err != nil {
				s.l.Printf("Error getting challenge %d peaks: %v", challenge.ID, err)
				continue
			}

			for _, peak := range peaks {
				if peak.PeakID == peakID {
					// This summit counts for this challenge
					err = s.RecordSummit(challenge.ID, userID, peakID, &activityID, summitedAt)
					if err != nil {
						s.l.Printf("Error recording summit for challenge %d: %v", challenge.ID, err)
					}
					break
				}
			}

		case models.GoalTypeSummitCount:
			// For summit_count, credit ANY summit within date range
			err = s.RecordSummit(challenge.ID, userID, peakID, &activityID, summitedAt)
			if err != nil {
				s.l.Printf("Error recording summit for challenge %d: %v", challenge.ID, err)
			}
		}
	}

	return nil
}

// ==================== Activities ====================

func (s *ChallengeService) GetChallengeActivities(challengeID int64) ([]models.ActivityWithUser, error) {
	return s.challengeDao.GetChallengeActivities(challengeID)
}
