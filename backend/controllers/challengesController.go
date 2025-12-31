package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"run-goals/dto"
	"run-goals/meta"
	"run-goals/models"
	"run-goals/services"
	"strconv"
)

type ChallengesControllerInterface interface {
	// Challenge CRUD
	CreateChallenge(rw http.ResponseWriter, r *http.Request)
	GetChallenge(rw http.ResponseWriter, r *http.Request)
	UpdateChallenge(rw http.ResponseWriter, r *http.Request)
	DeleteChallenge(rw http.ResponseWriter, r *http.Request)

	// Discovery
	GetUserChallenges(rw http.ResponseWriter, r *http.Request)
	GetFeaturedChallenges(rw http.ResponseWriter, r *http.Request)
	GetPublicChallenges(rw http.ResponseWriter, r *http.Request)
	SearchChallenges(rw http.ResponseWriter, r *http.Request)

	// Peaks
	GetChallengePeaks(rw http.ResponseWriter, r *http.Request)
	SetChallengePeaks(rw http.ResponseWriter, r *http.Request)

	// Participation
	JoinChallenge(rw http.ResponseWriter, r *http.Request)
	LeaveChallenge(rw http.ResponseWriter, r *http.Request)
	GetParticipants(rw http.ResponseWriter, r *http.Request)
	GetLeaderboard(rw http.ResponseWriter, r *http.Request)

	// Progress
	GetSummitLog(rw http.ResponseWriter, r *http.Request)
	RecordSummit(rw http.ResponseWriter, r *http.Request)

	// Groups
	AddGroupToChallenge(rw http.ResponseWriter, r *http.Request)
	RemoveGroupFromChallenge(rw http.ResponseWriter, r *http.Request)
	GetGroupChallenges(rw http.ResponseWriter, r *http.Request)
}

type ChallengesController struct {
	l                *log.Logger
	challengeService *services.ChallengeService
}

func NewChallengesController(
	l *log.Logger,
	challengeService *services.ChallengeService,
) *ChallengesController {
	return &ChallengesController{
		l:                l,
		challengeService: challengeService,
	}
}

// ==================== Challenge CRUD ====================

