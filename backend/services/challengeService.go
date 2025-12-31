package services

import (
	"errors"
	"log"
	"run-goals/daos"
	"run-goals/models"
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
	LeaveChallenge(challengeID int64, userID int64) error
	GetParticipants(challengeID int64) ([]models.ChallengeParticipantWithUser, error)
	GetLeaderboard(challengeID int64) ([]models.LeaderboardEntry, error)

	// Progress tracking
	RecordSummit(challengeID int64, userID int64, peakID int64, activityID *int64, summitedAt time.Time) error
	GetSummitLog(challengeID int64, userID *int64) ([]models.ChallengeSummitLogWithDetails, error)
	RefreshParticipantProgress(challengeID int64, userID int64) error

	// Group challenges
	AddGroupToChallenge(challengeID int64, groupID int64, deadlineOverride *time.Time) error
	RemoveGroupFromChallenge(challengeID int64, groupID int64) error
	GetGroupChallenges(groupID int64) ([]models.Challenge, error)
}

type ChallengeService struct {
	l            *log.Logger
	challengeDao *daos.ChallengeDao
}

func NewChallengeService(
	l *log.Logger,
	challengeDao *daos.ChallengeDao,
) *ChallengeService {
	return &ChallengeService{
		l:            l,
		challengeDao: challengeDao,
	}
}

// ==================== Challenge CRUD ====================

func (s *ChallengeService) CreateChallenge(userID int64, challenge models.Challenge, peakIDs []int64) (*models.Challenge, error) {
	// Set creator
	challenge.CreatedByUserID = &userID
	challenge.CreatedAt = time.Now()
	challenge.UpdatedAt = time.Now()

	// Set defaults
	if challenge.ChallengeType == "" {
		challenge.ChallengeType = models.ChallengeTypeCustom
	}
	if challenge.CompetitionMode == "" {
		challenge.CompetitionMode = models.CompetitionModeCollaborative
	}
	if challenge.Visibility == "" {
		challenge.Visibility = models.VisibilityPrivate
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

	// Get peak count
	peaks, err := s.challengeDao.GetChallengePeaks(id)
	if err == nil {
		result.TotalPeaks = len(peaks)
	}

	// If user ID provided, get their progress
	if userID != nil {
		isParticipant, _ := s.challengeDao.IsUserParticipant(id, *userID)
		result.IsJoined = isParticipant

		if isParticipant {
			summitLog, _ := s.challengeDao.GetChallengeSummitLog(id, userID)
			result.CompletedPeaks = len(summitLog)
			result.IsCompleted = result.CompletedPeaks >= result.TotalPeaks && result.TotalPeaks > 0
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

	return s.challengeDao.JoinChallenge(challengeID, userID)
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
	// Get total peaks
	peaks, err := s.challengeDao.GetChallengePeaks(challengeID)
	if err != nil {
		return err
	}
	totalPeaks := len(peaks)

	// Get completed peaks
	summitLog, err := s.challengeDao.GetChallengeSummitLog(challengeID, &userID)
	if err != nil {
		return err
	}
	completedPeaks := len(summitLog)

	// Update progress
	err = s.challengeDao.UpdateParticipantProgress(challengeID, userID, completedPeaks, totalPeaks)
	if err != nil {
		return err
	}

	// Check if completed
	if completedPeaks >= totalPeaks && totalPeaks > 0 {
		return s.challengeDao.MarkParticipantCompleted(challengeID, userID)
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

	// For each challenge, check if this peak is part of it
	for _, challenge := range challenges {
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
	}

	return nil
}