func (c *ChallengesController) CreateChallenge(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle POST challenges - creating new challenge")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	var request dto.CreateChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Convert request to model
	challenge := models.Challenge{
		Name:            request.Name,
		Description:     request.Description,
		ChallengeType:   request.ChallengeType,
		CompetitionMode: request.CompetitionMode,
		Visibility:      request.Visibility,
		StartDate:       request.StartDate,
		Deadline:        request.Deadline,
		TargetCount:     request.TargetCount,
		Region:          request.Region,
		Difficulty:      request.Difficulty,
	}

	created, err := c.challengeService.CreateChallenge(userID, challenge, request.PeakIDs)
	if err != nil {
		c.l.Printf("Error creating challenge: %v", err)
		http.Error(rw, "Failed to create challenge", http.StatusInternalServerError)
		return
	}

	response := dto.CreateChallengeResponse{ID: created.ID}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

func (c *ChallengesController) GetChallenge(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET challenge - getting challenge by ID")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	userID, _ := meta.GetUserIDFromContext(r.Context())
	var userIDPtr *int64
	if userID > 0 {
		userIDPtr = &userID
	}

	challenge, err := c.challengeService.GetChallenge(challengeID, userIDPtr)
	if err != nil {
		if errors.Is(err, services.ErrChallengeNotFound) {
			http.Error(rw, "Challenge not found", http.StatusNotFound)
			return
		}
		c.l.Printf("Error getting challenge: %v", err)
		http.Error(rw, "Failed to get challenge", http.StatusInternalServerError)
		return
	}

	// Get peaks with user progress
	peaks, _ := c.challengeService.GetChallengePeaks(challengeID, userIDPtr)

	// Get participants
	participants, _ := c.challengeService.GetParticipants(challengeID)

	response := dto.ChallengeDetailResponse{
		Challenge:    *challenge,
		Peaks:        peaks,
		Participants: participants,
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

func (c *ChallengesController) UpdateChallenge(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle PUT challenges - updating challenge")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	var request dto.UpdateChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	challenge := models.Challenge{
		ID:              request.ID,
		Name:            request.Name,
		Description:     request.Description,
		ChallengeType:   request.ChallengeType,
		CompetitionMode: request.CompetitionMode,
		Visibility:      request.Visibility,
		StartDate:       request.StartDate,
		Deadline:        request.Deadline,
		TargetCount:     request.TargetCount,
		Region:          request.Region,
		Difficulty:      request.Difficulty,
	}

	err := c.challengeService.UpdateChallenge(request.ID, userID, challenge)
	if err != nil {
		if errors.Is(err, services.ErrChallengeNotFound) {
			http.Error(rw, "Challenge not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, services.ErrNotChallengeOwner) {
			http.Error(rw, "Not authorized", http.StatusForbidden)
			return
		}
		c.l.Printf("Error updating challenge: %v", err)
		http.Error(rw, "Failed to update challenge", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (c *ChallengesController) DeleteChallenge(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle DELETE challenges - deleting challenge")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	userID, _ := meta.GetUserIDFromContext(r.Context())

	err = c.challengeService.DeleteChallenge(challengeID, userID)
	if err != nil {
		if errors.Is(err, services.ErrChallengeNotFound) {
			http.Error(rw, "Challenge not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, services.ErrNotChallengeOwner) {
			http.Error(rw, "Not authorized", http.StatusForbidden)
			return
		}
		c.l.Printf("Error deleting challenge: %v", err)
		http.Error(rw, "Failed to delete challenge", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// ==================== Discovery ====================

func (c *ChallengesController) GetUserChallenges(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET user-challenges")

	userID, _ := meta.GetUserIDFromContext(r.Context())

	challenges, err := c.challengeService.GetUserChallenges(userID)
	if err != nil {
		c.l.Printf("Error getting user challenges: %v", err)
		http.Error(rw, "Failed to get challenges", http.StatusInternalServerError)
		return
	}

	response := dto.ChallengeListResponse{
		Challenges: challenges,
		Total:      len(challenges),
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

func (c *ChallengesController) GetFeaturedChallenges(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET featured-challenges")

	challenges, err := c.challengeService.GetFeaturedChallenges()
	if err != nil {
		c.l.Printf("Error getting featured challenges: %v", err)
		http.Error(rw, "Failed to get challenges", http.StatusInternalServerError)
		return
	}

	response := dto.PublicChallengeListResponse{
		Challenges: challenges,
		Total:      len(challenges),
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

func (c *ChallengesController) GetPublicChallenges(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET public-challenges")

	region := r.URL.Query().Get("region")
	var regionPtr *string
	if region != "" {
		regionPtr = &region
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	challenges, err := c.challengeService.GetPublicChallenges(regionPtr, limit, offset)
	if err != nil {
		c.l.Printf("Error getting public challenges: %v", err)
		http.Error(rw, "Failed to get challenges", http.StatusInternalServerError)
		return
	}

	response := dto.PublicChallengeListResponse{
		Challenges: challenges,
		Total:      len(challenges),
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

func (c *ChallengesController) SearchChallenges(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET search-challenges")

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(rw, "Missing search query", http.StatusBadRequest)
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	challenges, err := c.challengeService.SearchChallenges(query, limit)
	if err != nil {
		c.l.Printf("Error searching challenges: %v", err)
		http.Error(rw, "Failed to search challenges", http.StatusInternalServerError)
		return
	}

	response := dto.PublicChallengeListResponse{
		Challenges: challenges,
		Total:      len(challenges),
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

// ==================== Peaks ====================

func (c *ChallengesController) GetChallengePeaks(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET challenge-peaks")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	userID, _ := meta.GetUserIDFromContext(r.Context())
	var userIDPtr *int64
	if userID > 0 {
		userIDPtr = &userID
	}

	peaks, err := c.challengeService.GetChallengePeaks(challengeID, userIDPtr)
	if err != nil {
		c.l.Printf("Error getting challenge peaks: %v", err)
		http.Error(rw, "Failed to get peaks", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(peaks)
}

func (c *ChallengesController) SetChallengePeaks(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle PUT challenge-peaks")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	userID, _ := meta.GetUserIDFromContext(r.Context())

	var request dto.SetChallengePeaksRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = c.challengeService.SetChallengePeaks(challengeID, userID, request.PeakIDs)
	if err != nil {
		if errors.Is(err, services.ErrChallengeNotFound) {
			http.Error(rw, "Challenge not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, services.ErrNotChallengeOwner) {
			http.Error(rw, "Not authorized", http.StatusForbidden)
			return
		}
		c.l.Printf("Error setting challenge peaks: %v", err)
		http.Error(rw, "Failed to set peaks", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// ==================== Participation ====================

func (c *ChallengesController) JoinChallenge(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle POST challenge-join")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	userID, _ := meta.GetUserIDFromContext(r.Context())

	err = c.challengeService.JoinChallenge(challengeID, userID)
	if err != nil {
		if errors.Is(err, services.ErrChallengeNotFound) {
			http.Error(rw, "Challenge not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, services.ErrAlreadyParticipant) {
			http.Error(rw, "Already a participant", http.StatusConflict)
			return
		}
		c.l.Printf("Error joining challenge: %v", err)
		http.Error(rw, "Failed to join challenge", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (c *ChallengesController) LeaveChallenge(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle DELETE challenge-leave")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	userID, _ := meta.GetUserIDFromContext(r.Context())

	err = c.challengeService.LeaveChallenge(challengeID, userID)
	if err != nil {
		if errors.Is(err, services.ErrNotParticipant) {
			http.Error(rw, "Not a participant", http.StatusBadRequest)
			return
		}
		c.l.Printf("Error leaving challenge: %v", err)
		http.Error(rw, "Failed to leave challenge", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (c *ChallengesController) GetParticipants(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET challenge-participants")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	participants, err := c.challengeService.GetParticipants(challengeID)
	if err != nil {
		c.l.Printf("Error getting participants: %v", err)
		http.Error(rw, "Failed to get participants", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(participants)
}

func (c *ChallengesController) GetLeaderboard(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET challenge-leaderboard")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	leaderboard, err := c.challengeService.GetLeaderboard(challengeID)
	if err != nil {
		c.l.Printf("Error getting leaderboard: %v", err)
		http.Error(rw, "Failed to get leaderboard", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(leaderboard)
}

// ==================== Progress ====================

func (c *ChallengesController) GetSummitLog(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET challenge-summit-log")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	// Optional user filter
	var userIDPtr *int64
	if userIDStr := r.URL.Query().Get("userId"); userIDStr != "" {
		if uid, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userIDPtr = &uid
		}
	}

	summitLog, err := c.challengeService.GetSummitLog(challengeID, userIDPtr)
	if err != nil {
		c.l.Printf("Error getting summit log: %v", err)
		http.Error(rw, "Failed to get summit log", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(summitLog)
}

func (c *ChallengesController) RecordSummit(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle POST challenge-summit - manually recording summit")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	userID, _ := meta.GetUserIDFromContext(r.Context())

	var request dto.RecordSummitRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = c.challengeService.RecordSummit(challengeID, userID, request.PeakID, request.ActivityID, request.SummitedAt)
	if err != nil {
		if errors.Is(err, services.ErrNotParticipant) {
			http.Error(rw, "Not a participant", http.StatusForbidden)
			return
		}
		c.l.Printf("Error recording summit: %v", err)
		http.Error(rw, "Failed to record summit", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// ==================== Groups ====================

func (c *ChallengesController) AddGroupToChallenge(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle POST challenge-group")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	var request dto.AddGroupToChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		c.l.Printf("Error unmarshalling data: %v", err)
		http.Error(rw, "Error unmarshalling data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = c.challengeService.AddGroupToChallenge(challengeID, request.GroupID, request.DeadlineOverride)
	if err != nil {
		if errors.Is(err, services.ErrChallengeNotFound) {
			http.Error(rw, "Challenge not found", http.StatusNotFound)
			return
		}
		c.l.Printf("Error adding group to challenge: %v", err)
		http.Error(rw, "Failed to add group", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (c *ChallengesController) RemoveGroupFromChallenge(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle DELETE challenge-group")

	challengeID, err := c.getChallengeIDFromURL(r)
	if err != nil {
		http.Error(rw, "Invalid challenge ID", http.StatusBadRequest)
		return
	}

	groupIDStr := r.URL.Query().Get("groupId")
	if groupIDStr == "" {
		http.Error(rw, "Missing groupId", http.StatusBadRequest)
		return
	}
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid group ID", http.StatusBadRequest)
		return
	}

	err = c.challengeService.RemoveGroupFromChallenge(challengeID, groupID)
	if err != nil {
		c.l.Printf("Error removing group from challenge: %v", err)
		http.Error(rw, "Failed to remove group", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (c *ChallengesController) GetGroupChallenges(rw http.ResponseWriter, r *http.Request) {
	c.l.Printf("Handle GET group-challenges")

	groupIDStr := r.URL.Query().Get("groupId")
	if groupIDStr == "" {
		http.Error(rw, "Missing groupId", http.StatusBadRequest)
		return
	}
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid group ID", http.StatusBadRequest)
		return
	}

	challenges, err := c.challengeService.GetGroupChallenges(groupID)
	if err != nil {
		c.l.Printf("Error getting group challenges: %v", err)
		http.Error(rw, "Failed to get challenges", http.StatusInternalServerError)
		return
	}

	response := dto.PublicChallengeListResponse{
		Challenges: challenges,
		Total:      len(challenges),
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

// ==================== Helpers ====================

func (c *ChallengesController) getChallengeIDFromURL(r *http.Request) (int64, error) {
	idStr := r.URL.Query().Get("challengeId")
	if idStr == "" {
		idStr = r.URL.Query().Get("id")
	}
	if idStr == "" {
		return 0, errors.New("missing challenge ID")
	}
	return strconv.ParseInt(idStr, 10, 64)
}
